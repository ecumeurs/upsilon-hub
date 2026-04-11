---
id: api_go_battle_engine
human_name: Go UpsilonBattle JSON API & Webhook Dispatcher
type: MODULE
layer: ARCHITECTURE
version: 1.1
status: STABLE
priority: 5
tags: [api, golang, rest, webhooks]
parents:
  - [[api_standard_envelope]]
dependents:
  - [[module_upsilonapi]]
  - [[api_go_battle_start]]
  - [[api_go_webhook_callback]]
  - [[battleui_upsilon_api_service]]
  - [[api_go_battle_action]]
---
# Go UpsilonBattle JSON API & Webhook Dispatcher

## INTENT
To define the external JSON boundary for UpsilonBattle, allowing the Laravel Gateway to instantiate arenas, proxy commands, and receive asynchronous state updates via webhooks.

## THE RULE / LOGIC
The Go Battle Engine API is composed of several specialized endpoints and a webhook dispatch system. All communications must follow the [[api_standard_envelope]].

### Core Components:
- **Arena Initialization:** [[api_go_battle_start]] (Receives the initial player roster).
- **In-Game Actions:** [[api_go_battle_action]]
- **State Notifications:** [[api_go_webhook_callback]] (Broadcasts the `BoardState`).

### Board State & Identification:
The engine acts as the **Single Source of Truth** for player identity during a match. The `BoardState` includes a `players` list documenting the ID-to-Nickname mapping for all participants, enabling clients to correctly identify and render entities without external lookups.

## TECHNICAL INTERFACE (The Bridge)
- **Base Path:** `/internal`
- **Port:** `8081`
- **Code Tag:** `@spec-link [[api_go_battle_engine]]`

## EXPECTATION (For Testing)
- All endpoints must return a [[api_standard_envelope]] with `success: true` on successful operations.
- All endpoints must handle invalid input by returning a [[api_standard_envelope]] with `success: false` and a descriptive message.
