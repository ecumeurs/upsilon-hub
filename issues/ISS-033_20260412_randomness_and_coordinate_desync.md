# Issue: Lack of Random Seeding and Coordinate Desync in Tactical State

**ID:** `20260412_randomness_and_coordinate_desync`
**Ref:** `ISS-033`
**Date:** 2026-04-12
**Severity:** Medium
**Status:** Open
**Component:** `upsilonbattle/battlearena/ruler`
**Affects:** `upsilonapi/api`, `upsiloncli`

---

## Summary

Combat matches in the Upsilon engine suffer from deterministic initialization because `math/rand` is never seeded. This results in identical "random" starting positions across consecutive matches. Additionally, there is a systemic desync where entity `Position` fields in API DTOs are often zeroed out `(0,0)`, despite the entities being occupying non-zero cells in the engine's grid.

---

## Technical Description

### Background
The tactical grid should provide diverse starting environments. Once placed, an entity's position must be accurately reflected in the board state shared with clients.

### The Problem Scenario
1. `ruler.go` initiates a match. Since `rand.Seed` is not called, the same position sequence is generated.
2. `NewEntity` mapping in `api/output.go` handles the conversion from `position.Position` to `dto.Position`.
3. Observed logs show entities with `position: {x:0, y:0}` in the DTO, even if they occupy other cells on the grid, leading to AI confusion and invalid move path planning.

### Where This Pattern Exists Today
- [/workspace/upsilonbattle/battlearena/ruler/ruler.go](/workspace/upsilonbattle/battlearena/ruler/ruler.go) - initialization logic.
- [/workspace/upsilonapi/api/output.go:110](/workspace/upsilonapi/api/output.go#L110) - position mapping.

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | High |
| Impact if triggered | Medium — Degraded tactical gameplay / AI confusion |
| Detectability | Medium — Visible in CLI/WebUI board data |
| Current mitigant | None |

---

## Recommended Fix

**Short term:** 
1. Call `rand.Seed(time.Now().UnixNano())` during ruler initialization.
2. Investigate and fix the `NewEntity` mapping to ensure `entity.Position` is correctly captured after engine placement.

---

## References

- [grid.go](/workspace/upsilonmapdata/grid/grid.go)
- [output.go](/workspace/upsilonapi/api/output.go)
