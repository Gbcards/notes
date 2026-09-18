# Tasks

## 1. Project Setup

- [x] 1.1 Initialize `go.mod` for the `notes-service` module and verify `go build ./...` succeeds with an empty `main.go`
- [x] 1.2 Add dependencies (`gin-gonic/gin`, `go.mongodb.org/mongo-driver`, `go.uber.org/fx`, `stretchr/testify`) and verify `go mod tidy` runs clean
- [x] 1.3 Create the package skeleton from design.md (`internal/domain/note`, `internal/application/note`, `internal/adapters/http`, `internal/adapters/mongo`, `internal/platform`, `cmd/api`) and verify each directory exists with a placeholder file

## 2. Domain Layer

- [x] 2.1 Implement the `Note` entity (ID, Content, CreatedAt, UpdatedAt) in `internal/domain/note/note.go` with a constructor that rejects empty content, and verify unit tests cover accepting non-empty content and rejecting empty content (spec: Create Note)
- [x] 2.2 Define domain sentinel errors (`ErrNotFound`, `ErrEmptyContent`) in `internal/domain/note/errors.go` and verify they are exported and referenced by the tests in 2.1

## 3. Application Layer (Ports and Use Cases)

- [x] 3.1 Define the `Repository` port interface (Create, List, Update, Delete) in `internal/application/note/ports.go` and verify it compiles with no import of Gin, `uber-fx`, or the Mongo driver
- [x] 3.2 Implement `CreateNote` use case and verify unit tests (mocking `Repository` with `testify/mock`) cover successful creation and rejection of empty content (spec: Create Note)
- [x] 3.3 Implement `ListNotes` use case and verify unit tests cover returning all notes and returning an empty list when none exist (spec: List Notes)
- [x] 3.4 Implement `UpdateNote` use case and verify unit tests cover successful update (content changes, `UpdatedAt` changes, `CreatedAt` unchanged), update of a nonexistent ID returning `ErrNotFound`, and rejection of empty content (spec: Update Note)
- [x] 3.5 Implement `DeleteNote` use case and verify unit tests cover successful deletion and deletion of a nonexistent ID returning `ErrNotFound` (spec: Delete Note)

## 4. MongoDB Adapter (Outbound)

- [x] 4.1 Implement the Mongo document model and ID mapping (string <-> `primitive.ObjectID`) in `internal/adapters/mongo/model.go`, verified by a unit test round-tripping a valid hex ID and rejecting a malformed one
- [x] 4.2 Implement `internal/adapters/mongo/note_repository.go` satisfying the `Repository` port, translating `mongo.ErrNoDocuments` into `domain/note.ErrNotFound`, and verify it compiles against the port interface (`var _ application.Repository = (*NoteRepository)(nil)`)

## 5. HTTP Adapter (Inbound)

- [x] 5.1 Define request/response DTOs in `internal/adapters/http/dto.go` matching the spec's fields (content, id, createdAt, updatedAt)
- [x] 5.2 Implement Gin handlers for create/list/update/delete in `internal/adapters/http/note_handler.go`, mapping `ErrEmptyContent` to 400 and `ErrNotFound` to 404, and verify unit tests per handler cover the success path and each error mapping
- [x] 5.3 Wire the routes in `internal/adapters/http/router.go` and verify a manual request against each of the 4 routes (via `httptest`) returns the expected status code

## 6. Dependency Injection Wiring

- [x] 6.1 Add an `fx.Module` per layer (`internal/platform/fx_modules.go`) providing the Mongo repository bound to the `Repository` interface, the use cases, and the HTTP handlers/router, and verify `cmd/api/main.go` only calls `fx.New(platform.Modules()...)` with no manual constructor calls
- [x] 6.2 Add Mongo connection config/lifecycle (`internal/platform/config.go`, using `fx.Lifecycle` for connect/disconnect) and verify the app starts and stops cleanly against a local MongoDB (manual run)

## 7. Documentation

- [x] 7.1 Write the project `CLAUDE.md` documenting: the layering rule (domain has no external imports), where ports/adapters live, the `uber-fx` module-per-layer DI convention, and the testing convention (mock ports, no real MongoDB in unit tests), and verify it matches the structure agreed in design.md

## 8. Verification

- [x] 8.1 Run `go test ./...` and verify all unit tests pass (domain, application, adapters model mapping)
- [x] 8.2 Manually exercise all 4 endpoints (create, list, update, delete) against a local MongoDB and verify responses match the scenarios in `specs/notes/spec.md`
