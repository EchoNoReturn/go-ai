## 1. Response Normalization

- [x] 1.1 Locate existing response parsing paths in `llm` providers and session flow
- [x] 1.2 Add a normalization helper that trims whitespace and removes UTF-8 BOM
- [x] 1.3 Integrate normalization into provider response decoding paths

## 2. Parsing and Error Handling

- [x] 2.1 Centralize JSON decoding into a shared helper with provider context
- [x] 2.2 Add structured parse errors with truncated body previews
- [x] 2.3 Handle documented provider envelope variants in the shared decoder

## 3. Tests

- [x] 3.1 Add unit tests for normalization (BOM/whitespace)
- [x] 3.2 Add tests for malformed JSON error surfacing
- [x] 3.3 Add tests covering provider-specific envelope variants
