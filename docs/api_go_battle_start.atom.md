---
id: api_go_battle_start
human_name: UpsilonBattle Arena Start API
type: API
version: 1.0
status: DRAFT
priority: CORE
tags: [api, golang, battle, initialization]
parents:
  - [[api_go_battle_engine]]
  - [[api_standard_envelope]]
dependents: []
---

# UpsilonBattle Arena Start API

## INTENT
To initialize a new battle arena instance with players, entities, and map data.

## THE RULE / LOGIC
**Endpoint:** `POST /internal/arena/start`

### Request (Wrapped in [[api_go_std_message]])
- `match_id`: `string` - Unique identifier for the match.
- `callback_url`: `string` - URL where webhook events will be sent.
- `players`: `Array<Player>`
  - `id`: `string` - Player controller ID.
  - `team`: `int` - Team number.
  - `ia`: `boolean` - Whether the player is AI-controlled.
  - `entities`: `Array<Entity>`
    - `id`: `string` - Unique entity ID.
    - `name`: `string` - Entity name.
    - `hp`: `int`, `max_hp`: `int`
    - `attack`: `int`, `defense`: `int`
    - `move`: `int`, `max_move`: `int`

### Response (Wrapped in [[api_go_std_message]])
- `arena_id`: `string` - The UUID of the newly created arena.
- `initial_state`: `BoardState`
  - `grid`: `Grid` (width, height, cells array)
  - `entities`: `Array<Entity>` (including positions)
  - `turn`: `Array<Turn>` (turn sequence)
  - `current_player_id`: `string`
  - `current_entity_id`: `string`

## TECHNICAL INTERFACE (The Bridge)
- **API Endpoint:** `POST /internal/arena/start`
- **Code Tag:** `@spec-link [[api_go_battle_start]]`
- **Go Handler:** `handler.HandleArenaStart`
- **Request Type:** `api.ArenaStartRequest`
- **Response Type:** `api.ArenaStartResponse`

## EXPECTATION (For Testing)
- Valid `ArenaStartRequest` -> Returns `200 OK` with `ArenaStartResponse`.
- Invalid JSON or missing required fields -> Returns `400 Bad Request` with `Success: false`.
