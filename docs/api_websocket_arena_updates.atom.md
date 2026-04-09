---
id: api_websocket_arena_updates
human_name: "WebSocket Arena Updates (Private)"
type: API
layer: ARCHITECTURE
version: 1.0
status: STABLE
priority: 2
tags: [websocket, battle, tactical, updates]
parents:
  - [[api_websocket]]
  - [[api_battle_proxy]]
dependents: []
---

# WebSocket Arena Updates (Private)

## INTENT
To synchronize tactical game state and turn changes to all participants of a specific match in real-time.

## THE RULE / LOGIC
1. **Channel Name**: `private-arena.{match_id}`
   - `{match_id}` matches the UUID of the active match.
2. **Authorization**: Only participants assigned to the match can subscribe.
3. **Core Events**:
   - `board.updated`: Triggered by engine change.
     - **Payload**:
       - `match_id`: `string (UUID)`
       - `data`: `BoardStateDTO` (See [[battleui_api_dtos]])
   - `game.started`: Initial sync trigger.

## TECHNICAL INTERFACE (The Bridge)
- **Channel Pattern:** `private-arena.*`
- **Code Tag:** `@spec-link [[api_websocket_arena_updates]]`
- **Laravel Event:** `App\Events\BoardUpdated`

## EXPECTATION (For Testing)
- Game Action processed -> Engine Webhook hits BattleUI -> `board.updated` broadcasted.
- Client on match page -> Subscribed to `private-arena.{id}` -> Board state updates without refresh.
