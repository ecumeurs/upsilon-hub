# Issue: Unified Scripting Lifecycle and CI Testing Framework

**ID:** `20260415_scripting_unified_lifecycle_and_ci`
**Ref:** `ISS-044`
**Date:** 2026-04-15
**Severity:** High
**Status:** Open
**Component:** `upsiloncli`
**Affects:** `upsiloncli/scripting.md`, `CI/CD pipelines`, `Developer Experience`

---

## Summary

The current scripting environment in `upsiloncli` requires boilerplate code for common tasks like bot registration, matchmaking, and cleanup. We need to provide unified high-level tools to handle the full bot lifecycle safely and establish a robust framework for CI testing with proper assertions and sanctions.

---

## Technical Description

### Background
Currently, developers manually script sequences like `auth_register`, `matchmaking_join`, and `auth_delete` within each scenario. `upsiloncli/scripting.md` provides examples of these patterns, but they are not encapsulated in the core `upsilon` JS object provided to scripts.

### The Problem Scenario
1.  **Boilerplate Overload:** Every script must manually register a `teardown` to avoid "ghost accounts" in the database.
2.  **Safety Risks:** If a script crashes or a timeout is reached, accounts might be left in the matchmaking queue or in-game, causing state pollution.
3.  **Missing CI Infrastructure:** There is no standard way to define "what is a test," which bots to use for specific scenarios, or how to sanction a test failure in a CI environment (e.g., exiting with non-zero code based on specific assertions).

### Where This Pattern Exists Today
- [scripting.md](file:///workspace/upsiloncli/scripting.md) (Lines 40-87 shows the manual teardown logic)
- [progression_test_winner.js](file:///workspace/upsiloncli/samples/progression_test_winner.js)
- [progression_test_loser.js](file:///workspace/upsiloncli/samples/progression_test_loser.js)

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | High |
| Impact if triggered | Medium (DB bloat, unreliable CI results) |
| Detectability | High (Pipeline failures, dangling DB records) |
| Current mitigant | Manual teardown blocks in script samples |

---

## Recommended Fix

**Short term:**
- Implement `upsilon.bootstrapBot()` and `upsilon.shutdownBot()` helper functions in the CLI's JS bridge.
- Add `upsilon.joinWaitMatch(game_mode)` to handle the wait-pattern for matchmaking events.

**Medium term:**
- Define a JSON/YAML schema for CI test suites that specifies bot profiles, scripts, and expected outcomes.
- Integrate these suites into the `upsilontest` or a dedicated CI command.

**Long term:**
- Move account lifecycle management to a purely ephemeral system (e.g., in-memory DB or automatic TTL for bots).

---

## References

- [scripting.md](file:///workspace/upsiloncli/scripting.md)
- [atd/script_farm.atom.md](file:///workspace/docs/atd/script_farm.atom.md) (assumed path)
