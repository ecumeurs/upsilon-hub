# Issue: Remove Silent Failure Default Values

**ID:** `20260420_remove_silent_failure_defaults`
**Ref:** `ISS-062`
**Date:** 2026-04-20
**Severity:** Medium
**Status:** Open
**Component:** `battleui/core`
**Affects:** `UpsilonApiService`, `BattleArena.vue`, `config/services.php`

---

## Summary

The codebase contains multiple instances where missing configuration or failed API calls default to fallback values (e.g., empty arrays, default grid sizes, or local dev URLs) instead of failing fast. This pattern hides critical integration errors and makes debugging production issues significantly harder.

---

## Technical Description

### Background
In a robust production environment, the system should "fail fast" when a mandatory dependency or data point is missing. This ensures that the root cause (e.g., missing environment variable) is immediately visible in logs.

### The Problem Scenario
1. An environment variable like `UPSILON_API_URL` is missing in production. 
2. The system defaults to `http://localhost:8081` instead of throwing an error.
3. The application "works" but fails with a cryptic "Connection Refused" error that points to the wrong host, hiding the fact that the configuration is missing entirely.
4. In the UI, failing to fetch a grid might default to a 10x10 empty grid, leading to a "ghost" board state where the user sees obstacles but no tactical data, with no clear error message.

### Where This Pattern Exists Today
- **`battleui/config/services.php`**: Defaults for environment variables that should be mandatory.
- **`battleui/app/Services/UpsilonApiService.php`**: `sendEnvelopeRequest` returns `[]` on failure (line 108) or a hardcoded "success: false" array (lines 113-120) instead of throwing an exception.
- **`battleui/resources/js/Pages/BattleArena.vue`**: Computed properties like `grid` (line 146) and `allEntities` (line 147) default to safe-looking but incorrect empty states.

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | High |
| Impact if triggered | Medium (extended debug time) |
| Detectability | Low — manifests as missing data or "weird" UI behavior |
| Current mitigant | None |

---

## Recommended Fix

**Short term:** Audit `config/` files and remove default values for environment variables that do not have a sensible "global" default (like API URLs).
**Medium term:** Refactor `UpsilonApiService` to throw custom exceptions (e.g., `EngineConnectionException`) that can be caught and handled by the global Laravel Exception Handler.
**Long term:** Implement strict Prop types and validation in Vue components to ensure they do not render unless a valid schema-compliant state is provided.

---

## References

- [UpsilonApiService.php](file:///workspace/battleui/app/Services/UpsilonApiService.php)
- [BattleArena.vue](file:///workspace/battleui/resources/js/Pages/BattleArena.vue)
- [services.php](file:///workspace/battleui/config/services.php)
