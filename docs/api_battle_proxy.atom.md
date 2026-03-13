---
id: api_battle_proxy
human_name: Battle Proxy & Webhook API
type: API
version: 1.0
status: DRAFT
priority: CORE
tags: [battle, proxy, webhook, api]
parents:
  - [[api_laravel_gateway]]
  - [[api_standard_envelope]]
dependents:
  - [[api_go_battle_action]]
  - [[api_go_webhook_callback]]
---

# Battle Proxy & Webhook API

## INTENT
To proxy user actions to the Go engine and ingest engine state updates back into the Laravel ecosystem.

## THE RULE / LOGIC
**Endpoints:**
- `POST /api/v1/game/{id}/action`: User action proxy.
- `POST /api/webhook/upsilon`: Engine callback ingestion.

### Action Proxy Logic:
1. Validate request token.
2. Verify player belongs to the arena.
3. Forward payload to Go `POST /internal/arena/{id}/action`.
4. Return Go's response directly or mapped.

### Webhook Ingestion Logic:
1. Receive state update from Go.
2. Update Redis key `arena:{id}:state`.
3. Broadcast via Reverb to channel `arena.{id}`.

## TECHNICAL INTERFACE (The Bridge)
- **API Endpoint:** `/api/v1/game/*`, `/api/webhook/*`
- **Code Tag:** `@spec-link [[api_battle_proxy]]`
- **Related Issue:** `ISS-007`
- **Test Names:** `TestActionProxying`, `TestWebhookUpdatesStateAndBroadcasts`

## EXPECTATION (For Testing)
- Action forwarded to Go carries the same `request_id`.
- Webhook receipt triggers `BoardUpdated` event in Laravel.
