# Issue: Action Feedback Infrastructure Verification Plan

**ID:** `20260427_action_feedback_verification_plan`
**Ref:** `ISS-092`
**Date:** 2026-04-27
**Severity:** High
**Status:** Open
**Component:** `upsilonapi/handler`, `upsilonapi/api`, `upsilonbattle/battlearena/ruler/rules`
**Affects:** `CLI`, `BattleUI`, `E2E Tests`

---

## Summary

The Action Feedback system has been refactored (ISS-090) to support multi-target results and detailed impact metrics (damage, HP deltas, credits) in synchronous API responses. This change alters the communication layer between the Engine and the Gateway/UI. A comprehensive verification plan is required to ensure no regressions in existing tactical flows and to validate the new data structures.

---

## Technical Description

### Background
Previously, `ControllerAttackReply` only returned the `Entity` (attacker) state. Multi-target skills were not fully supported in synchronous feedback, relying instead on asynchronous webhooks which were often delayed or deduplicated.

### The Problem Scenario
1. **Field Renaming:** `Entity` was renamed to `Attacker` in response structs.
2. **Structure Change:** `Results []ActionResult` was added to provide per-target feedback.
3. **Implicit Usage:** The CLI and E2E tests might rely on the old `data` structure of the synchronous response.

---

## Verification Plan

### 1. Unit Testing (Engine)
- [x] Verify `EffectApplicator` returns correct `[]ActionResult`.
- [x] Verify `Attack` logic populates `Results`.
- [x] Verify `UseSkill` logic populates `Results` for multiple targets.
- [x] Fix all broken tests in `battlearena/ruler/rules`.

### 2. Integration Testing (Bridge/API)
- [x] Verify `upsilonapi` maps engine results to `api.ActionResult`.
- [x] Verify `HandleArenaAction` includes `results` in JSON response.
- [/] Verify `TestBattleFullRoundtrip` (Currently flaky, see ISS-091).

### 3. E2E Regression Testing (Post-Commit)
The following CI tests MUST be run and passed:
- `e2e_match_resolution_standard.js`: Verify damage leads to correct match end.
- `e2e_credit_economy.js`: Verify credits are awarded and visible in feedback.
- `e2e_friendly_fire_prevention.js`: Verify AOE skills handle targets correctly.
- `e2e_inventory_equip_battle.js`: Verify item-based damage boosts are reflected.

### 4. CLI Validation
- [ ] Update CLI to display `Damage` and `HP` changes from the new `results` array.
- [ ] Verify `upsilon-cli attack` prints impact details.

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | High |
| Impact if triggered | High (Broken UI/CLI feedback) |
| Detectability | High (Tests fail or CLI shows no damage) |
| Current mitigant | Synchronous tests in `upsilonapi` |

---

## Recommended Fix

**Short term:** Complete the manual verification of the CLI and run the specific E2E tests identified above.
**Medium term:** Update the CLI's result parser to prioritize the `results` array over legacy fields.
**Long term:** Standardize all tactical feedback to use the `ActionResult` pattern, including positional effects.

---

## References

- [communication.md](file:///workspace/communication.md)
- [handler.go](file:///workspace/upsilonapi/handler/handler.go)
- [rules/attack.go](file:///workspace/upsilonbattle/battlearena/ruler/rules/attack.go)
- [rules/skill.go](file:///workspace/upsilonbattle/battlearena/ruler/rules/skill.go)
