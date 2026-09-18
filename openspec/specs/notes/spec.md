# Notes Specification

## Purpose

Lets clients create, list, edit, and delete simple text notes, each tracked with its creation and last-updated timestamps.

## Requirements

### Requirement: Create Note
The system SHALL allow creating a note by submitting its content. The system SHALL reject a note with empty content. On success, the system SHALL assign the note an identifier and set both its creation timestamp and last-updated timestamp to the time of creation.

#### Scenario: Successful creation
- **WHEN** a client submits a note with non-empty content
- **THEN** the system creates the note, assigns it an identifier, and returns the note with its content, creation timestamp, and last-updated timestamp (equal to the creation timestamp)

#### Scenario: Rejects empty content
- **WHEN** a client submits a note with empty or missing content
- **THEN** the system rejects the request and creates no note

### Requirement: List Notes
The system SHALL allow listing all existing notes, each with its content, creation timestamp, and last-updated timestamp.

#### Scenario: List returns all notes
- **WHEN** a client requests the list of notes
- **THEN** the system returns every existing note with its content, creation timestamp, and last-updated timestamp

#### Scenario: List when no notes exist
- **WHEN** a client requests the list of notes and none exist
- **THEN** the system returns an empty list

### Requirement: Update Note
The system SHALL allow updating the content of an existing note by its identifier. The system SHALL reject an update to empty content. On success, the system SHALL update the last-updated timestamp while preserving the original creation timestamp.

#### Scenario: Successful update
- **WHEN** a client submits new non-empty content for an existing note's identifier
- **THEN** the system updates the note's content, updates its last-updated timestamp, and keeps its original creation timestamp unchanged

#### Scenario: Update of nonexistent note
- **WHEN** a client submits new content for an identifier that does not match any note
- **THEN** the system rejects the request and changes nothing

#### Scenario: Rejects empty content on update
- **WHEN** a client submits empty or missing content for an existing note's identifier
- **THEN** the system rejects the request and leaves the note unchanged

### Requirement: Delete Note
The system SHALL allow deleting an existing note by its identifier.

#### Scenario: Successful deletion
- **WHEN** a client requests deletion of an existing note's identifier
- **THEN** the system deletes the note so it no longer appears when listing notes

#### Scenario: Deletion of nonexistent note
- **WHEN** a client requests deletion of an identifier that does not match any note
- **THEN** the system rejects the request and deletes nothing
