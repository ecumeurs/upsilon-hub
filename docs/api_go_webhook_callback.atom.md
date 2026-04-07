---
id: api_go_webhook_callback
human_name: UpsilonBattle Webhook Callback
type: API
layer: ARCHITECTURE
version: 1.0
status: DRAFT
priority: 5
tags: [api, golang, callback, webhooks]
parents:
  - [[api_go_battle_engine]]
  - [[api_standard_envelope]]
dependents: []
---
# UpsilonBattle Webhook Callback

## INTENT
To asynchronously notify the Laravel Gateway of state changes, turn start/end, and battle results.

## THE RULE / LOGIC
**Destination:** The `callback_url` provided during [[api_go_battle_start]].

### Payload Structure (Wrapped in [[api_standard_envelope]])
- `event_type`: `string` - The type of event.
- `data`: `object` - Event-specific payload.

### Event Types mapping:
- `game.started` <- `rulermethods.BattleStart`
- `turn.started` <- `rulermethods.ControllerNextTurn`
- `game.ended` <- `rulermethods.BattleEnd`
- `board.updated` <- `rulermethods.EntitiesStateChanged`
- `attacked` <- `rulermethods.ControllerAttacked`

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[api_go_webhook_callback]]`
- **Go Dispatcher:** `bridge.HTTPController.forwardToWebhook`
- **Payload Type:** `api.ArenaEvent` (in some paths) or `map[string]interface{}` (in `forwardToWebhook`).

## EXPECTATION (For Testing)
- Ruler broadcasts `BattleStart` -> Dispatcher sends `POST` to `callback_url` with `event_type: "game.started"`.
- Dispatcher should handle non-200 responses from the callback URL (though current implementation just logs it).
