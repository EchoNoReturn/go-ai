## Context

The llm package integrates multiple providers with slightly different response envelopes. Current parsing assumes well-formed JSON bodies and consistent top-level shapes, which breaks on whitespace/BOM issues and provider-specific variants. Errors can be opaque, making debugging and recovery difficult.

## Goals / Non-Goals

**Goals:**
- Normalize raw response bodies before JSON parsing.
- Centralize response decoding to handle provider variants consistently.
- Improve error context to speed triage and reduce silent failures.
- Add tests that cover parsing edge cases and provider variants.

**Non-Goals:**
- Changing provider APIs or request payload structures.
- Introducing new external dependencies beyond standard library.
- Redesigning the overall session interface or public API shape.

## Decisions

- **Introduce a normalization step before JSON decoding.**
  - Rationale: Eliminates common failures from BOM and whitespace without touching provider logic.
  - Alternatives: Provider-specific trimming; rejected due to duplicated logic and inconsistent behavior.

- **Centralize response decoding in llm package helpers.**
  - Rationale: Ensures consistent handling and error formatting across providers.
  - Alternatives: Keep decoding inside each provider file; rejected because bug fixes would need to be replicated and could drift.

- **Structured parse errors with provider context and body preview.**
  - Rationale: Helps users and developers identify malformed responses quickly.
  - Alternatives: Return raw JSON errors; rejected due to poor debuggability.

## Risks / Trade-offs

- **[Risk]** Over-normalizing could mask provider-side issues → **Mitigation:** Keep normalization minimal (BOM + whitespace) and return clear parse errors on failure.
- **[Risk]** Variant handling might accept unintended shapes → **Mitigation:** Restrict to documented provider envelopes and add tests per variant.
- **[Risk]** Error previews could include sensitive data → **Mitigation:** Truncate previews and avoid logging by default; only return in error struct.
