## Why

LLM provider responses occasionally fail to parse due to response body formatting edge cases, which causes chat completion failures and inconsistent behavior. Fixing this now improves reliability across providers and prevents silent data loss in downstream message handling.

## What Changes

- Normalize and validate raw HTTP response bodies before parsing.
- Add stricter JSON decoding with clear error surfaces and fallback handling for known response variants.
- Improve response formatting/parsing coverage in the llm package with targeted tests.

## Capabilities

### New Capabilities
- `llm-response-normalization`: Normalize and parse provider response bodies consistently across OpenAI/Anthropic/DeepSeek.

### Modified Capabilities
- (none)

## Impact

- llm package response parsing and error handling.
- Provider-specific request/response adapters.
- Tests for response body decoding and formatting edge cases.
