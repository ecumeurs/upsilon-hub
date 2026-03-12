---
id: api_go_battle_action
human_name: UpsilonBattle Arena Action API
type: API
version: 1.0
status: DRAFT
priority: CORE
tags: [api, golang, battle, action]
parents:
  - [[api_go_battle_engine]]
  - [[api_standard_envelope]]
dependents: []
---

# UpsilonBattle Arena Action API

## INTENT
To allow players to perform actions (Move, Attack, Skill) within an active battle arena.

## THE RULE / LOGIC
**Endpoint:** `POST /internal/arena/{id}/action`

### Request (Wrapped in [[api_standard_envelope]])
- `id`: `string` (UUID) - The ID of the arena (passed in URL).
- `player_id`: `string` - The ID of the player performing the action.
- `entity_id`: `string` - The ID of the entity performing the action.
- `type`: `string` - The action type (e.g., `Move`, `Attack`).
- `target_coords`: `Array<Position>`
  - `x`: `int`, `y`: `int`

### Response (Wrapped in [[api_standard_envelope]])
- Returns `200 OK` with the updated entity state or a success message.
- If the action is invalid, returns `400 Bad Request` with `Success: false`.

## TECHNICAL INTERFACE (The Bridge)
- **API Endpoint:** `POST /internal/arena/:id/action`
- **Code Tag:** `@spec-link [[api_go_battle_action]]`
- **Go Handler:** `handler.HandleArenaAction`
- **Request Type:** `api.ArenaActionRequest`
- **Response Map:**
  - `rulermethods.ControllerAttackReply` -> `api.NewEntity(d.Entity)`
  - `rulermethods.ControllerMoveReply` -> `api.NewEntity(d.Entity)`
  - Default -> `stdmessage.DataNil{}`

## EXPECTATION (For Testing)
- Valid `ArenaActionRequest` -> Ruler processes action -> Returns `200 OK`.
- Action target out of range -> Returns `400 Bad Request`.
- Arena ID not found -> Returns `400 Bad Request`.
