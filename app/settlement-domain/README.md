# Settlement domain invariant lab

This is the first executable domain-model milestone for the Institutional Settlement Network. It uses only the Go standard library so accounting rules can be tested before Cosmos SDK wiring.

## Run

    go test ./...
    go test -race ./...

## Invariants covered

- A transfer amount must be positive and the source and destination must differ.
- Both accounts must exist; the source must have sufficient funds.
- A successful transfer debits and credits the same amount.
- Rejected transfers do not change balances.
- Total balances are conserved by transfers, excluding explicit account opening/genesis-like setup.
- Repeating an identical idempotency key returns the original receipt; conflicting reuse is rejected.
- Integer overflow is rejected before mutation.
- Concurrent requests cannot spend the same funds twice within this process.

## Important limitations

This is **not yet a Cosmos SDK module** and is not durable or distributed. The mutex protects only one process. A real chain must enforce invariants in deterministic application state, define authorization and denomination rules, persist transaction identity according to consensus semantics, and test failure/commit behavior against the pinned SDK. OpenAccount is a test/genesis-like setup operation, not a public issuance API.

Learning reference: COSMOS-047 in the separate cosmos-learning-lab repository.
