---
id: mech_move_validation_jump_limitations
human_name: "Movement Jump Limitations"
type: MECHANIC
layer: IMPLEMENTATION
status: DRAFT
priority: 3
version: 1.2
tags: [movement, validation, jump]
parents:
  - [[shared:requirement_req_trpg_game_definition]]
dependents: []
---

# Movement Jump Limitations

## INTENT
Reject movement paths whose vertical step between consecutive cells exceeds the entity's jump height.

## THE RULE / LOGIC
- **Property:** `JumpHeight` on the entity (see `property.JumpHeight`, int).
- **Default:** `2`. Absence of `JumpHeight` on an entity resolves to `2`.
- **Surface Projection:** `upsilonapi/bridge/bridge.go` projects incoming `X, Y` coordinates to the `TopMostCellAt(x, y)` height before validation. This allows clients to omit `Z` coordinates.
- **Walkable Surfaces:** Both `cell.Ground` and `cell.Dirt` types are considered walkable/targetable by the engine.
- **Validation step:** for each pair of consecutive cells `(a, b)` in the requested path, `|b.Z - a.Z| > JumpHeight` → path rejected with `entity.path.notvalid`.
- **Start step:** the entity's current position to the first path cell is checked with the same constraint.
- **A* pathfinding:** `grid.AStarPath(start, end, jumpHeight, exclude)` prunes neighbours whose Z delta exceeds `jumpHeight`.

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[mech_move_validation_jump_limitations]]`
- **Code:** `upsilonbattle/battlearena/ruler/rules/move.go`, `upsilonapi/bridge/bridge.go`, `grid.AStarPath`.
- **Test Names:** `edge_movement_jump_limitations` (EC-09).

## EXPECTATION
A path step exceeding `JumpHeight` in absolute Z delta is rejected with `entity.path.notvalid`; steps within `JumpHeight` are permitted.
