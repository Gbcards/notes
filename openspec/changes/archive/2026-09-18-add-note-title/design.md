# Design

## Context

`Note` (`internal/domain/note/note.go`) currently has `ID`, `Content`, `CreatedAt`, `UpdatedAt`. `Content` is required; `domain.New` and the update method reject empty content via `ErrEmptyContent`. The HTTP DTOs (`internal/adapters/http/dto.go`) and the Mongo document (`internal/adapters/mongo/model.go`) mirror those same fields with no partial-update semantics: a create or update request always carries the full set of fields the client wants stored (there's no way to distinguish "field omitted" from "field sent empty" since DTOs use plain `string`, not pointers). See proposal.md - Why/What Changes for the motivation.

## Goals / Non-Goals

**Goals:**
- Add `Title` end-to-end (domain, application, HTTP, Mongo) following the existing full-replacement request shape.
- Keep `Title` unvalidated (any string, including empty, is valid) — no new domain error.
- Preserve backward compatibility for notes persisted before this change (no migration script).

**Non-Goals:**
- No partial-update ("only change title, leave content untouched") semantics — that's a bigger behavior change than this proposal asks for, and the existing `Content` field doesn't support it either.
- No title-based search/filtering — out of scope, not requested.
- No length limit or format validation on `Title`.

## Decisions

- **Parameter order is `title` before `content`** in every changed signature (`domain.New(title, content string, now time.Time)`, `CreateNoteUseCase.Execute(ctx, title, content string)`, `UpdateNoteUseCase.Execute(ctx, id, title, content string)`, and the `NoteCreator`/`NoteUpdater` port interfaces in `internal/adapters/http/note_handler.go`). Alternative considered: keep `content` first and append `title` last to minimize call-site diffs — rejected because `title` reads first in the JSON payload and in the domain model's natural "identify, then describe" order, and every call site needs updating regardless of position.
- **Rename `Note.UpdateContent` to `Note.Update(title, content string, now time.Time) error`** rather than adding a separate `SetTitle` method, since `UpdateNoteUseCase` always resends both fields together (full-replacement, matching current `Content`-only behavior) and a single method keeps the "set `UpdatedAt` exactly once per update" invariant in one place.
- **No new domain error for title.** `Title` has no validation rule, so `domain.New` and `Note.Update` only ever return `ErrEmptyContent` (unchanged), never a title-related error.
- **Mongo field is a plain `Title string` with `bson:"title"`, no `omitempty`.** Documents written before this change simply lack the `title` key; the driver decodes a missing key into the zero value (`""`), so old notes read back with an empty title with no migration step. New/updated documents always have the key (even if `""`), consistent with `Content`'s existing tag.
- **JSON field is `title` on all three DTOs** (`createNoteRequest`, `updateNoteRequest`, `noteResponse`), consistent with `content`'s lowercase tag.

## Risks / Trade-offs

- [Full-replacement semantics surprise a client that expects PATCH-like partial updates on title] -> Mitigation: this matches the existing `Content` behavior exactly, so it's a consistent (if minimal) API contract, not a new inconsistency.
- [Every call site of `CreateNoteUseCase.Execute` / `UpdateNoteUseCase.Execute` and the `NoteCreator`/`NoteUpdater` interfaces must change in lockstep] -> Mitigation: the Go compiler catches every missed call site at build time; task breakdown updates all of them together.
