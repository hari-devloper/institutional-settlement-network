# Agent and contributor instructions

## Learning and approval gates

- The learner studies the concept in GPT before implementation is delegated to Codex.
- Link every meaningful product concept to a stable ID in the separate learning repository (for example, `COSMOS-001`).
- Do not start implementation while relevant compatibility, domain, or architecture decisions remain unapproved.
- State assumptions, acceptance criteria, failure cases, and verification steps before substantive changes.
- Do not describe the project as production-ready.

## Architecture boundaries

- Treat on-chain state as authoritative for settlement only after the application architecture is approved.
- Treat PostgreSQL/indexer data as a derived, rebuildable read model; never silently make it authoritative for settlement.
- Keep settlement and reconciliation as distinct workflows with explicit state transitions and audit trails.
- Model duplicate requests, idempotency, rejected transactions, timeouts, acknowledgements, refunds, recovery, and reconciliation explicitly where applicable.
- Do not assume that a timeout alone proves a remote chain performed no action.
- Keep the learning repository independent: no shared source directories, submodules, or cross-repository build dependencies.

## Engineering quality

- Define domain invariants before implementation and encode them as tests.
- Prefer small changes with focused tests and reviewable commits.
- Never fabricate test results, benchmarks, source paths, compatibility claims, or security assurances.
- Record exact dependency versions and justify compatibility using official release metadata and module files.
- Keep secrets, private keys, credentials, real customer data, and production endpoints out of the repository.
- Record actual test commands and results; distinguish expected output from observed output.
- Update docs and diagrams with behavior changes.

## Security and release posture

This is an educational prototype. Before any production use, require independent security review, operational readiness assessment, legal/regulatory review, privacy review, governance decisions, and explicit release approval.
