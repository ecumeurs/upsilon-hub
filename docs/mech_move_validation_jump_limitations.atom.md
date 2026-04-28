---
id: temp
status: DRAFT
version: 1.2
parents: []
dependents: []
---

# New Atom

## INTENT

## THE RULE / LOGIC
- **Property:** `JumpHeight` on the entity (see `property.JumpHeight`, int).
- **Default:** `2`. Absence of `JumpHeight` on an entity resolves to `2`.
- **Surface Projection:** `upsilonapi/bridge/bridge.go` projects incoming `X, Y` coordinates to the `TopMostCellAt(x, y)` height before validation. This allows clients to omit `Z` coordinates.
- **Walkable Surfaces:** Both `cell.Ground` and `cell.Dirt` types are considered walkable/targetable by the engine.
- **Validation step:** for each pair of consecutive cells `(a, b)` in the requested path, `|b.Z - a.Z| > JumpHeight` → path rejected with `entity.path.notvalid`.
- **Start step:** the entity's current position to the first path cell is checked with the same constraint.
- **A* pathfinding:** `grid.AStarPath(start, end, jumpHeight, exclude)` prunes neighbours whose Z delta exceeds `jumpHeight`.

## TECHNICAL INTERFACE

## EXPECTATION
