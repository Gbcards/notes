# Proposal

## Why

Note titles currently have no length limit, so a client can submit an arbitrarily long title. Capping it at 200 characters keeps titles reasonable for display and storage.

## What Changes

- Creating or updating a note rejects a title longer than 200 Unicode characters (counted as runes, not bytes, so accented characters like `ñ`/`á` count as one each).
- The HTTP API returns `400 Bad Request` when the title exceeds the limit, consistent with the existing empty-content rejection.

## Capabilities

### New Capabilities

(none)

### Modified Capabilities

- `notes`: the create and update requirements gain a constraint that the title must not exceed 200 characters.

## Impact

- `internal/domain/note`: new `ErrTitleTooLong`; `New` and `Update` validate title length.
- `internal/adapters/http`: `writeError` maps `ErrTitleTooLong` to `400 Bad Request`.
- No API shape change (still a `title` string field); no new dependencies.
