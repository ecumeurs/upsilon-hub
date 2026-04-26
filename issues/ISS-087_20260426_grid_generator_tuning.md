# Issue: Grid Generator Tuning - Large and Flat Maps

**ID:** `20260426_grid_generator_tuning`
**Ref:** `ISS-087`
**Date:** 2026-04-26
**Severity:** Medium
**Status:** Open
**Component:** `upsilonmapmaker/gridgenerator`, `upsilonapi/bridge`
**Affects:** Battle map generation, tactical gameplay, player experience

---

## Summary

Since the integration of the `gridgenerator`, battle maps have been observed to be consistently too large and lacking in tactical depth due to a lack of obstacles and vertical variation. The current default configuration in the `ArenaBridge` produces "flat" maps that don't leverage the 3D grid's potential. We need to tune the generative options and expose them to the API.

---

## Technical Description

### Background
The `GridGenerator` provides several types of terrain (`Flat`, `Hill`, `River`, `Mountain`) and supports obstacle generation via an `ObstructionRate`. However, the current implementation in `ArenaBridge.go` uses hardcoded values that don't enable obstructions and default to a 10x10 size.

### The Problem Scenario
1. **Maps are too large**: A 10x10 grid is often too sparse for 1v1 character skirmishes, leading to many turns spent just moving.
2. **Lack of Obstacles**: The `GenerateObstrcution` flag is currently `false` by default, resulting in empty plains.
3. **Flatness**: Even with the `Hill` generator, the terrain often feels too uniform.
4. **Hardcoded Config**: The bridge doesn't accept map parameters from the `ArenaStartRequest`.

### Proposed Fixes

#### 1. Bridge Configuration (`upsilonapi/bridge/bridge.go`)
- Enable `GenerateObstrcution = true`.
- Set a default `ObstructionRate` (e.g., `tools.NewIntRange(5, 15)`).
- Reduce default size for smaller matches (e.g., 7x7 or 8x8).

#### 2. API Extension (`upsilonapi/api/input.go`)
- Add `MapOptions` to `ArenaStartRequest`:
  ```go
  type MapOptions struct {
      Width           int    `json:"width"`
      Length          int    `json:"length"`
      Type            string `json:"type"` // flat, hill, river, mountain
      ObstacleDensity int    `json:"obstacle_density"`
  }
  ```

#### 3. Generator Refinement (`upsilonmapmaker/gridgenerator/gridgenerator.go`)
- Review `generateFlat` and `generateHill` to ensure higher frequency of height variations.
- Ensure `GenerateObstrcution` logic correctly replaces ground with obstacles at tactical positions.

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | High |
| Impact if triggered | Medium |
| Detectability | High |
| Current mitigant | Manual map overrides possible in tests, but production uses defaults. |

---

## Recommended Fix

**Short term:**
- Update `ArenaBridge.go` to enable obstacles and reduce default size to 8x8.
- Set a baseline `ObstructionRate`.

**Medium term:**
- Expose map configuration via the API.
- Add "Map Preset" support (e.g., "Arena", "Forest", "Mountain Pass").

---

## References

- [gridgenerator.go](file:///workspace/upsilonmapmaker/gridgenerator/gridgenerator.go)
- [bridge.go](file:///workspace/upsilonapi/bridge/bridge.go)
- [ISS-074_20260423_comprehensive_item_system.md](file:///workspace/issues/ISS-074_20260423_comprehensive_item_system.md) (related to battle integration)
