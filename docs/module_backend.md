---
id: module_backend
human_name: UpsilonBattle Backend Component
type: MODULE
version: 1.0
status: REVIEW
priority: CORE
tags: [backend, go, api, combat]
parents: []
dependents:
  - [[module_game]]
  - [[mech_board_generation]]
  - [[mech_action_economy]]
  - [[mech_initiative]]
---

# UpsilonBattle Backend Component

## INTENT
To govern and calculate all strict battle-related logic in an isolated Go application providing a JSON API.

## THE RULE / LOGIC
- Technology Stack: Must be implemented in Go.
- Contract Type: Must expose operations via a JSON API.
- Scope of Responsibility:
  - **Included:** Board generation math, Initiative evaluation, Combat math (HP resolution), validation of Turn Action Economy, and Delay constraints.
  - **Excluded:** Matchmaking queues or Pairing logic (this must be handled by the UI/Frontend orchestrator).
- Stateless Processing: State enforcement during an active battle must be verified against the secure source of truth defined by this service.

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[module_backend]]`
- **Test Names:** `TestBackendCombatMath`, `TestBackendJSONContract`

## EXPECTATION (For Testing)
- UI requests a fire action via JSON API -> Go backend validates rules, subtracts HP, adds step delay, and returns updated board state JSON.
