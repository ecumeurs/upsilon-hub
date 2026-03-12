---
id: battleui_api_dtos
human_name: BattleUI API Data Transfer Objects
type: DATA
version: 1.0
status: DRAFT
priority: CORE
tags: [battleui, dto, api, types]
parents:
  - [[battleui_upsilon_api_service]]
  - [[api_go_battle_start]]
  - [[api_go_battle_action]]
---

# BattleUI API Data Transfer Objects

## INTENT
To provide strongly-typed representations of the JSON payloads exchanged with the Go Battle Engine, ensuring that Laravel's implementation matches the Go `api` package exactly.

## THE RULE / LOGIC
The following DTOs must mirror the Go types defined in `api/input.go` and `api/output.go`:

### Request DTOs (Outgoing)
- **ArenaStartRequest:** Mirrors `api.ArenaStartRequest`.
    - `match_id`: string
    - `callback_url`: string
    - `players`: Array<PlayerDTO>
- **ArenaActionRequest:** Mirrors `api.ArenaActionRequest`.
    - `player_id`: string
    - `entity_id`: string
    - `type`: string
    - `target_coords`: Array<PositionDTO>

### Response DTOs (Incoming)
- **ArenaStartResponse:** Mirrors `api.ArenaStartResponse`.
    - `arena_id`: string
    - `initial_state`: BoardStateDTO
- **ArenaActionResponse:** Mirrors `api.ArenaActionResponse`.
    - `status`: string

## TECHNICAL INTERFACE (The Bridge)
- **Namespace:** `App\DTOs` or `App\Http\Resources`
- **Code Tag:** `@spec-link [[battleui_api_dtos]]`

## EXPECTATION (For Testing)
- Mapping a Go `BoardState` JSON to `BoardStateDTO` must not lose data.
- All DTOs must be serializable to JSON in a format accepted by the Go `gin` handlers.
