# Tasks

## 1. Domain

- [x] 1.1 Add `ErrTagTooLong` ("note tag must not exceed 20 characters") and `ErrTooManyTags` ("note must not have more than 5 tags") to the `var` block in `internal/domain/note/errors.go`, and verify `go build ./internal/domain/...` succeeds
- [x] 1.2 In `internal/domain/note/note.go`, add `Tags []string` to `Note`, add the constants `maxTagRunes = 20` and `maxTags = 5` next to `maxTitleRunes`, and implement the unexported `normalizeTags` exactly as described in design.md (trim, lower, drop empty, rune-length check, dedupe in first-occurrence order, count after the loop, non-nil and non-aliasing result). Verify with `go build ./internal/domain/...`
- [x] 1.3 Change `New(title, content string, tags []string, now time.Time)` and `Update(title, content string, tags []string, now time.Time)` to call `normalizeTags` after the content and title checks and before assigning any field. Update every existing call in `note_test.go`, and verify `go test ./internal/domain/...` passes with the new tests listed under Domain in design.md's Testing section

## 2. Application

- [x] 2.1 Add a `tags []string` parameter to `CreateNoteUseCase.Execute` (`create_note.go`) and `UpdateNoteUseCase.Execute` (`update_note.go`) and pass it to `domain.New` and `n.Update`. Update the existing calls in `create_note_test.go` and `update_note_test.go`, and add tests that use `mock.MatchedBy` to assert normalized tags reach `repo.Create`/`repo.Update` and `AssertNotCalled` to assert nothing is persisted on `ErrTagTooLong`/`ErrTooManyTags`. Verify with `go test ./internal/application/...`

## 3. HTTP Adapter

- [x] 3.1 In `internal/adapters/http/dto.go`, add `Tags []string \`json:"tags"\`` to `createNoteRequest`, `updateNoteRequest`, and `noteResponse`, and make `toNoteResponse` return `[]string{}` when `n.Tags` is nil. Verify with `go build ./internal/adapters/http/...`
- [x] 3.2 In `note_handler.go`, add `tags []string` to the `NoteCreator` and `NoteUpdater` `Execute` signatures, pass `req.Tags` from `Create` and `Update`, and add `writeError` cases mapping `ErrTagTooLong` and `ErrTooManyTags` to `400`. Verify with `go build ./internal/adapters/http/...`
- [x] 3.3 In `note_handler_test.go`, update `stubCreator` and `stubUpdater` to the new signatures and have them record the tags they receive. Add the tests listed under HTTP in design.md's Testing section (tags forwarded, tags echoed, `"tags":[]` for nil on create/update/list, 400 for `"go"` and `[1, 2]` without calling the stub, 400 for both new errors), and verify with `go test ./internal/adapters/http/...`

## 4. Mongo Adapter

- [x] 4.1 In `internal/adapters/mongo/model.go`, add `Tags []string \`bson:"tags"\`` to `noteDocument`, copy it in `newDocument` and `toDomain`, and add `updateDocument(n *domain.Note) bson.M` returning the `$set` with `title`, `content`, `tags`, and `updated_at`. Switch `NoteRepository.Update` in `note_repository.go` to use it. Verify with `go build ./internal/adapters/mongo/...`
- [x] 4.2 In `model_test.go`, add tests for tags in `newDocument` and `toDomain`, for a BSON document with no `tags` key decoding to nil tags without error, and for `updateDocument` containing all four fields with the note's tags. Verify with `go test ./internal/adapters/mongo/...`

## 5. Verification

- [x] 5.1 Run `go vet ./...`, `go build ./...` and `go test ./...` for the whole repo, including `internal/platform` with no changes to it, and confirm everything passes
