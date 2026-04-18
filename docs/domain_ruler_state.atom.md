---
id: domain_ruler_state
human_name: Ruler State Machine Domain
type: MODULE
layer: ARCHITECTURE
version: 1.1
status: STABLE
priority: 5
tags: []
parents: []
dependents:
  - [[api_ruler_methods]]
  - [[domain_ruler_state_action_validation]]
  - [[domain_ruler_state_data_custody]]
  - [[domain_ruler_state_game_states]]
  - [[domain_ruler_state_technical_interface]]
  - [[rule_turn_clock]]
---
# Ruler State Machine Domain

## INTENT
To aggregate the constituent rules of Ruler State Machine Domain.

## THE RULE / LOGIC
Ensures the Ruler maintains a consistent state of the game, managing transitions and input validation.

**Lifecycle Phases:**
1. **Creation:** `NewRuler` initializes the actor but DOES NOT start the message loop.
2. **Setup:** The creator (e.g., `ArenaBridge`) has exclusive direct access to the `GameState` (Grid, Entities, Controllers) while the actor is stopped. Direct mutations are only safe in this phase.
3. **Activation:** The creator must explicitly call `Start()` to begin the actor loop.
4. **Ownership:** Once `Start()` is called, the Ruler takes **True Ownership** of the `GameState`. Direct external access is strictly prohibited and will cause race conditions. All subsequent interactions must occur via message passing. [[api_ruler_methods]]

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[domain_ruler_state]]`
