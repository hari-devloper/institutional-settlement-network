package settlement

import (
	"errors"
	"fmt"
	"math"
	"sync"
)

var (
	ErrInvalidAccount        = errors.New("account ID must not be empty")
	ErrAccountExists         = errors.New("account already exists")
	ErrAccountNotFound       = errors.New("account not found")
	ErrInvalidAmount         = errors.New("amount must be positive")
	ErrSameAccount           = errors.New("source and destination must differ")
	ErrInsufficientFunds     = errors.New("insufficient funds")
	ErrIdempotencyConflict   = errors.New("idempotency key was reused with different transfer data")
	ErrBalanceOverflow       = errors.New("destination balance would overflow int64")
	ErrInvalidIdempotencyKey = errors.New("idempotency key must not be empty")
)

type Receipt struct {
	ID       string
	Sequence uint64
	From     string
	To       string
	Amount   int64
}

type transferRequest struct {
	from   string
	to     string
	amount int64
}

// Ledger is an in-memory teaching model, not persistent blockchain state.
// Its mutex makes each operation atomic within this process only.
type Ledger struct {
	mu           sync.Mutex
	balances     map[string]int64
	transfers    map[string]Receipt
	requests     map[string]transferRequest
	nextSequence uint64
}

func NewLedger() *Ledger {
	return &Ledger{
		balances:  make(map[string]int64),
		transfers: make(map[string]Receipt),
		requests:  make(map[string]transferRequest),
	}
}

// OpenAccount is a fixture/genesis-like operation. Production issuance and
// account-creation authority must be defined separately.
func (l *Ledger) OpenAccount(accountID string, openingBalance int64) error {
	if accountID == "" {
		return ErrInvalidAccount
	}
	if openingBalance < 0 {
		return errors.New("opening balance must not be negative")
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if _, exists := l.balances[accountID]; exists {
		return ErrAccountExists
	}
	l.balances[accountID] = openingBalance
	return nil
}

func (l *Ledger) Balance(accountID string) (int64, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	balance, ok := l.balances[accountID]
	if !ok {
		return 0, ErrAccountNotFound
	}
	return balance, nil
}

func (l *Ledger) Transfer(key, from, to string, amount int64) (Receipt, error) {
	if key == "" {
		return Receipt{}, ErrInvalidIdempotencyKey
	}
	if from == "" || to == "" {
		return Receipt{}, ErrInvalidAccount
	}
	if amount <= 0 {
		return Receipt{}, ErrInvalidAmount
	}
	if from == to {
		return Receipt{}, ErrSameAccount
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	request := transferRequest{from: from, to: to, amount: amount}
	if receipt, exists := l.transfers[key]; exists {
		if l.requests[key] != request {
			return Receipt{}, ErrIdempotencyConflict
		}
		return receipt, nil
	}

	sourceBalance, sourceExists := l.balances[from]
	if !sourceExists {
		return Receipt{}, fmt.Errorf("%w: %s", ErrAccountNotFound, from)
	}
	destinationBalance, destinationExists := l.balances[to]
	if !destinationExists {
		return Receipt{}, fmt.Errorf("%w: %s", ErrAccountNotFound, to)
	}
	if sourceBalance < amount {
		return Receipt{}, ErrInsufficientFunds
	}
	if destinationBalance > math.MaxInt64-amount {
		return Receipt{}, ErrBalanceOverflow
	}
	if l.nextSequence == math.MaxUint64 {
		return Receipt{}, errors.New("transfer sequence exhausted")
	}

	// Validate everything before mutating either balance.
	l.nextSequence++
	receipt := Receipt{ID: key, Sequence: l.nextSequence, From: from, To: to, Amount: amount}
	l.balances[from] = sourceBalance - amount
	l.balances[to] = destinationBalance + amount
	l.transfers[key] = receipt
	l.requests[key] = request
	return receipt, nil
}

func (l *Ledger) TotalBalance() (int64, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	var total int64
	for _, balance := range l.balances {
		if balance > math.MaxInt64-total {
			return 0, errors.New("total balance exceeds int64")
		}
		total += balance
	}
	return total, nil
}
