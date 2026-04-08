---
id: entity_player
human_name: Player Account Entity
type: MODULE
layer: ARCHITECTURE
version: 1.0
status: STABLE
priority: 5
tags: []
parents: []
dependents:
  - [[entity_player_entity_character_rules_apply]]
  - [[uc_admin_user_management]]
  - [[entity_player_entity_player_stats_tracking]]
  - [[rule_admin_access_restriction]]
  - [[entity_player_entity_player_registration]]
  - [[entity_player_entity_player_initial_setup]]
  - [[infra_seed_admin]]
---
# Player Account Entity

## INTENT
To aggregate the constituent rules of Player Account Entity.

## THE RULE / LOGIC
Initial setup and registration for player accounts.
Core attributes for identity:
- `account_name` (Public/Unique)
- `full_address` (Private)
- `birth_date` (Private)
- `role` (Admin, Player)

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[entity_player]]`
