# Spec Delta

## ADDED Requirements

### Requirement: Note Tags
The system SHALL allow a note to carry an optional list of tags, submitted when creating or updating the note and returned whenever the note is returned. The system SHALL normalize each submitted tag by removing leading and trailing whitespace and converting it to lowercase. The system SHALL discard tags that are empty after normalization. The system SHALL remove duplicate tags without error, keeping each tag at the position of its first occurrence. The system SHALL reject a tag longer than 20 characters after normalization, counting each Unicode character once regardless of its byte encoding. The system SHALL reject a note with more than 5 tags after normalization and duplicate removal. The system SHALL reject a tags value that is not a list of strings. On update, the submitted tags SHALL replace the note's entire tag list, and omitting tags SHALL leave the note with no tags. A note without tags SHALL be returned with an empty list of tags, never a null value, including notes stored before tags existed.

#### Scenario: Creation with tags
- **WHEN** a client creates a note with non-empty content and tags `["go", "backend"]`
- **THEN** the system creates the note and returns it with tags `["go", "backend"]`

#### Scenario: Creation without tags
- **WHEN** a client creates a note with non-empty content and no tags
- **THEN** the system creates the note and returns it with an empty list of tags

#### Scenario: Duplicate tags are removed
- **WHEN** a client creates a note with tags `["go", "go", "backend"]`
- **THEN** the system creates the note without error and returns it with tags `["go", "backend"]`

#### Scenario: Tags are normalized
- **WHEN** a client creates a note with tags `["Go", " go ", "GO"]`
- **THEN** the system stores and returns the single tag `["go"]`

#### Scenario: First occurrence order is preserved
- **WHEN** a client creates a note with tags `["backend", "go", "Backend", "api"]`
- **THEN** the system returns tags `["backend", "go", "api"]`

#### Scenario: Empty and whitespace-only tags are discarded
- **WHEN** a client creates a note with tags `["go", "", "  "]`
- **THEN** the system returns tags `["go"]`

#### Scenario: Accepts a tag at exactly 20 characters
- **WHEN** a client creates a note with a tag exactly 20 characters long after trimming surrounding whitespace
- **THEN** the system creates the note with that tag

#### Scenario: Rejects a tag over 20 characters
- **WHEN** a client creates or updates a note with a tag longer than 20 characters after trimming surrounding whitespace
- **THEN** the system rejects the request and creates or changes no note

#### Scenario: Rejects more than 5 tags
- **WHEN** a client creates or updates a note with tags that are still more than 5 after normalization and duplicate removal
- **THEN** the system rejects the request and creates or changes no note

#### Scenario: Duplicates do not count toward the limit
- **WHEN** a client creates a note with 6 submitted tags that reduce to 5 distinct tags after normalization
- **THEN** the system creates the note with those 5 tags

#### Scenario: Rejects tags that are not a list of strings
- **WHEN** a client creates or updates a note with a tags value such as `"go"` or `[1, 2]`
- **THEN** the system rejects the request as invalid and creates or changes no note

#### Scenario: Update replaces the tag list
- **WHEN** a client updates an existing note that has tags `["go", "backend"]`, submitting tags `["API", "api", "db"]`
- **THEN** the system replaces the note's tags with `["api", "db"]` and updates its last-updated timestamp

#### Scenario: Update without tags clears them
- **WHEN** a client updates an existing note that has tags, submitting non-empty content and no tags
- **THEN** the system sets the note's tags to an empty list and updates its last-updated timestamp

#### Scenario: Notes stored before tags existed
- **WHEN** a client lists notes that include a note stored before tags existed
- **THEN** the system returns that note with an empty list of tags and no error

## MODIFIED Requirements

### Requirement: List Notes
The system SHALL allow listing all existing notes, each with its title, content, tags, creation timestamp, and last-updated timestamp.

#### Scenario: List returns all notes
- **WHEN** a client requests the list of notes
- **THEN** the system returns every existing note with its title, content, tags, creation timestamp, and last-updated timestamp

#### Scenario: List when no notes exist
- **WHEN** a client requests the list of notes and none exist
- **THEN** the system returns an empty list
