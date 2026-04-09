---
id: api_websocket_game_events
human_name: "WebSocket Game Events"
type: API
layer: ARCHITECTURE
version: 1.0
status: STABLE
priority: 3
tags: [websocket, real-time, api]
parents:
  - [[api_go_webhook_callback]]
dependents: []
---

# WebSocket Game Events

## INTENT
To provide real-time updates to players regarding matchmaking status and active game state by forwarding events from the Go engine via Laravel Reverb.

## THE RULE / LOGIC
1. **User Channel**: `private-user.{user_id}`
   - Used for events specific to a player before they are assigned a match.
   - **Event**: `MatchFound` - Triggered when the matchmaking service pairs the player. Contains `match_id`.
2. **Arena Channel**: `private-arena.{match_id}`
   - Used for events shared by all participants in a specific match.
   - **Event**: `BoardUpdated` - Forwarded from `api_go_webhook_callback`. Contains full or delta board state.
   - **Event**: `GameStarted` - Forwarded when the engine confirms the match is active.

## TECHNICAL INTERFACE (The Bridge)
- **Broadcaster:** Laravel Reverb
- **Frontend Client:** Laravel Echo / Pusher JS
- **Code Tag:** `@spec-link [[api_websocket_game_events]]`
- **Related Events:** `App\Events\MatchFound`, `App\Events\BoardUpdated`

## EXPECTATION (For Testing)
- Webhook `game.started` received -> `BoardUpdated` broadcasted to `private-arena.{match_id}`.
- PVP Match successfully created -> `MatchFound` broadcasted to all matched `private-user.{user_id}` channels.
