---
id: rule_mapmaker_seed_determinism
status: STABLE
priority: 2
version: 1.0
parents: []
human_name: "Procedural Seed Determinism"
type: RULE
dependents: []
layer: BUSINESS
---

# New Atom

## INTENT
Ensure that procedural map generation is perfectly deterministic given a specific seed.

## THE RULE / LOGIC
- Every random choice in the generation algorithm must be derived from the provided seed.
- Sequential calls with the same seed must produce bit-identical grid structures.
- Use the shared `upsilontools/tools` package for all randomization.

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[mapmaker_seed_determinism]]`
- **Test Names:** `TestSeedDeterminism`

## EXPECTATION
