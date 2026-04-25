# Issue: Standardize cell access on Y-major layout

**ID:** `20260424_cell_access_y_major_standard`
**Ref:** `ISS-079`
**Date:** 2026-04-24
**Severity:** Medium
**Status:** Open
**Component:** `upsilonapi/api/output.go`, `upsiloncli/internal/dto/types.go`, `upsiloncli/internal/script/bridge.go`, `upsiloncli/tests/scenarios/`, `battleui/`
**Affects:** Every caller that iterates `BoardState.grid.cells`

---

## Summary

The tactical grid is currently serialized width-major (`cells[x][y]`) by the Go engine (see `upsilonapi/api/output.go:210-231`) and documented as such in `communication.md` §4.2. The convention in almost every other game/UI stack — and the shape most JS/TS code naturally writes when iterating rows first — is Y-major (`cells[y][x]`). Multiple test scenarios already iterate the grid as if it were Y-major, silently reading the wrong cell whenever `width != height`, or walking out of bounds when the grid is non-square. We must pick one axis order and enforce it everywhere via a helper so future contributors cannot reintroduce this drift.

The decision here is to migrate the whole ecosystem to **Y-major** (`cells[y][x]`) — rows of columns — since that is what humans, the frontend, and JSON-oriented tooling default to.

---

## Technical Description

### Background

- Engine output: `upsilonapi/api/output.go:194-233` builds `Cells` as `make([][]Cell, g.Width)` then fills `bs.Grid.Cells[x][y]` — width-major.
- CLI DTO: `upsiloncli/internal/dto/types.go:37-42` inherits that shape.
- CLI helper: `upsiloncli/internal/script/bridge.go:479-505` (`jsCellContentAt`) uses `board.Grid.Cells[x][y]` — also width-major, consistent with the engine.
- Tests: several scripts use `cells[y][x]` — inconsistent, producing wrong cells on non-square boards.
- Communication spec: `communication.md` §4.2 states width-major explicitly; any change here is a communication-layer alteration (per `CLAUDE.md` forewords, must warn the user before shipping).

### The Problem Scenario

```
Engine board 10 wide × 6 tall:
  cells[x=0..9][y=0..5]   // current contract
Test code:
  for (y=0; y<h; y++) for (x=0; x<w; x++) cell = cells[y][x];
  // When y=6..9 this reads beyond the inner array length 6 → undefined.
  // When y<6 and x<10 it reads a transposed cell: the obstacle we "found"
  // is at a different column than the one we hand to planTravelToward.
```

This means EC-01 (`edge_movement_obstacle_collision`), EC-09 (`edge_movement_jump_limitations`), EC-15 (`edge_attack_target_invalid_cell`), and EC-16 (`edge_attack_target_no_entity`) are not testing what they advertise. They can pass or fail for unrelated reasons depending on grid dimensions and seed.

### Where This Pattern Exists Today

- Source of truth (width-major today):
  - `upsilonapi/api/output.go:194-233`
  - `upsiloncli/internal/dto/types.go:37-42`
  - `upsiloncli/internal/script/bridge.go:479-505`
- Callers that silently do Y-major (wrong relative to current contract):
  - `upsiloncli/tests/scenarios/edge_movement_obstacle_collision.js:37`
  - `upsiloncli/tests/scenarios/edge_movement_jump_limitations.js:35,68,82`
  - `upsiloncli/tests/scenarios/edge_attack_target_invalid_cell.js:32`
  - `upsiloncli/tests/scenarios/edge_attack_target_no_entity.js:32`
- Documentation:
  - `communication.md` §4.2 ("width-major 2D matrix: `cells[x][y]`") — will flip to Y-major.
  - `upsilonbattle/docs/entity_grid.atom.md` (if it restates the shape)

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | High — tests are already inconsistent with the contract |
| Impact if triggered | Medium — silent wrong-cell reads mask real failures and falsely advertise coverage |
| Detectability | Low — no runtime panic; tests often still "PASS" for the wrong reason |
| Current mitigant | None; spec is correct but unenforced |

---

## Recommended Fix

**Short term (this PR):**
- Introduce a bridge-level helper `upsilon.cellAt(board, x, y)` (and a `forEachCell(board, cb)` iterator) so scenario scripts never index the matrix directly. Helper encapsulates the current width-major access. Re-route all edge/e2e scripts through it.
- Document the axis choice at a single place (`communication.md` §4.2 + a comment above `Grid.Cells` in `output.go`).

**Medium term:**
- Migrate the serialization in `upsilonapi/api/output.go` to emit Y-major (`cells[y][x]`). The helper becomes a no-op layer but stays as the only supported access path. Bump `BoardState.version` semantics or add a `grid.layout: "y-major"` discriminator during the transition window.
- Update `upsiloncli/internal/dto/types.go`, `jsCellContentAt`, `battleui` grid readers, and any Laravel consumers in lockstep.
- Update `communication.md` §4.2, the `entity_grid` atom, and the `[[api_standard_envelope]]`-adjacent ATD docs — this is a **communication layer change** and per `CLAUDE.md` must be flagged to the maintainer before merging.

**Long term:**
- Add an ATD lint / CI guard that fails the build if a scenario script directly subscripts `board.grid.cells`; only the helper is allowed.
- Provide typed DTOs (Go + TS) that expose `CellAt(x, y)` as the only public accessor so the underlying storage layout becomes an implementation detail.

---

## References

- `upsilonapi/api/output.go:194-233` — current width-major serialization
- `upsiloncli/internal/script/bridge.go:479-505` — existing `cellContentAt` helper (width-major)
- `upsiloncli/internal/dto/types.go:37-42` — Grid DTO
- `communication.md` §1.4 & §4.2 — communication contract
- `upsilonbattle/docs/mech_move_validation_move_validation_jump_limitations.atom.md` — consumer that depends on correct cell access
- Tests that currently violate the contract: `edge_movement_obstacle_collision.js`, `edge_movement_jump_limitations.js`, `edge_attack_target_invalid_cell.js`, `edge_attack_target_no_entity.js`
