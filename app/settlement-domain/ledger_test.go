package settlement

import (
	"errors"
	"math"
	"sync"
	"sync/atomic"
	"testing"
)

func mustOpen(t *testing.T, l *Ledger, id string, balance int64) {
	t.Helper()
	if err := l.OpenAccount(id, balance); err != nil {
		t.Fatalf("OpenAccount(%q): %v", id, err)
	}
}

func TestTransferMovesFundsAndConservesTotal(t *testing.T) {
	l := NewLedger()
	mustOpen(t, l, "bank-a", 1000)
	mustOpen(t, l, "bank-b", 250)
	before, err := l.TotalBalance()
	if err != nil { t.Fatal(err) }

	receipt, err := l.Transfer("settlement-001", "bank-a", "bank-b", 300)
	if err != nil { t.Fatal(err) }
	after, err := l.TotalBalance()
	if err != nil { t.Fatal(err) }

	if receipt.Sequence != 1 || receipt.Amount != 300 { t.Fatalf("unexpected receipt: %+v", receipt) }
	a, _ := l.Balance("bank-a")
	b, _ := l.Balance("bank-b")
	if a != 700 || b != 550 { t.Fatalf("balances = (%d, %d), want (700, 550)", a, b) }
	if before != after { t.Fatalf("total changed from %d to %d", before, after) }
}

func TestRejectsInvalidTransfersWithoutChangingBalances(t *testing.T) {
	tests := []struct {
		name string
		key, from, to string
		amount int64
		want error
	}{
		{"empty key", "", "a", "b", 1, ErrInvalidIdempotencyKey},
		{"zero amount", "k", "a", "b", 0, ErrInvalidAmount},
		{"negative amount", "k", "a", "b", -1, ErrInvalidAmount},
		{"same account", "k", "a", "a", 1, ErrSameAccount},
		{"unknown source", "k", "missing", "b", 1, ErrAccountNotFound},
		{"unknown destination", "k", "a", "missing", 1, ErrAccountNotFound},
		{"insufficient funds", "k", "a", "b", 101, ErrInsufficientFunds},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			l := NewLedger()
			mustOpen(t, l, "a", 100)
			mustOpen(t, l, "b", 50)
			before, _ := l.TotalBalance()
			_, err := l.Transfer(tc.key, tc.from, tc.to, tc.amount)
			if !errors.Is(err, tc.want) { t.Fatalf("error = %v, want %v", err, tc.want) }
			after, _ := l.TotalBalance()
			if before != after { t.Fatalf("rejected transfer changed total: %d -> %d", before, after) }
			a, _ := l.Balance("a")
			b, _ := l.Balance("b")
			if a != 100 || b != 50 { t.Fatalf("rejected transfer changed balances: a=%d b=%d", a, b) }
		})
	}
}

func TestIdempotentReplayDoesNotTransferTwice(t *testing.T) {
	l := NewLedger()
	mustOpen(t, l, "a", 100)
	mustOpen(t, l, "b", 0)
	first, err := l.Transfer("same-key", "a", "b", 40)
	if err != nil { t.Fatal(err) }
	second, err := l.Transfer("same-key", "a", "b", 40)
	if err != nil { t.Fatal(err) }
	if first != second { t.Fatalf("replay receipt differs: first=%+v second=%+v", first, second) }
	a, _ := l.Balance("a")
	b, _ := l.Balance("b")
	if a != 60 || b != 40 { t.Fatalf("replay applied twice: a=%d b=%d", a, b) }
}

func TestIdempotencyKeyCannotBeReusedForDifferentRequest(t *testing.T) {
	l := NewLedger()
	mustOpen(t, l, "a", 100)
	mustOpen(t, l, "b", 0)
	mustOpen(t, l, "c", 0)
	if _, err := l.Transfer("key", "a", "b", 20); err != nil { t.Fatal(err) }
	_, err := l.Transfer("key", "a", "c", 20)
	if !errors.Is(err, ErrIdempotencyConflict) { t.Fatalf("error = %v, want %v", err, ErrIdempotencyConflict) }
	a, _ := l.Balance("a")
	b, _ := l.Balance("b")
	c, _ := l.Balance("c")
	if a != 80 || b != 20 || c != 0 { t.Fatalf("unexpected balances: a=%d b=%d c=%d", a, b, c) }
}

func TestDestinationOverflowIsRejectedAtomically(t *testing.T) {
	l := NewLedger()
	mustOpen(t, l, "a", 1)
	mustOpen(t, l, "b", math.MaxInt64)
	_, err := l.Transfer("overflow", "a", "b", 1)
	if !errors.Is(err, ErrBalanceOverflow) { t.Fatalf("error = %v, want %v", err, ErrBalanceOverflow) }
	a, _ := l.Balance("a")
	b, _ := l.Balance("b")
	if a != 1 || b != math.MaxInt64 { t.Fatalf("overflow changed balances: a=%d b=%d", a, b) }
}

func TestConcurrentTransfersCannotOverspend(t *testing.T) {
	l := NewLedger()
	mustOpen(t, l, "a", 100)
	mustOpen(t, l, "b", 0)
	var successes atomic.Int32
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := string(rune('A' + i))
			if _, err := l.Transfer(key, "a", "b", 10); err == nil {
				successes.Add(1)
			}
		}(i)
	}
	wg.Wait()
	if got := successes.Load(); got != 10 { t.Fatalf("successful transfers = %d, want 10", got) }
	a, _ := l.Balance("a")
	b, _ := l.Balance("b")
	if a != 0 || b != 100 { t.Fatalf("balances = (%d, %d), want (0, 100)", a, b) }
}
