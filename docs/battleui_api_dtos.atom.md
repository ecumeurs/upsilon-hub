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
  - [[api_go_battle_action]]
  - [[api_go_battle_start]]
  - [[battleui_upsilon_api_service]]
dependents: []
---
# BattleUI API Data Transfer Objects

## INTENT
To provide strongly-typed representations of the JSON payloads exchanged with the Go Battle Engine, ensuring that Laravel's implementation matches the Go `api` package exactly.

## THE RULE / LOGIC
The following DTOs represent the **Secure Client-Facing** representations. Laravel acts as a Masking Gateway, stripping internal UUIDs before transmission.

### Outgoing Requests (Frontend -> Laravel)
- **ArenaActionRequest**
    - `entity_id`: `string (UUID)`
    - `type`: `string` ("move", "attack", "pass", "forfeit")
    - `target_coords`: `Array<PositionDTO>`
    - *Note: `player_id` is automatically injected/verified by Laravel.*

### Core Structures (Masked for Privacy)
- **UserDTO**
    - `account_name`: `string`
    - `role`: `string`
    - `ws_channel_key`: `string (UUID)` (Pseudonym)
    - `email`: `string`
    - `total_wins`: `integer`
    - `ratio`: `string`
- **LeaderboardEntryDTO**
    - `account_name`: `string`
    - `wins`: `integer`
    - `losses`: `integer`
    - `score`: `float`
    - `rank`: `integer`
    - `is_self`: `boolean`
- **PlayerDTO**
    - `is_self`: `boolean`
    - `nickname`: `string`
    - `team`: `int`
    - `entities`: `Array<EntityDTO>`
- **EntityDTO**: `{id: string (UUID), is_self: boolean, name: string, hp: int, max_hp: int, attack: int, defense: int, move: int, max_move: int, position: PositionDTO}`
- **BoardStateDTO**: `{entities: Array<EntityDTO>, grid: GridDTO, turn: Array<TurnDTO>, is_my_turn: boolean, current_entity_id: string (UUID), timeout: string (ISO8601), start_time: string (ISO8601), is_winner: boolean|null, players: Array<PlayerDTO>}`
- **TurnDTO**: `{is_self: boolean, entity_id: string (UUID), delay: int}`

## TECHNICAL INTERFACE (The Bridge)
- **Namespace:** `App\DTOs` or `App\Http\Resources`
- **Code Tag:** `@spec-link [[battleui_api_dtos]]`

## EXPECTATION (For Testing)
- Mapping a Go `BoardState` JSON to `BoardStateDTO` must not lose data.
- All DTOs must be serializable to JSON in a format accepted by the Go `gin` handlers.
