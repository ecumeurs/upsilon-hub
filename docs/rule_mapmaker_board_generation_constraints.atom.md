---
id: rule_mapmaker_board_generation_constraints
status: STABLE
human_name: "Board Generation Constraints"
version: 1.0
parents: []
type: RULE
layer: BUSINESS
priority: 2
dependents: []
---

# New Atom

## INTENT
Define the physical and logical constraints for generated tactical boards.

## THE RULE / LOGIC
- Generated boards must stay within the specified Width, Length, and Height ranges.
- Every ground tile must be reachable or correctly marked as an obstacle.
- The top-most layer must be playable (no visible dirt islands).

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[mapmaker_board_generation_constraints]]`
- **Test Names:** `TestGridBoundaries`

## EXPECTATION
