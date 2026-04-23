# Issue: Player Choosing Facing Direction on Pass

**ID:** `20260423_pass_choose_facing`
**Ref:** `ISS-072`
**Date:** 2026-04-23
**Severity:** Medium
**Status:** Open
**Component:** `battleui`, `upsilonapi`, `upsilonbattle`
**Affects:** `upsilonbattle/battlearena/ruler/rules/pass.go`, `battleui` ActionPanel

---

## Summary

When a player chooses to "Pass" their turn, they must be given the option to select their final facing direction (Up, Right, Down, Left). This adds a tactical layer to passing, allowing players to prevent exposing their back to enemies.

---

## Technical Description

### Background
Currently, passing a turn ends the actor's activation immediately, leaving them facing whatever direction they were last looking.

### The Problem Scenario
1. A player moves to a position and has no valid actions left.
2. They "Pass".
3. They are left facing a wall or away from the enemy.
4. An enemy rogue uses this to easily backstab (ISS-070).

### Proposed Solution
*   **API Update:** The `Pass` command should accept an optional `orientation` parameter.
*   **Engine Update:** `upsilonbattle` should update the entity's orientation before completing the pass.
*   **UI Update:** When "Pass" is clicked, show a small cardinal direction selector (similar to movement/targeting).

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | High |
| Impact if triggered | Medium (Tactical depth) |
| Detectability | High |
| Current mitigant | None |

---

## Recommended Fix

**Short term:** Update the `Pass` API to accept orientation and update the engine logic.
**Medium term:** Implement the UI selection overlay in the `battleui`.
**Long term:** Add "Auto-Facing" options for defensive stances.

---

## References

- [ISS-070_20260422_backstabbing_mechanics.md](file:///workspace/issues/ISS-070_20260422_backstabbing_mechanics.md)
- [mec_backstabbing_mechanic.atom.md](file:///workspace/docs/mec_backstabbing_mechanic.atom.md)
