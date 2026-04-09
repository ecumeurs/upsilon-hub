---
id: mech_action_economy_time_constraint_rules
human_name: Time Constraint Rules Mechanic
type: MECHANIC
layer: IMPLEMENTATION
version: 1.0
status: STABLE
priority: 5
tags: []
parents: 
  - [[mech_action_economy]]
dependents: []
---
# Time Constraint Rules Mechanic

## INTENT
Defines the temporal constraints for a turn.

## THE RULE / LOGIC
- Turn duration is strictly capped at 30 seconds.

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[mech_action_economy_time_constraint_rules]]`
