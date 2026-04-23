# Issue: Shielding Credit Attribution System

**ID:** `20260423_shielding_credit_attribution`
**Ref:** `ISS-078`
**Date:** 2026-04-23
**Severity:** Medium
**Status:** Open
**Component:** `upsilonbattle/battlearena`
**Affects:** Credit Earning, Buff System

---

## Summary

Design and implement a robust system for attributing credits earned through damage mitigation (shields/blocking). The current system lacks a way to track the original caster of a shield when it absorbs damage later in the battle. This requires linking buffs to effects and effects to casters, and handling credit calculation when the mitigation is "popped".

---

## Technical Description

### Background
In [ISS-067](ISS-067_20260422_credit_economy_shop.md), we established the rule: "1 HP mitigated = 1 credit for the caster". However, shields are currently stored as a numeric counter on the target entity, losing the identity of the player who provided the shield.

### The Problem Scenario
1. **Player A** casts a shield on **Ally B**.
2. **Ally B** is later attacked by **Enemy C**.
3. The shield absorbs 10 damage.
4. **Player A** should receive 10 credits, but the engine only sees B's shield value decrease.

### Proposed Direction
- **Buff to Effect Link:** Shields should be treated as buffs/effects.
- **Caster Persistence:** The effect must persist the `CasterID` even after the initial application.
- **Earning Trigger:** When damage is subtracted from a shield, a trigger must look up the associated effect's caster and award credits.
- **Multi-Caster Handling:** Define how overlapping shields from different casters should prioritize credit earning (e.g., FIFO, LIFO, or proportional).

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | High |
| Impact if triggered | Medium |
| Detectability | High |
| Current mitigant | Shielding credits are currently disabled in [ISS-067]. |

---

## Recommended Fix

**Short term:** Design the `Effect` -> `Buff` -> `Caster` traceability in the engine.
**Medium term:** Update `ApplyDirectEffect` and `Attack` logic to use this traceability for credit calculation.

---

## References

- [ISS-067_20260422_credit_economy_shop.md](ISS-067_20260422_credit_economy_shop.md)
- `upsilonbattle/battlearena/property/effect/effect.go`
- `upsilonbattle/battlearena/ruler/rules/attack.go`
