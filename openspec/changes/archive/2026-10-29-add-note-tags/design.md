# Design

## Context

See proposal.md - Why. `Note.New` and `Note.Update` (`internal/domain/note/note.go`) already enforce two rules (`ErrEmptyContent`, `ErrTitleTooLong`). Each validates everything before mutating anything. `writeError` in `internal/adapters/http/note_handler.go` maps each domain error to `400`. The Mongo adapter persists through `newDocument` and `toDomain` (`model.go`). The repository's `Update` (`note_repository.go`) builds its `$set` by hand, field by field, so a new field is not persisted on update unless it is added there explicitly.

`limit-note-title-length` is complete but not yet archived. Its delta MODIFIES the *Create Note* and *Update Note* requirements.

## Goals / Non-Goals

**Goals:**
- Keep the whole tag rule (normalization, deduplication, limits) in `internal/domain/note`, with only standard library imports.
- Guarantee `"tags": []` (never `null`) in every response, including documents that have no `tags` field.

**Non-Goals:**
- Unicode normalization (NFC/NFD) of tags. This is consistent with the title decision in `limit-note-title-length`.
- Collapsing internal whitespace (`"machine  learning"` is kept as-is apart from trimming and lowercasing).
- An index on `tags`, since filtering by tag is out of scope.
- Changes to `internal/platform`. See Decisions.

## Decisions

**One unexported normalizer in the domain.** `normalizeTags(raw []string) ([]string, error)` lives in `note.go`. `New` and `Update` call it after the existing `ErrEmptyContent` and `ErrTitleTooLong` checks and before any field is assigned, so a rejected update leaves the note unchanged. It processes tags in this order:
1. `strings.TrimSpace`, then `strings.ToLower`.
2. Skip the tag if it is empty.
3. If `utf8.RuneCountInString(tag) > 20`, return `ErrTagTooLong`.
4. Skip the tag if it has already been seen (tracked with a `map[string]struct{}`). Otherwise append it, which preserves first-occurrence order.
5. After the loop, if the result has more than 5 tags, return `ErrTooManyTags`.

The loop never stops early on the count, so when a request breaks both limits, `ErrTagTooLong` wins. The result is always a newly allocated, non-nil slice (`make([]string, 0, len(raw))`) that never aliases the caller's input.

**Where nil is allowed.** A `Note` built or updated through `New`/`Update` always has non-nil `Tags`. Only `toDomain` can produce `nil` (legacy documents without the field), and only `toNoteResponse` converts `nil` to `[]`. Because every write goes through `New` or `Update`, `newDocument` and the update document always receive non-nil tags and add no guard of their own. The 5-tag count runs after deduplication, because counting raw entries would turn duplicates into errors, and the story says duplicates never fail. The 20-character limit is measured after trimming, so surrounding padding never causes a rejection. The constants `maxTagRunes = 20` and `maxTags = 5` sit next to `maxTitleRunes`.

**Positional `tags []string` parameter instead of a command struct.** The signatures become `New(title, content string, tags []string, now time.Time)`, `Update(title, content string, tags []string, now time.Time)`, `CreateNoteUseCase.Execute(ctx, title, content, tags)` and `UpdateNoteUseCase.Execute(ctx, id, title, content, tags)`, and `NoteCreator` and `NoteUpdater` change to match. A struct such as `NoteInput` would scale better, but it would be a refactor across every layer for a single added field. That can be revisited if more fields arrive.

**Two new domain errors, mapped like the existing ones.** `ErrTagTooLong = errors.New("note tag must not exceed 20 characters")` and `ErrTooManyTags = errors.New("note must not have more than 5 tags")` go in the `var` block in `errors.go`. Like `ErrTitleTooLong`, each message includes its limit. Each also gets its own `case` in `writeError` returning `400`.

**Type validation stays at the HTTP boundary.** The request DTOs declare `Tags []string \`json:"tags"\``. `ShouldBindJSON` already fails with `400` on `"go"` or `[1, 2]`, so no new code is needed, only tests. `"tags": null` and an omitted field both decode to `nil`, which normalizes to `[]`. A `null` inside the array (`["go", null]`) decodes to `""` and is discarded, which is acceptable.

**The "never null" guarantee lives in `toNoteResponse`.** Legacy Mongo documents decode to a `nil` slice through `toDomain`, which bypasses the domain constructor, so the domain cannot guarantee non-nil on its own. `toNoteResponse` converts `nil` to `[]string{}`. That is the single place where the JSON contract is produced, and it covers create, update, list, and legacy data. Alternative considered: also coercing in `toDomain`. Rejected to keep one owner of the guarantee.

**Mongo document field without `omitempty`.** The document gets `Tags []string \`bson:"tags"\``. `newDocument` and `toDomain` copy it. New notes therefore always store an array, while legacy documents simply lack the field. No migration is needed.

**Move the update document into a pure helper so it can be tested.** The hand-built `$set` in `note_repository.go` `Update` moves into `updateDocument(n *domain.Note) bson.M` in `model.go`, which now also includes `"tags"`. `Update` calls it. A forgotten `$set` field is the most likely way this change breaks, and integration tests are out of scope, so this lets `model_test.go` catch it with a pure unit test. Alternative considered: leave the `$set` inline and rely on review. Rejected because that leaves acceptance criterion 7 untested.

**Spec delta uses ADDED plus a MODIFIED List, not MODIFIED Create or Update.** All tag behavior sits in a new *Note Tags* requirement, so this change does not touch the *Create Note* and *Update Note* blocks that `limit-note-title-length` modifies. That avoids one change's archive overwriting the other's.

**No platform changes.** The `httpModule` conversion functions (`func(uc *application.CreateNoteUseCase) httpadapter.NoteCreator { return uc }`) keep compiling as long as the use case and inbound port signatures change together. No new providers are needed.

## Testing

The signature change breaks existing call sites in `note_test.go` (every `New` and `Update` call), `create_note_test.go`, `update_note_test.go`, and the `stubCreator` and `stubUpdater` in `note_handler_test.go`. All of them are updated, and the mock `Repository` is unchanged because the port does not change.

- **Domain**: covers every rule in the spec. That means duplicates, normalization, empty tags, first-occurrence order, a tag of exactly 20 runes (multi-byte included), a 21-rune tag, 6 tags that dedupe to 5 being accepted, 6 distinct tags being rejected, `ErrTagTooLong` winning over `ErrTooManyTags`, nil input producing non-nil output, an output slice that does not alias the input, and a rejected `Update` leaving `Tags` unchanged.
- **Application**: tags are passed through. `mock.MatchedBy` asserts that `repo.Create` and `repo.Update` receive the normalized tags, and `AssertNotCalled` confirms nothing is persisted when tag validation fails.
- **HTTP**: the stubs record the tags they receive, so tests can assert the handler forwards `req.Tags`. Tests decode the response body to assert `"tags":[]` when the stub returns nil `Tags` (on create, update, and list) and assert that tags are echoed back. Other tests check that `"tags": "go"` and `"tags": [1, 2]` return 400 without calling the stub, and that `ErrTagTooLong` and `ErrTooManyTags` map to 400.
- **Mongo**: `newDocument` and `toDomain` copy tags. A BSON document with no `tags` key decodes to nil (same pattern as the missing-title test). `updateDocument` includes `title`, `content`, `tags`, and `updated_at`.

## Risks / Trade-offs

- [Existing clients that send PUT without `tags` will wipe a note's tags] -> This is intended (full-replacement semantics, matching `title`). It is documented in the spec scenario *Update without tags clears them* and should be communicated to API consumers.
- [The signature change breaks compilation of downstream layers until each one is updated] -> `tasks.md` groups run sequentially and each group verifies with a layer-scoped `go test`. `go build ./...` is only expected to pass at the Verification step.
- [`strings.ToLower` is a simple case mapping, not full Unicode case folding (e.g. `ß`)] -> Accepted. It does not matter for realistic tag input.

## Migration Plan

None. Existing documents without `tags` are read as an empty list. Rollback means redeploying the previous version, which ignores the extra `tags` field in documents written meanwhile.
