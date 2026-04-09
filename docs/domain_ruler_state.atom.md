---
id: domain_ruler_state
human_name: Ruler State Machine Domain
type: MODULE
layer: ARCHITECTURE
version: 1.0
status: STABLE
priority: 5
tags: []
parents: []
dependents:
  - [[domain_ruler_state_technical_interface]]
  - [[rule_turn_clock]]
  - [[domain_ruler_state_action_validation]]
  - [[domain_ruler_state_game_states]]
  - [[api_ruler_methods]]
  - [[domain_ruler_state_data_custody]]
---
# Ruler State Machine Domain

## INTENT
To aggregate the constituent rules of Ruler State Machine Domain.

## THE RULE / LOGIC
Ensures the Ruler maintains a consistent state of the game, managing transitions and input validation.

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[domain_ruler_state]]`
