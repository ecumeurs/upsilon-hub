---
id: data_persistence
human_name: PostgreSQL Database Persistence
type: DATA
layer: IMPLEMENTATION
version: 1.0
status: STABLE
priority: 5
tags: [database, postgresql, state]
parents: []
dependents:
  - [[entity_game_match]]
---
# PostgreSQL Database Persistence

## INTENT
Serve as the centralized, persistent source of truth for accounts, characters, and historical match statistics.

## THE RULE / LOGIC
- Technology Stack: Must be strictly deployed on PostgreSQL.
- Primary Entities Supported:
  - Users (authentication credentials, win/loss metrics, ratio calculation material).
  - Characters (HP, Movement, Attack, Defense stats linked to a User via player_id).
  - Game Matches (matches historical state, board_state caching, turn tracking).
  - Matchmaking Queues (active queues with JSON-based character selection).
- Integration Note: Since Laravel orchestrates authentication and Go orchestrates active combat, both services may require explicit interface schemas or distinct responsibility bounded contexts inside this PostgreSQL instance.

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[data_persistence]]`
- **Test Names:** `TestPostgresPlayerSchema`, `TestPostgresCharacterSchema`

## EXPECTATION (For Testing)
- Game Ends via Go API -> Service updates Player Win/Loss record in PostgreSQL -> Laravel queries updated stats for the Leaderboard.
