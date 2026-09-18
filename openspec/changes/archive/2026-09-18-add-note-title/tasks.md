# Tasks

## 1. Domain

- [x] 1.1 Add `Title string` to `Note` (`internal/domain/note/note.go`); change `New(title, content string, now time.Time) (*Note, error)` to set it (still validates only `content`)
- [x] 1.2 Rename `UpdateContent` to `Update(title, content string, now time.Time) error`, setting both `Title` and `Content` (still validates only `content`)
- [x] 1.3 Update `internal/domain/note/note_test.go` for the new signatures and title assertions; verify with `go test ./internal/domain/...`

## 2. Application

- [x] 2.1 Change `CreateNoteUseCase.Execute` to `Execute(ctx context.Context, title, content string) (*domain.Note, error)`, passing `title` into `domain.New` (`internal/application/note/create_note.go`)
- [x] 2.2 Change `UpdateNoteUseCase.Execute` to `Execute(ctx context.Context, id, title, content string) (*domain.Note, error)`, passing both into `Note.Update` (`internal/application/note/update_note.go`)
- [x] 2.3 Update `create_note_test.go` and `update_note_test.go` call sites and mock expectations for the new parameter; verify with `go test ./internal/application/...`

## 3. HTTP adapter

- [x] 3.1 Add `Title` to `createNoteRequest`, `updateNoteRequest`, and `noteResponse` (`json:"title"`) and populate it in `toNoteResponse` (`internal/adapters/http/dto.go`)
- [x] 3.2 Update `NoteCreator.Execute` and `NoteUpdater.Execute` interface signatures and the `Create`/`Update` handler methods to pass `req.Title` through (`internal/adapters/http/note_handler.go`)
- [x] 3.3 Update `note_handler_test.go` stubs and request/response fixtures for the new field; verify with `go test ./internal/adapters/http/...`

## 4. Mongo adapter

- [x] 4.1 Add `Title string` with `bson:"title"` to `noteDocument`, and set it in `newDocument` and `toDomain` (`internal/adapters/mongo/model.go`)
- [x] 4.2 Include `"title"` in the `$set` of `NoteRepository.Update` (`internal/adapters/mongo/note_repository.go`)
- [x] 4.3 Update `model_test.go` for the new field, including a case asserting a document without a stored `title` key decodes to an empty title; verify with `go test ./internal/adapters/mongo/...`

## 5. Verification

- [x] 5.1 Run `go build ./...` to confirm every call site (including `internal/platform/fx_modules.go` bindings) compiles against the new signatures
- [x] 5.2 Run `go test ./...` and confirm all packages pass
