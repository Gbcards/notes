# Design

## Context

See proposal.md - Why. `Note.New` and `Note.Update` (`internal/domain/note/note.go`) already enforce one business rule on content (`ErrEmptyContent`) and the HTTP adapter (`internal/adapters/http/note_handler.go`, `writeError`) already maps that domain error to `400 Bad Request`. The title length limit follows the same existing pattern.

## Goals / Non-Goals

**Goals:**
- Reject titles over 200 characters at creation and update, counting characters the way a person reading them would (accented letters like `ñ`/`á` count as one each).

**Non-Goals:**
- Unicode normalization (NFC/NFD) of titles. Out of scope - see Decisions.
- Truncating or trimming titles. The system rejects over-limit titles rather than silently modifying client input, consistent with how empty content is rejected rather than defaulted.

## Decisions

**Count length in runes, not bytes.** Go's `len(string)` counts UTF-8 bytes; `ñ`, `á`, etc. are multi-byte, so a byte-based limit would cut off titles shorter than 200 visible characters. Use `utf8.RuneCountInString` instead, applied in `Note.New` and `Note.Update` alongside the existing content check.

**Do not add Unicode normalization.** A title's `ñ` could theoretically arrive as a decomposed sequence (base letter + combining tilde = 2 runes) rather than the single precomposed code point produced by normal keyboard input and browsers. Handling that would require normalizing to NFC first (e.g. `golang.org/x/text/unicode/norm`), a new dependency, for a case that doesn't occur in practice for this client base. Plain rune counting is sufficient.

**New domain error, mapped like the existing one.** Add `ErrTitleTooLong` next to `ErrEmptyContent` in `internal/domain/note/errors.go`. `writeError` in the HTTP adapter gains one more `case` mapping it to `400 Bad Request`, mirroring the existing `ErrEmptyContent` case.

## Risks / Trade-offs

- [Decomposed Unicode titles count 1 rune too long per decomposed character] -> Accepted; not a realistic input for this service, and avoids an extra dependency.
