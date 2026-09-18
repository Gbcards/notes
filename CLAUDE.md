# notes-service

Go microservice for creating, listing, editing, and deleting notes. Built with hexagonal architecture (ports & adapters), `uber-fx` for dependency injection, and MongoDB for persistence.

## Layout

```
cmd/api/main.go                     -- fx.New(platform.Modules()...).Run() -- nothing else
internal/domain/note/                -- Note entity + domain errors, NO external imports
internal/application/note/           -- ports.go (Repository port) + use cases, unit tests
internal/adapters/http/              -- Gin handlers + DTOs (adapter IN)
internal/adapters/mongo/             -- Mongo repository implementation (adapter OUT)
internal/platform/                   -- fx modules wiring providers per layer, Mongo/HTTP server lifecycle
```

Request flow: HTTP request -> Gin handler (adapter IN) -> use case (application layer) -> `Note` (domain) -> `Repository` port -> Mongo adapter (adapter OUT) -> MongoDB.

## Layering rules

- **`internal/domain`** imports nothing outside the Go standard library. No Gin, no Mongo driver, no `uber-fx`. Business rules (e.g. rejecting empty content) live here as methods on `Note`.
- **`internal/application`** defines the outbound port (`Repository` interface in `ports.go`) and the use cases that depend on it. It never imports Gin, `uber-fx`, or the Mongo driver — only `internal/domain`.
- **`internal/adapters/http`** and **`internal/adapters/mongo`** are the only packages allowed to import Gin and the Mongo driver, respectively. The HTTP adapter defines its own small inbound-port interfaces (`NoteCreator`, `NoteLister`, `NoteUpdater`, `NoteDeleter`) that a use case's concrete type satisfies structurally.
- **`internal/platform`** is the only place that wires concrete adapters to ports. It is the only package that imports both an adapter package and `uber-fx`.

If a change to `internal/domain` or `internal/application` would require importing `gin-gonic/gin`, `go.mongodb.org/mongo-driver`, or `go.uber.org/fx`, that is a signal the logic belongs in an adapter instead.

## Dependency injection (`uber-fx`)

- One `fx.Module` per layer in `internal/platform/fx_modules.go` (`mongoModule`, `applicationModule`, `httpModule`), each with its own `fx.Provide(...)` list.
- A concrete adapter is bound to the interface it satisfies with `fx.Annotate(ctor, fx.As(new(Interface)))` (see `mongoModule`'s binding of `*mongoadapter.NoteRepository` to `application.Repository`), or with a small conversion function when the interface lives in a different package than the one importing it back (see `httpModule`'s `func(uc *application.CreateNoteUseCase) httpadapter.NoteCreator { return uc }`).
- `cmd/api/main.go` never constructs a dependency by hand — it only calls `fx.New(platform.Modules()...).Run()`. Anything that must run even though nothing else depends on it (e.g. the HTTP server) is forced into existence with `fx.Invoke` inside `platform.Modules()`, not in `main.go`.
- Long-running resources (the Mongo client, the HTTP server) attach `OnStart`/`OnStop` hooks via `fx.Lifecycle` in their constructor, in `internal/platform`.

## Testing

- **Domain**: plain `testing` unit tests on `Note`'s rules — no external dependencies needed.
- **Application**: unit tests per use case, mocking the `Repository` port with `testify/mock`. Never touch a real MongoDB in these tests.
- **HTTP adapter**: unit tests per handler using `net/http/httptest`, with hand-written stubs for `NoteCreator`/`NoteLister`/`NoteUpdater`/`NoteDeleter` — no real use cases or Mongo involved.
- **Mongo adapter**: unit tests are limited to pure logic (e.g. ID mapping in `model_test.go`). Integration tests against a real MongoDB are out of scope for this service so far.

## Note identifier

A note's `ID` is an opaque string everywhere outside `internal/adapters/mongo`. The Mongo adapter is the only place that knows it is a hex-encoded `primitive.ObjectID`; it converts at the boundary (`parseID` / `.Hex()`).
