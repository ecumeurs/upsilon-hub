# Issue: Engine movement validation bypass on non-ground cells

**ID:** `20260428_engine_movement_obstacle_validation_bypass`
**Ref:** `ISS-091`
**Date:** 2026-04-28
**Severity:** High
**Status:** Resolved
**Component:** `upsilonbattle/battlearena/ruler/rules/move.go`
**Affects:** All tactical movement actions in the engine.

---

## Summary

The engine's movement rule fails to reject paths that end on non-walkable cells (obstacles) when those cells are not ground tiles. This is caused by a silent failure in `CellsForPositions` which skips invalid positions, allowing the validation loop in `move.go` to complete without encountering the obstacle cell. Furthermore, errors from `Grid.MoveEntity` are ignored, allowing the entity state to be updated even if the destination is technically invalid.

---

## Technical Description

### Background

The `move` rule is supposed to validate that every step of a path is walkable (not an obstacle) and that all steps are adjacent. If a cell is an obstacle, it should return `entity.path.obstacle`.

### The Problem Scenario

1.  A client sends a move request to an obstacle cell at `(x, y)`.
2.  `bridge.go` (in `upsilonapi`) resolves the `Z` coordinate using `TopMostGroundAt(x, y)`. Since the cell is an obstacle (not ground), it returns `-1`.
3.  The position `(x, y, -1)` is passed to the engine.
4.  `move.go` calls `Grid.CellsForPositions(path)`.
5.  In `upsilonmapdata/grid/grid.go`, `CellsForPositions` iterates over the path. For `(x, y, -1)`, `g.Cells[p]` returns `nil, false` (out of bounds).
6.  `CellsForPositions` **silently skips** this position and returns a shortened slice of cells.
7.  The validation loop in `move.go` iterates over the shortened slice, which only contains the valid steps BEFORE the obstacle. It finds no errors.
8.  `move.go` calls `ctx.Grid.MoveEntity(from, dest, id)` where `dest` is `(x, y, -1)`.
9.  `MoveEntity` returns an error because `(x, y, -1)` is not in the grid, but `move.go` **ignores the return value**.
10. The entity remains at its old position or in an inconsistent state, but the API returns `200 OK` and `"action move accepted"`.

```
Client Request -> Move to (0,2) [Obstacle]
      |
      v
Bridge -> TopMostGroundAt(0,2) = -1  => Position(0,2,-1)
      |
      v
Rule (move.go) -> CellsForPositions([..., (0,2,-1)])
      |           |--> Skips (0,2,-1)
      |
      v
Validation -> Loop only sees valid cells -> SUCCESS
      |
      v
Execution -> MoveEntity(..., (0,2,-1)) -> ERROR (ignored)
      |
      v
Reply -> 200 OK (Move Accepted)
```

### Where This Pattern Exists Today

- `upsilonbattle/battlearena/ruler/rules/move.go:134` (Call to `CellsForPositions`)
- `upsilonbattle/battlearena/ruler/rules/move.go:136-159` (Validation loop)
- `upsilonbattle/battlearena/ruler/rules/move.go:46` (Ignored `MoveEntity` error)
- `upsilonmapdata/grid/grid.go:402` (`CellsForPositions` skipping logic)
- `upsilonapi/bridge/bridge.go:403` (Use of `TopMostGroundAt` for moves)

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | Medium |
| Impact if triggered | High — Desync between client/server and illegal game moves |
| Detectability | High — Manifests as tests failing to receive expected error keys |
| Current mitigant | None |

---

## Recommended Fix

**Short term:** 
1. Modify `move.go` to ensure `len(cells) == len(req.Path)`. If not, return a path error.
2. Check the return value of `MoveEntity` in `move.go` and handle the error.
3. Modify `bridge.go` to use `TopMostCellAt` instead of `TopMostGroundAt` to ensure we get the actual cell even if it's an obstacle, allowing the rule to see the `Obstacle` type.

**Medium term:** 
Change `CellsForPositions` to return an error or a slice of the same length containing `nil` for invalid positions, forcing the caller to handle them.

---

## References

- [move.go](file:///workspace/upsilonbattle/battlearena/ruler/rules/move.go)
- [grid.go](file:///workspace/upsilonmapdata/grid/grid.go)
- [bridge.go](file:///workspace/upsilonapi/bridge/bridge.go)
