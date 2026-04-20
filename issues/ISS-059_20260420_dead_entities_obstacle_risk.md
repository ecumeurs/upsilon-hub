# Issue: Dead Entities Considered Obstacles

**ID:** `ISS-059_20260420_dead_entities_obstacle_risk`
**Ref:** `ISS-059`
**Date:** 2026-04-20
**Severity:** High
**Status:** Open
**Component:** `upsilonbattle/battlearena/ruler/rules`, `battleui/resources/js`
**Affects:** `battleui/resources/js/Pages/BattleArena.vue`, `upsilonbattle/battlearena/controller/controllers/aggressive.go`

---

## Summary

Dead entities (HP <= 0) are incorrectly treated as obstacles on both the frontend and the backend. This prevents units from moving onto or through tiles containing corpses, despite corpses typically being non-blocking. The root cause is a combination of missing cleanup logic in the backend `UseSkill` handler and missing vitality checks in the frontend pathfinding neighbors calculation.

---

## Technical Description

### Background

In a tactical battle, once an entity's HP reaches zero, it should be removed from the active grid occupancy map to allow other entities to pass through or stop on that cell. While the `Attack` rule correctly handles this, other logic paths (like skills) appear to leave the entity in a "zombie" state on the grid.
Note though that `Attack` is the only method used right now, so while `Skill` does it wrongly it doesn't matter much for now.

### The Problem Scenario

1.  An entity is killed by a Skill (e.g., area damage or direct damage from a skill).
2.  The `UseSkill` logic in `skill.go` updates the entity's HP to 0 but does not remove it from `gs.Grid`, `gs.Entities`, or `gs.Turner`.
3.  The Frontend receives a board update where the entity exists with `hp: 0`.
4.  The Frontend's `getNeighbors` function (used for movement calculation) sees the entity at the position and marks it as `isOccupied = true` without checking HP.
5.  The AI's `AggressiveController` attempts to move, but its `isPathStepBlocked` check sees the `EntityID` still present on the `Grid` cell and considers it blocked.

### Where This Pattern Exists Today

-   **Frontend**: `battleui/resources/js/Pages/BattleArena.vue:306`. The `some` check lacks `&& e.hp > 0`.
-   **Backend**: `upsilonbattle/battlearena/ruler/rules/skill.go:56-66`. The application of damage (`ApplyDirectEffect`) updates HP but there is no subsequent check to remove dead entities.
-   **Backend**: `upsilonbattle/battlearena/controller/controllers/aggressive.go:432`. It checks `c.EntityID != uuid.Nil` without verifying if that entity is actually alive in `KnownEntities`.
-   **Backend**: `upsilonbattle/battlearena/ruler/rules/move.go:123`. The Ruler's movement validation also checks `c.EntityID` without vitality verification.

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | High |
| Impact if triggered | High — breaks tactical movement and AI pathing |
| Detectability | High — units simply cannot move where dead bodies are |
| Current mitigant | None |

---

## Recommended Fix

**Short term:** Update `BattleArena.vue` and `aggressive.go` to ignore entities with `hp <= 0` during occupancy checks.

**Medium term:** Implement a centralized `gs.CheckAndRemoveDead(entityID)` method in the Ruler's GameState and call it in `Attack`, `UseSkill`, and `EndOfTurn` (to handle DoTs like poison).

**Long term:** Ensure the `Grid` occupancy is the single source of truth for movement, and that the Ruler strictly maintains this occupancy in sync with entity vitality.

---

## References

- [BattleArena.vue](battleui/resources/js/Pages/BattleArena.vue)
- [skill.go](upsilonbattle/battlearena/ruler/rules/skill.go)
- [aggressive.go](upsilonbattle/battlearena/controller/controllers/aggressive.go)
- [move.go](upsilonbattle/battlearena/ruler/rules/move.go)
