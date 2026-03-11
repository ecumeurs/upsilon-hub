# Issue: Duplicate Turn Info received in AggressiveController

**ID:** `20260311_duplicate_turn_info_aggressive_controller`
**Ref:** `ISS-011`
**Date:** 2026-03-11
**Severity:** Low
**Status:** Resolved
**Component:** `upsilonbattle/battlearena/controller/controllers`
**Affects:** `AggressiveController`

---

## Summary

During arena initialization, the `AggressiveController` receives two consecutive `GetEntitiesStateReply` messages. 

**Resolution:** This is a natural occurrence in the broader architecture. The first update happens at inscription (`SetQueue`), and the second provides the final authoritative state at `BattleStart`. While redundant in rapid automated tests, this ensures consistency in live scenarios where controllers might join one at a time.

---

## Technical Description

### Background
When a battle starts, controllers (including the `AggressiveController`) receive a `BattleStart` signal. Human-controlled players via `HTTPController` bridge these to webhooks. IA-driven controllers like `AggressiveController` respond to sequence signals to update their internal state and decide actions.

### The Problem Scenario
In `main_test.go` from `upsilonapi`, logs show:
```
32: time="2026-03-11T13:52:55Z" level=info msg="New Turn Info Received" ... message_type=rulermethods.GetEntitiesStateReply ...
33: time="2026-03-11T13:52:55Z" level=info msg="New Turn Info Received" ... message_type=rulermethods.GetEntitiesStateReply ...
```
Both pulses happen at the exact same timestamp with identical content.

### Where This Pattern Exists Today
Log lines 32 and 33 in the test output of `TestArenaStartEndpoint`.

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | High |
| Impact if triggered | Low |
| Detectability | High |
| Current mitigant | None, appears to be harmless but redundant. |

---

## Recommended Fix

**Short term:** Investigate the `AggressiveController`'s initialization logic and `BattleStart` handler to see if it triggers redundant state requests.
**Medium term:** Ensure state requests are debounced or only triggered once upon entry into the battle.

---

## References

- [test_output.log](file:///workspace/upsilonapi/test_output.log)
