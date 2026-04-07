---
id: uc_match_resolution
human_name: Match Resolution Use Case
type: MODULE
layer: ARCHITECTURE
version: 1.0
status: STABLE
priority: 5
tags: []
parents: []
dependents:
  - [[uc_match_resolution_dashboard_redirect]]
  - [[uc_match_resolution_match_persistence]]
  - [[uc_match_resolution_progression_reward]]
  - [[uc_match_resolution_win_detection]]
---
# Match Resolution Use Case

## INTENT
To handle the immediate conclusion of a game match, win detection, and backend reward persistence.

## THE RULE / LOGIC
1. Game match concludes when a win condition is met [[spec_match_format_win_condition_rule]] or a player **forfeits** [[uc_combat_turn]].
2. **Win Detection / Next Turn**: 
   - If a winner is identified: System persists match data [[uc_match_resolution_match_persistence]] and awards progression points [[uc_match_resolution_progression_reward]].
   - If **NO** winner is detected after a character turn: The flow returns to **Combat Turn Management (UC-4)** for the next character's turn.
3. Upon match conclusion, User is given the choice to enter **Progression & Stat Allocation (UC-6)** or return to the **Dashboard**.

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[uc_match_resolution]]`
