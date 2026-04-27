# Issue: Flaky TestBattleFullRoundtrip Grid Adjacency

**ID:** `20260427_flaky_battle_roundtrip_grid_adjacency`
**Ref:** `ISS-091`
**Date:** 2026-04-27
**Severity:** Medium
**Status:** Open
**Component:** `upsilonapi`
**Affects:** `upsilonapi/TestBattleFullRoundtrip`

---

## Summary

`TestBattleFullRoundtrip` occasionally fails with the error "Entity is not adjacent to the first move on a randomly-generated grid". This is a pre-existing flaky test issue caused by the random grid seed producing invalid or incompatible configurations for the test's move logic.

---

## Technical Description

### Background
The `TestBattleFullRoundtrip` test performs an end-to-end simulation of a battle. This includes generating a grid, adding controllers, and executing moves.

### The Problem Scenario
1. The test initializes a battle which triggers a random grid generation.
2. The test logic attempts to perform a move.
3. If the random grid layout places the entity and the target move in non-adjacent cells (or otherwise violates adjacency rules expected by the engine for that specific seed), the test fails.
4. The failure is non-deterministic, as it depends on the random seed used during grid generation.

### Where This Pattern Exists Today
- `upsilonapi/battle_test.go` (presumably, based on the `TestBattleFullRoundtrip` name)
- The grid generation logic in the engine used by the API tests.

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | Medium |
| Impact if triggered | Medium |
| Detectability | High — fails visibly in CI with adjacency error |
| Current mitigant | Retrying the test usually passes (as seen in recent logs) |

---

## Recommended Fix

**Short term:** 
- Use a fixed seed for the random grid generator within `TestBattleFullRoundtrip` to ensure consistent test behavior.
- Alternatively, add a retry mechanism at the test level if it depends on randomness.

**Medium term:** 
- Implement "suitability checks" in the grid generator to ensure it doesn't generate "broken" states for standard E2E tests.

**Long term:** 
- Transition roundtrip tests to use predefined, static grid layouts instead of random ones to ensure 100% determinism.

---

## References

- [battle_test.go](file:///workspace/upsilonapi/battle_test.go) (Target test file)
- [grid_generator.go](file:///workspace/upsilon-hub/upsilonbattle/pkg/game/grid_generator.go) (Potential source of randomness, path needs verification)
