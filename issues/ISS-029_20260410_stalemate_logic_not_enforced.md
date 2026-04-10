# Issue: PVP Stalemate Logic Not Enforced

**ID:** `20260410_stalemate_logic_not_enforced`
**Ref:** `ISS-029`
**Date:** 2026-04-10
**Severity:** Medium
**Status:** Open
**Component:** `battleui/app/Http/Controllers/API/WebhookController.php`, `upsilonapi/`
**Affects:** `game_matches` lifecycle

---

## Summary

Currently, there is no logic to detect and resolve a stalemate when remaining characters on both sides have a Defense stat exceeding the maximum possible Attack of the opposition. This can lead to "infinite" matches where neither side can deal damage.

---

## Technical Description

### Background
Matches in Upsilon Battle are resolved through `Attack - Defense = Damage`. If `Defense >= Attack`, damage is zero or negligible.

### The Problem Scenario
A match reaches a state where:
- Team A's Max Attack <= Team B's Min Defense
- Team B's Max Attack <= Team A's Min Defense
Neither side can effectively reduce the other's HP, and the match stays "in progress" indefinitely.

### Where This Pattern Exists Today
- The Go engine (`upsilonapi`) continues to process turns but no meaningful state change occurs.
- The `WebhookController.php` in Laravel simply forwards the state without checking for terminal stalemate conditions.

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | Low (primarily affects high-defense tank builds) |
| Impact if triggered | Medium (match soft-lock) |
| Detectability | Medium (match duration metrics) |
| Current mitigant | None |

---

## Recommended Fix

**Short term:** Implement a stalemate checker in `WebhookController.php` that marks a match as a `DRAW` if the condition is met after turn processing.  
**Medium term:** Implement the logic directly in the Go Engine (`upsilonapi`) to send a `MATCH_TERMINATED` event with a `DRAW` status.  
**Long term:** Introduce active abilities or armor-piercing mechanics to ensure a stalemate is impossible.

---

## References

- [[rule_pvp_stalemate_draw]]
- [WebhookController.php](file:///workspace/battleui/app/Http/Controllers/API/WebhookController.php)
