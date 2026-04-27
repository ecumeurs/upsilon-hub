# Issue: Action Feedback Infrastructure Upgrade (Multi-Target & Synchronous Results)

**ID:** `20260426_action_feedback_infrastructure_upgrade`
**Ref:** `ISS-090`
**Date:** 2026-04-26
**Severity:** High
**Status:** Open
**Component:** `upsilonapi`, `upsilonbattle`, `battleui`, `upsiloncli`
**Affects:** Communication Layer, E2E Tests, UI Tactical Feedback

---

## Summary

The current action feedback system (ActionFeedback DTO) is single-target oriented and provides insufficient data in synchronous API replies. This creates friction in E2E testing (ISS-088) and lacks the architecture needed for multi-target skills. This issue tracks a comprehensive upgrade of the communication layer to support a list of results per action and provide immediate feedback in the `game_action` endpoint.

---

## Technical Description

### Background
Currently, an "Attack" action returns only the attacker's state. Results (damage, HP changes) are only sent via asynchronous broadcasts. This makes it difficult for synchronous callers (like the JS test agent) to verify action consequences immediately.

### The Problem Scenario
1.  **Multi-target skills**: AOE skills cannot be described by the current `ActionFeedback` which only has one `target_id`.
2.  **Synchronous verification**: Tests must wait for a separate `board.updated` event to see if an attack landed, instead of checking the API response directly.
3.  **Documentation Desync**: `communication.md` and Postman collections do not reflect these needed changes.

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | High (Impacts all future skill development) |
| Impact if triggered | High (Breaks current E2E logic and UI reporting) |
| Detectability | High (API response schema change) |
| Current mitigant | None (Top-level fields will be maintained temporarily for BC) |

---

## Recommended Fix

**Short term:**
- Expand `ControllerAttackReply` and `ControllerUseSkillReply` in the Go Engine.
- Update `ActionFeedback` DTO in `upsilonapi` to include a `results` array.
- Update `game_action` handler to return the feedback object.

**Medium term:**
- Update `TacticalActionReport.vue` to iterate over results.
- Update `communication.md` and Postman collections.
- Audit all E2E/Edge tests for compatibility.

**Long term:**
- Deprecate top-level `damage`/`hp` fields in `ActionFeedback` in favor of the `results` list.

---

## References

- [implementation_plan.md](file:///home/vscode/.gemini/antigravity/brain/9444671f-a06b-443e-b6e7-5321dd3b813b/implementation_plan.md)
- [ISS-088](file:///workspace/issues/ISS-088_20260426_credit_economy_payload_mismatch.md)
- [communication.md](file:///workspace/communication.md)
