# Issue: Credit Economy E2E Test Payload Misinterpretation

**ID:** `20260426_credit_economy_payload_mismatch`
**Ref:** `ISS-088`
**Date:** 2026-04-26
**Severity:** Medium
**Status:** Open
**Component:** `upsiloncli/tests/scenarios`
**Affects:** `e2e_credit_economy.js`

---

## Summary

The `e2e_credit_economy.js` test fails because it misinterprets the response from the `game_action` (attack) API. It expects the response to contain the state of the *target* entity (to calculate damage dealt), but the engine currently returns the state of the *attacker* entity.

---

## Technical Description

### Background
The `game_action` endpoint for an "attack" type action returns the updated state of an entity. The `e2e_credit_economy.js` script uses this response to determine how much damage was dealt by comparing the returned HP with the HP before the attack.

### The Problem Scenario
1.  Bot A (Herald, 30 HP) attacks Bot B (Herald, 30 HP).
2.  The engine processes the attack, dealing damage to Bot B.
3.  The engine returns a `ControllerAttackReply` containing the state of **Bot A** (the attacker).
4.  The JS script receives the response and reads `result.hp` (which is 30, Bot A's HP).
5.  The script compares this with `foeHpBefore` (which was 30).
6.  `myDamageDealt` is calculated as `30 - 30 = 0`.
7.  The script asserts `myDamageDealt > 0` and fails.

### Where This Pattern Exists Today
- [e2e_credit_economy.js](file:///workspace/upsiloncli/tests/scenarios/e2e_credit_economy.js#L67-L68)
- [attack.go](file:///workspace/upsilonbattle/battlearena/ruler/rules/attack.go#L147-L149)
- [handler.go](file:///workspace/upsilonapi/handler/handler.go#L65-L66)

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | High (Always fails when both entities have same HP) |
| Impact if triggered | Medium (Blocks credit economy verification) |
| Detectability | High (Test fails with "Attack landed but no damage was reported") |
| Current mitigant | None |

---

## Recommended Fix

**Short term:** Update `e2e_credit_economy.js` to fetch the foe's state from the character list or wait for the `board.updated` event instead of relying on the `game_action` return value.
**Medium term:** Standardize the `ControllerAttackReply` to include both attacker and target states, or a summary of the action effects (damage dealt, buffs applied, etc.).

---

## References

- [e2e_credit_economy.js](file:///workspace/upsiloncli/tests/scenarios/e2e_credit_economy.js)
- [attack.go](file:///workspace/upsilonbattle/battlearena/ruler/rules/attack.go)
