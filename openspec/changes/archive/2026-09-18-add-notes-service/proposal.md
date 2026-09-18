# Proposal

## Why

There is no notes service yet. We need a small, well-structured microservice to create, list, edit, and delete notes, built with hexagonal architecture from the start so the domain logic stays testable and independent of Gin and MongoDB.

## What Changes

- New Go microservice (`notes-service`) implementing hexagonal architecture: `domain`, `application` (ports + use cases), and `adapters` (HTTP inbound, MongoDB outbound).
- CRUD HTTP API for notes: create, list, update, delete.
- Dependency injection wired with `uber-fx` across all layers (no manual wiring in `main.go` beyond `fx.New`).
- MongoDB persistence adapter implementing the repository port.
- Unit tests (`testify`) for domain rules and application use cases, mocking the repository port — no real MongoDB in unit tests.
- Project `CLAUDE.md` documenting the hexagonal architecture conventions (layer boundaries, where ports/adapters live, DI and testing rules) for this service.

## Capabilities

### New Capabilities
- `notes`: create, list, update, and delete notes, each with content, creation timestamp, and last-updated timestamp.

### Modified Capabilities
(none — greenfield project, no existing specs)

## Impact

- New Go codebase under this repo (no existing code to migrate).
- New MongoDB collection for notes.
- New HTTP API surface (create/list/update/delete endpoints).
- New `CLAUDE.md` at the project root capturing architecture conventions for future work on this service.
