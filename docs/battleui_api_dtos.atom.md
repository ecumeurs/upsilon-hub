---
id: battleui_api_dtos
human_name: BattleUI API Data Transfer Objects
type: DATA
layer: IMPLEMENTATION
version: 1.0
status: STABLE
priority: 5
tags: [battleui, dto, api, types]
parents:
  - [[battleui_upsilon_api_service]]
  - [[api_go_battle_start]]
  - [[api_go_battle_action]]
dependents: []
---
# BattleUI API Data Transfer Objects

## INTENT
To provide strongly-typed representations of the JSON payloads exchanged with the Go Battle Engine, ensuring that Laravel's implementation matches the Go `api` package exactly.

## THE RULE / LOGIC
The following DTOs must mirror the Go types defined in `api/input.go` and `api/output.go` exactly.

### Outgoing Requests
- **ArenaStartRequest**
    - `match_id`: `string (UUID)`
    - `callback_url`: `string`
    - `players`: `Array<PlayerDTO>`
- **ArenaActionRequest**
    - `player_id`: `string (UUID)`
    - `entity_id`: `string (UUID)`
    - `type`: `string` ("MOVE", "ATTACK", "PASS", "FORFEIT")
    - `target_coords`: `Array<PositionDTO>`

### Incoming Responses
- **ArenaStartResponse**
    - `arena_id`: `string (UUID)`
    - `initial_state`: `BoardStateDTO`
- **ArenaActionResponse**
    - `status`: `string` ("accepted" | "rejected")

### Core Structures
- **PlayerDTO**: `{id: string (UUID), entities: Array<EntityDTO>, team: int, ia: boolean}`
- **EntityDTO**: `{id: string (UUID), player_id: string (UUID), name: string, hp: int, max_hp: int, attack: int, defense: int, move: int, max_move: int, position: PositionDTO}`
- **PositionDTO**: `{x: int, y: int}`
- **BoardStateDTO**: `{entities: Array<EntityDTO>, grid: GridDTO, turn: Array<TurnDTO>, current_player_id: string (UUID), current_entity_id: string (UUID), timeout: string (ISO8601), start_time: string (ISO8601), winner_id: string (UUID)|null}`
- **GridDTO**: `{width: int, height: int, cells: Array<Array<CellDTO>>}`
- **CellDTO**: `{entity_id: string (UUID)|null, obstacle: boolean}`
- **TurnDTO**: `{player_id: string (UUID), entity_id: string (UUID), delay: int}`

## TECHNICAL INTERFACE (The Bridge)
- **Namespace:** `App\DTOs` or `App\Http\Resources`
- **Code Tag:** `@spec-link [[battleui_api_dtos]]`

## EXPECTATION (For Testing)
- Mapping a Go `BoardState` JSON to `BoardStateDTO` must not lose data.
- All DTOs must be serializable to JSON in a format accepted by the Go `gin` handlers.
