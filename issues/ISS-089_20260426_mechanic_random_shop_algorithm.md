# Issue: Deterministic Daily Random Shop

**ID:** `20260426_mechanic_random_shop_algorithm`
**Ref:** `ISS-089`
**Date:** 2026-04-26
**Severity:** Medium
**Status:** Open
**Component:** `upsilonbattle/mechanics/shop`
**Affects:** `upsilonbattle/engine`, `battleui/shop`

---

## Summary

Implementation of a daily rotating shop system that provides a deterministic set of items based on the player's ID, their account creation date, and the current calendar date. This ensures that every player has a unique but stable shop for the day, preventing exploits while maintaining roguelike variety and progression.

---

## Technical Description

### Background

Currently, the item system and skill system are being established (see ISS-073). A "Shop" mechanism is needed to allow players to acquire these assets. To prevent players from simply restarting the client or refreshing to get better items (save-scumming), the shop inventory must be tied to the player's identity and the passage of time.

### The Problem Scenario

The generation algorithm must be stable for a 24-hour window and unique per user.

```
1. Input: UserID (UUID), CreationDate (ISO), CurrentDate (YYYY-MM-DD)
2. Seed = Sha256(UserID + CreationDate + CurrentDate)
3. Initialize PRNG with Seed
4. Loop 3 times:
    a. Roll Category (Weapon, Armor, Utility, Skill)
    b. Roll Item from Category Registry
    c. Roll Rarity/Stats
5. Roll 1 time (Bonus Slot):
    a. IF Rand(0, 100) < HighValueThreshold:
        i. Generate High Value Item
    b. ELSE: Slot is Empty
6. Return List[Item]
```

### Where This Pattern Exists Today

This is a new mechanic. Related registry logic is being implemented in `upsilonbattle/battlearena/ruler` and item entities in `upsilonbattle/battlearena/entity`.

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | Medium |
| Impact if triggered | High |
| Detectability | High — users will report "same items every day" if the date isn't included, or "unfair shops" if seed is skewed. |
| Current mitigant | None |

---

## Recommended Fix

**Short term:** Define the `MECHANIC` atom in ATD and implement a Go-based `ShopGenerator` that takes the required inputs and returns a deterministic slice of items.

**Medium term:** Integrate with the `upsilonapi` to expose the shop contents via a dedicated endpoint.

**Long term:** Add "re-roll" tokens or progression-based shop upgrades (higher high-value chance).

---

## References

- [ISS-073_20260423_roguelike_skill_system.md](ISS-073_20260423_roguelike_skill_system.md)
- [api_standard_envelope](docs/api_standard_envelope.atom.md)
