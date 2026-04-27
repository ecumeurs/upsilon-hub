# Issue: Engine Hang during Arena Start (Infinite Loop in RandomPosition)

**ID:** `20260427_arena_start_hang_random_position`
**Ref:** `ISS-093`
**Date:** 2026-04-27
**Severity:** Critical
**Status:** Resolved
**Component:** `upsilonmapdata/grid`, `upsilonmapmaker/gridgenerator`
**Affects:** `upsilonapi/HandleArenaStart`, `TestBattleFullRoundtrip`

---

## Summary

The Upsilon Engine could enter an infinite loop during arena initialization (entity placement) when the generated grid's height was exactly equal to the ground elevation. This caused the API to hang indefinitely.

---

## Technical Description

### Background
When an arena starts, the `ArenaBridge` generates a grid and then places entities using `grid.RandomPosition()`. 

### The Problem Scenario
1. `gridgenerator.generateFlat()` used `tools.RandomInt(2, gr.Height)` to determine ground elevation.
2. If `gr.Height` was 2 (a common value for small tests), `ground_height` would be 2.
3. Ground cells were created at `z=2`.
4. `grid.TopMostGroundAt(x, y)` scans from `g.Height` downwards. If `g.Height` is 2, it finds the ground at `z=2`.
5. `grid.RandomPosition()` receives `z=2` and calls `grid.CellAt(pos)`.
6. `grid.PositionIsInGrid()` checks `p.Z < g.Height`. Since `2 < 2` is false, it returns false.
7. `RandomPosition()` loops forever, thinking it hasn't found a valid tile.

### Where This Pattern Exists Today
Fixed in:
- `upsilonmapmaker/gridgenerator/gridgenerator.go` (Ensured `ground_height < gr.Height`)
- `upsilonmapdata/grid/grid.go` (Hardened `RandomPosition` with a retry limit and fallback)

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | Medium (Depends on random seed/grid height range) |
| Impact if triggered | Critical (Infinite loop, service hang) |
| Detectability | High (API request never returns) |
| Current mitigant | None prior to fix. |

---

## Recommended Fix

**Short term:** 
- Apply the boundary fix to `gridgenerator.go`.
- Add retry limits to `RandomPosition` in `grid.go`. (DONE)

**Medium term:** 
- Audit all grid generation types for similar boundary risks.

**Long term:** 
- Implement a more robust "Spawn Registry" that doesn't rely on random sampling for critical placement.

---

## References

- [gridgenerator.go](file:///workspace/upsilonmapmaker/gridgenerator/gridgenerator.go)
- [grid.go](file:///workspace/upsilonmapdata/grid/grid.go)
