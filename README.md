# Institutional Settlement Network

A learning-oriented prototype exploring institutional settlement using the Cosmos ecosystem. The project is developed alongside a separate learning lab so that concepts are understood and documented before corresponding implementation work begins.

## Development rule

**GPT teaches first. Codex implements only after the relevant concept has been learned and the design gate explicitly approved.** Every concept used in this product must have corresponding learning documentation in the separate `cosmos-learning-lab` repository.

## Intended prototype scope

- Simulated institutions, accounts, balances, and settlement transactions
- Explicit authorization and accounting/state invariants
- Cosmos-based authoritative settlement state, subject to approved architecture and compatibility decisions
- PostgreSQL read model/indexer that is derived and rebuildable
- Separate settlement and reconciliation workflows
- APIs and potential React/TypeScript consoles for bank participants, treasury, network operations, and exploration
- Tests for success, rejection, duplicate requests, replay/rebuild, and failure/recovery paths

## Non-goals and cautions

This repository is a learning prototype, not a production-ready financial network. No real funds, credentials, or production data should be used. Production use would require independent security, operational, legal, regulatory, privacy, and governance reviews.

IBC behavior, optional CosmWasm support, dependency versions, module boundaries, and deployment topology are not finalized until compatibility and architecture gates are passed.

## Repository map

- `AGENTS.md` — implementation and quality rules
- `docs/ARCHITECTURE_PROPOSAL.md` — provisional architecture and invariants to resolve
- `docs/decisions/` — architecture decision records
- `docs/api/` — approved API contracts
- `app/` — Cosmos application, only after approval
- `indexer/` — derived read model and reconciliation tooling
- `web/` — participant and operations interfaces, if approved
- `tests/` — unit, integration, end-to-end, and failure tests
- `ops/` — local development and operational documentation

Directories are added when substantive content is introduced.

## Status

Foundation documentation only. Application code has not started. See the architecture proposal and decision records before implementation.
