## ADDED Requirements

### Requirement: Normalize raw response body before JSON parsing
The system SHALL normalize the raw HTTP response body by trimming leading and trailing whitespace and removing UTF-8 BOM before JSON parsing.

#### Scenario: Body contains leading BOM
- **WHEN** the response body begins with a UTF-8 BOM
- **THEN** the system removes the BOM before attempting JSON parsing

#### Scenario: Body contains trailing whitespace
- **WHEN** the response body contains trailing whitespace or newlines
- **THEN** the system trims the whitespace before attempting JSON parsing

### Requirement: Surface parsing failures with context
The system MUST return a structured error that includes provider name and a short parsing context when JSON decoding fails.

#### Scenario: Malformed JSON response
- **WHEN** the response body is not valid JSON
- **THEN** the system returns a parsing error that includes provider name and a truncated body preview

### Requirement: Handle known provider response variants
The system SHALL support parsing known provider response variants such as top-level objects and alternative envelope fields used by supported providers.

#### Scenario: Provider returns alternative envelope
- **WHEN** a provider returns a response with a documented alternate envelope field
- **THEN** the system extracts the response payload from that envelope and parses it successfully
