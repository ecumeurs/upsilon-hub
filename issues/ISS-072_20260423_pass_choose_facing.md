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

**Short term:** Update the `Pass` API to accept orientation and update the engine logic. (will need a dedicated endpoint and atd atom, and communication and postman collection update.)
**Medium term:** Implement the UI selection overlay in the `battleui`.
Ensure the CLI is up to date as well. 
Create a proper E2E scenario for this 
---

---

## Implementation Progress

### Facing Indicator (2026-04-25) — `battleui` ✅

A facing direction indicator has been added to `ThreeGrid.vue` (the 3D board component).

**How it works:**
- Each entity may carry an optional `facing` field in the wire format with values `"Up"` | `"Right"` | `"Down"` | `"Left"`.
  - `"Up"` = grid Y+1 (Three.js +Z)
  - `"Right"` = grid X+1 (Three.js +X)
  - `"Down"` = grid Y-1 (Three.js −Z)
  - `"Left"` = grid X-1 (Three.js −X)
- When `facing` is present and non-null, `ThreeGrid.vue` renders a small dark-green triangle (`#1a5c1a`) flat on the cell surface, tip pointing toward the facing edge.
- When `facing` is absent or null the indicator is silently omitted — fully backward-compatible.

**Wire format change needed (upsilonapi / upsilonbattle):**
Add `facing` to the entity payload returned by the turn/state endpoint (same shape as `EntityOrientation` already used internally: `"Up"` / `"Right"` / `"Down"` / `"Left"`). No breaking change — absence of the field is safe.

**Test fixture:** `/__test/component/grid-facing` — four pawns each facing a different cardinal direction, with the dark-green triangle visible. Playwright snapshot baseline committed at `tests/playwright/__snapshots__/components.spec.ts-snapshots/component-grid-facing-chromium.png`.

**Still to do:**
- Backend: expose `facing` in entity wire format (`upsilonapi`)
- Engine: honour `orientation` parameter on the `Pass` command (`upsilonbattle`)
- UI: facing picker overlay in ActionPanel (separate task — the indicator above is the display half)
- CLI: update to send `orientation` with pass command

---

## References

- [ISS-070_20260422_backstabbing_mechanics.md](file:///workspace/issues/ISS-070_20260422_backstabbing_mechanics.md)
- [mec_backstabbing_mechanic.atom.md](file:///workspace/docs/mec_backstabbing_mechanic.atom.md)
