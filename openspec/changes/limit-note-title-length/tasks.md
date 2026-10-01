# Tasks

## 1. Domain

- [x] 1.1 Add `ErrTitleTooLong` to `internal/domain/note/errors.go`, next to `ErrEmptyContent`, and verify `go build ./...` succeeds
- [x] 1.2 In `internal/domain/note/note.go`, reject a title whose `utf8.RuneCountInString` exceeds 200 in both `New` and `Update`, returning `ErrTitleTooLong`; verify with domain unit tests covering a 200-character title (accepted), a 201-character title (rejected), and a title containing multi-byte characters (e.g. `ñ`, `á`) at exactly 200 runes (accepted)

## 2. HTTP Adapter

- [x] 2.1 Add a `case errors.Is(err, domain.ErrTitleTooLong)` to `writeError` in `internal/adapters/http/note_handler.go` mapping to `http.StatusBadRequest`; verify with `httptest`-based handler unit tests for both the create and update handlers submitting an over-limit title

## 3. Verification

- [x] 3.1 Run `go build ./...` and `go test ./...` for the whole repo and confirm everything passes
