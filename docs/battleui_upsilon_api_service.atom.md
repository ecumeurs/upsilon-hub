---
id: battleui_upsilon_api_service
human_name: BattleUI UpsilonAPI Service
type: SERVICE
layer: ARCHITECTURE
version: 1.0
status: STABLE
priority: 5
tags: [battleui, service, api, integration]
parents:
  - [[api_go_battle_engine]]
dependents:
  - [[battleui_api_dtos]]
---
# BattleUI UpsilonAPI Service

## INTENT
To centralize and manage all communication between the Laravel-based `battleui` gateway and the Go-based `upsilonapi` (Battle Engine). This service ensures that all outgoing requests and incoming responses are validated against the defined API contracts.

## THE RULE / LOGIC
- **Ownership:** The `UpsilonApiService` is the sole owner of the HTTP client configuration and endpoint resolution for the `upsilonapi`.
- **Envelope Adherence:** All communications must be wrapped/unwrapped using the [[api_standard_envelope]].
- **DTO Mapping:** The service must use [[battleui_api_dtos]] for all request payloads and response unpacking.
- **Error Handling:** Any non-standard response or `success: false` from the Go engine must be translated into meaningful Laravel exceptions or logged with sufficient context.

## TECHNICAL INTERFACE (The Bridge)
- **Class:** `App\Services\UpsilonApiService`
- **Spec Links:**
    - [[api_go_battle_start]] -> `startArena(ArenaStartRequest $dto)`
    - [[api_go_battle_action]] -> `sendAction(ArenaActionRequest $dto)`
- **Code Tag:** `@spec-link [[battleui_upsilon_api_service]]`

## EXPECTATION (For Testing)
- `startArena` returns a valid `ArenaStartResponse` DTO or throws an exception.
- `sendAction` returns a valid `ArenaActionResponse` DTO or throws an exception.
- The service correctly handles timeouts and connection failures to the Go engine.
