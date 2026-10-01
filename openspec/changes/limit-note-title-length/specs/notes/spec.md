# Spec Delta

## MODIFIED Requirements

### Requirement: Create Note
The system SHALL allow creating a note by submitting its content and an optional title. The system SHALL reject a note with empty content. The system SHALL accept an empty or omitted title. The system SHALL reject a title longer than 200 characters, counting each Unicode character once regardless of its byte encoding. On success, the system SHALL assign the note an identifier and set both its creation timestamp and last-updated timestamp to the time of creation.

#### Scenario: Successful creation
- **WHEN** a client submits a note with non-empty content and a title
- **THEN** the system creates the note, assigns it an identifier, and returns the note with its title, content, creation timestamp, and last-updated timestamp (equal to the creation timestamp)

#### Scenario: Successful creation without a title
- **WHEN** a client submits a note with non-empty content and no title
- **THEN** the system creates the note with an empty title and returns it alongside its content, creation timestamp, and last-updated timestamp

#### Scenario: Rejects empty content
- **WHEN** a client submits a note with empty or missing content
- **THEN** the system rejects the request and creates no note, regardless of the title submitted

#### Scenario: Rejects a title over 200 characters
- **WHEN** a client submits a note with non-empty content and a title longer than 200 characters
- **THEN** the system rejects the request and creates no note

#### Scenario: Accepts a title at exactly 200 characters
- **WHEN** a client submits a note with non-empty content and a title exactly 200 characters long
- **THEN** the system creates the note with that title

### Requirement: Update Note
The system SHALL allow updating the title and content of an existing note by its identifier. The system SHALL reject an update to empty content. The system SHALL accept an empty or omitted title. The system SHALL reject a title longer than 200 characters, counting each Unicode character once regardless of its byte encoding. On success, the system SHALL update the last-updated timestamp while preserving the original creation timestamp.

#### Scenario: Successful update
- **WHEN** a client submits new non-empty content and a title for an existing note's identifier
- **THEN** the system updates the note's title and content, updates its last-updated timestamp, and keeps its original creation timestamp unchanged

#### Scenario: Successful update without a title
- **WHEN** a client submits new non-empty content and no title for an existing note's identifier
- **THEN** the system updates the note's content, sets its title to empty, and updates its last-updated timestamp

#### Scenario: Update of nonexistent note
- **WHEN** a client submits new content for an identifier that does not match any note
- **THEN** the system rejects the request and changes nothing

#### Scenario: Rejects empty content on update
- **WHEN** a client submits empty or missing content for an existing note's identifier
- **THEN** the system rejects the request and leaves the note unchanged, regardless of the title submitted

#### Scenario: Rejects a title over 200 characters on update
- **WHEN** a client submits a title longer than 200 characters for an existing note's identifier
- **THEN** the system rejects the request and leaves the note unchanged
