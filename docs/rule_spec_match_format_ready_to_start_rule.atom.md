---
id: rule_spec_match_format_ready_to_start_rule
status: DRAFT
type: RULE
layer: BUSINESS
parents: []
dependents: []
human_name: "Ready to Start Rule"
priority: 3
tags: [match,readiness]
version: 1.0
---

# New Atom

## INTENT
Define the mandatory conditions that must be met before a match is permitted to transition to the 'InProgress' state.

## THE RULE / LOGIC
- A valid, non-nil Grid must be initialized and associated with the GameState.
- The required number of Controllers (NbControllers) must be registered.
- All entities must be positioned within the Grid boundaries.

## TECHNICAL INTERFACE

## EXPECTATION
