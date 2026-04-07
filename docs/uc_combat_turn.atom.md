---
id: uc_combat_turn
human_name: Combat Turn Use Case
type: MODULE
layer: ARCHITECTURE
version: 1.0
status: STABLE
priority: 5
tags: []
parents: []
dependents:
  - [[uc_combat_turn_action_selection]]
  - [[uc_combat_turn_delay_costs_accumulation]]
  - [[uc_combat_turn_initiative_evaluation]]
  - [[uc_combat_turn_shot_clock_expiration]]
  - [[uc_combat_turn_shot_clock_management]]
  - [[uc_combat_turn_turn_ending]]
  - [[us_take_combat_turn]]
---
# Combat Turn Use Case

## INTENT
To aggregate the constituent rules of Combat Turn Use Case.

## THE RULE / LOGIC
Character turns are managed based on initiative and shot clock. During their turn, players can choose to Move, Attack, Pass, or **Forfeit** the match.

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[uc_combat_turn]]`
