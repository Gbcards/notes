# Design

## Context

Greenfield project — no existing code. See proposal.md - Why. The service must keep domain logic free of framework/infrastructure dependencies (hexagonal architecture) while using Go, Gin, MongoDB, and `uber-fx` for dependency injection, per the chosen stack.

## Goals / Non-Goals

**Goals:**
- Domain and application layers have zero import dependency on Gin, the MongoDB driver, or `uber-fx`.
- Every cross-layer dependency is expressed as a Go interface (port) and satisfied by an adapter, wired through `uber-fx` providers — no dependency is constructed and passed by hand outside of `fx`.
- Application-layer use cases are unit-testable against a mocked repository port, with no real MongoDB involved.

**Non-Goals:**
- Authentication/authorization on the API (not requested).
- Integration tests against a real MongoDB instance (unit tests only, per the confirmed scope).
- Multi-tenant or multi-user note ownership — notes have no owner concept.

## Decisions

### Layering and package layout
```
cmd/api/main.go                     -- fx.New(...) entrypoint
internal/domain/note/                -- Note entity + domain errors, no external imports
internal/application/note/           -- ports.go (repository port), use cases, unit tests
internal/adapters/http/              -- Gin handlers + DTOs (adapter IN)
internal/adapters/mongo/             -- Mongo repository implementation (adapter OUT)
internal/platform/                   -- fx modules wiring providers per layer, config
```
Rationale: keeps the dependency rule enforceable by package boundaries (Go tooling flags an accidental import from `domain` into `adapters`), and mirrors the diagram already agreed with the user (Gin handler -> use case -> domain -> repository port -> Mongo adapter).

### Dependency injection with uber-fx
Each layer exposes an `fx.Module` (e.g. `mongoModule`, `applicationModule`, `httpModule`) that provides constructors keyed by interface type (`fx.Provide(NewMongoNoteRepository)` bound to the `application/note.Repository` interface). `main.go` only calls `fx.New(platform.Modules()...)`.
Alternative considered: manual constructor wiring in `main.go` (idiomatic minimal Go). Rejected because the user explicitly wants `uber-fx`.

### Note identifier
Notes are stored with MongoDB's native `ObjectID`, but the domain and application layers treat the identifier as an opaque `string` (hex-encoded) — the Mongo adapter converts to/from `primitive.ObjectID` at the boundary. This keeps `domain/note` free of the Mongo driver import.
Alternative considered: application-generated UUID. Rejected as unnecessary extra complexity; nothing in the requirements needs client-supplied or globally-unique-across-systems IDs.

### Error mapping
Domain/application errors are typed sentinel values in `domain/note/errors.go` (e.g. `ErrNotFound`, `ErrEmptyContent`). The HTTP adapter is the only layer that maps these to status codes (`ErrNotFound` -> 404, `ErrEmptyContent` -> 400, anything else -> 500). The Mongo adapter translates `mongo.ErrNoDocuments` into `domain/note.ErrNotFound` before returning, so the application layer never sees a Mongo-specific error type.

### Testing strategy
- `domain/note`: pure unit tests on validation rules (e.g. rejecting empty content).
- `application/note`: unit tests per use case, mocking the repository port with `testify/mock`; no real MongoDB.
- `adapters/*`: left out of scope for this change (non-goal) — could be added later as a separate change if integration coverage is wanted.

### CLAUDE.md
Written at the project root during implementation (tasks.md) to document: the layering rule above, where ports/adapters live, the DI convention (module-per-layer, provide by interface), and the testing convention (mock ports, no real MongoDB in unit tests). This is documentation output, not a spec — it doesn't change externally observable behavior.

## Risks / Trade-offs

- [`uber-fx`'s reflection-based wiring can obscure "who provides what" for newcomers] -> Mitigation: keep one `fx.Module` per layer with an explicit, small provider list; document the convention in `CLAUDE.md`.
- [Exposing Mongo's `ObjectID` format indirectly through a hex-string ID couples the API's ID format to Mongo] -> Mitigation: acceptable trade-off given MongoDB is the only planned store; revisit if a store change is ever proposed.
- [No integration tests against real MongoDB means the adapter's query/update logic is unverified by automated tests] -> Mitigation: explicit non-goal for this change; flag as a follow-up if defects appear in the adapter.
