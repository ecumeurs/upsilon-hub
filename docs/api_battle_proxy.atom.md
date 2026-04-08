---
id: api_battle_proxy
human_name: Battle Proxy & Webhook API
type: API
layer: ARCHITECTURE
version: 1.0
status: DRAFT
priority: 5
tags: [battle, proxy, webhook, api]
parents:
  - [[api_laravel_gateway]]
  - [[api_standard_envelope]]
dependents: []
---
# Battle Proxy & Webhook API

## INTENT
To proxy user actions to the Go engine and ingest engine state updates back into the Laravel ecosystem.

## THE RULE / LOGIC
**Endpoints:**
- `GET /api/v1/game/{match_id}`: Retrieve the cached board state/match details.
- `POST /api/v1/game/{match_id}/action`: Proxy tactical user commands to the Go engine.

### Response - Match Details (Wrapped in [[api_standard_envelope]])
- `id`: `string (UUID)`
- `game_mode`: `string`
- `started_at`: `string (ISO8601)`
- `concluded_at`: `string (ISO8601)|null`
- `winning_team_id`: `int|null`

### Request - Action (Wrapped in [[api_standard_envelope]])
- `payload`: `ArenaActionRequest` (See [[api_go_battle_action]])

### Response - Action (Wrapped in [[api_standard_envelope]])
- `status`: `string` ("accepted" | "rejected")

## TECHNICAL INTERFACE (The Bridge)
- **API Endpoint:** `/api/v1/game/*`, `/api/webhook/*`
- **Code Tag:** `@spec-link [[api_battle_proxy]]`
- **Related Issue:** `ISS-007`
- **Test Names:** `TestActionProxying`, `TestWebhookUpdatesStateAndBroadcasts`

## EXPECTATION (For Testing)
- Action forwarded to Go carries the same `request_id`.
- Webhook receipt triggers `BoardUpdated` event in Laravel.
