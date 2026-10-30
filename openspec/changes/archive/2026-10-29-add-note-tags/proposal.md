# Proposal

## Why

Users have no way to classify or organize their notes (story WR-1234). Adding optional tags lets a client label each note. Normalizing and deduplicating tags means a note never carries repeated or cosmetically different labels (`Go`, ` go `, `GO`).

## What Changes

- Notes gain an optional `tags` field (a list of strings) that clients submit when creating or updating a note and that the API returns on create, update, and list.
- Tags are normalized: leading and trailing whitespace is trimmed and the tag is lowercased. Tags left empty are discarded, and duplicates are removed, each tag keeping the position of its first occurrence. Duplicates never cause an error.
- A tag longer than 20 characters (after trimming, counted in Unicode characters) is rejected with `400 Bad Request`.
- A note with more than 5 tags (counted after normalization and duplicate removal) is rejected with `400 Bad Request`.
- `PUT /notes/{id}` replaces the whole tag list. Omitting `tags` (or sending `null`) leaves the note with no tags, the same way `title` already behaves on update.
- A `tags` value that is not a list of strings (e.g. `"go"` or `[1, 2]`) is rejected with `400 Bad Request`.
- A response always contains `"tags": []` rather than `null`, including for notes stored before this change that have no `tags` field in MongoDB.
- The 20-character and 5-tag limits are additions to WR-1234, agreed during exploration. They need to be reflected in the story's acceptance criteria.

## Capabilities

### New Capabilities

(none)

### Modified Capabilities

- `notes`: adds a *Note Tags* requirement (submission, normalization, deduplication, limits, replacement on update, empty-list default), and *List Notes* now returns each note's tags.

## Impact

- `internal/domain/note`: `Note.Tags`, a tag normalization rule applied in `New` and `Update`, and the new errors `ErrTagTooLong` and `ErrTooManyTags`.
- `internal/application/note`: the create and update use cases take tags and pass them to the domain.
- `internal/adapters/http`: a `tags` field in the request and response DTOs, the inbound port signatures (`NoteCreator`, `NoteUpdater`), and `writeError` mappings for the new errors.
- `internal/adapters/mongo`: a `tags` field in the document, and the update `$set` in `note_repository.go`.
- API: additive and backward compatible for create and list. On update, clients that omit `tags` will clear them.
- No new dependencies, and no data migration.

## Out of Scope

- Filtering or searching notes by tag.
