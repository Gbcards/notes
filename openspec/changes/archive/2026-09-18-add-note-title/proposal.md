# Proposal

## Why

Notes currently have no way to summarize themselves beyond their full content. Clients need a short, optional title to identify a note without reading its content.

## What Changes

- Add an optional `title` field to notes, settable when creating a note and editable when updating it.
- `title` has no validation rule (it may be empty or omitted); `content` keeps its existing "must not be empty" rule unchanged.
- Return `title` alongside `content` and the timestamps in note responses (create, list, update).
- Notes created before this change have no stored title and are treated as having an empty title.

## Capabilities

### New Capabilities
(none)

### Modified Capabilities
- `notes`: Create Note and Update Note requirements gain an optional `title` field; List Note (and the response shape used by Create/Update) includes `title` in the returned note.

## Impact

- `internal/domain/note`: `Note` struct gains `Title`; `New` and the content-update method take a title parameter.
- `internal/application/note`: `CreateNoteUseCase.Execute` and `UpdateNoteUseCase.Execute` gain a `title` parameter.
- `internal/adapters/http`: `createNoteRequest`, `updateNoteRequest`, and `noteResponse` gain a `title` field; `NoteCreator`/`NoteUpdater` port interfaces change signature.
- `internal/adapters/mongo`: `noteDocument` gains a `title` bson field; `Update` persists it. No data migration needed — existing documents decode with an empty title.
- Existing unit tests for domain, application, HTTP adapter, and mongo model need updating for the new field/signatures.
