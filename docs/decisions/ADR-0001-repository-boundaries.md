# ADR-0001: Keep learning and product repositories independent

- **Status:** Accepted
- **Date:** 2026-10-24
- **Decision owners:** Project learner/maintainer

## Context

The project has two complementary goals: deeply learn the Cosmos stack and build a settlement-network prototype. Combining notes, labs, and product implementation in one repository would blur the learning record and product boundaries.

## Decision

Maintain two independent public repositories:

1. `hari-devloper/cosmos-learning-lab` contains the curriculum, concept handbook, source-code explorations, labs, failure experiments, interview preparation, and mastery evidence.
2. `hari-devloper/institutional-settlement-network` contains product requirements, architecture decisions, settlement domain design, application code, tests, read-model/indexer, interfaces, security, and operations documentation.

The repositories must not share source directories, Git submodules, or build dependencies. The product repository may reference learning concepts by stable IDs and links.

## Consequences

- Learning notes remain independently reviewable and useful for interview preparation.
- Product changes must keep corresponding learning documentation current.
- Some concepts will be cross-referenced rather than duplicated.
- Each repository has its own tests, dependency declarations, and release history.
- Changes to architecture and dependencies must pass explicit review gates before implementation.

## Revisit when

Revisit only if the project goals change materially or a concrete maintenance issue demonstrates that this boundary no longer serves them.
