---
id: module_frontend
human_name: BattleUI Frontend Component
type: MODULE
version: 1.0
status: REVIEW
priority: CORE
tags: [frontend, laravel, vue, tailwind, ui]
parents: []
dependents:
  - [[ui_landing]]
  - [[ui_registration]]
  - [[ui_dashboard]]
  - [[ui_waiting_room]]
  - [[ui_board]]
  - [[ui_leaderboard]]
  - [[req_matchmaking]]
---

# BattleUI Frontend Component

## INTENT
To present the interactive TRPG client, manage user sessions, and handle matchmaking orchestration before passing battle state to the Backend.

## THE RULE / LOGIC
- Technology Stack: Must be implemented as a Laravel application serving Vue.js views styled strictly with Tailwind CSS.
- Scope of Responsibility:
  - Session and Authentication Management (JWT distribution and validation).
  - Orchestrating the queuing system (`req_matchmaking`) and pairing clients prior to game instantiation.
  - Rendering the board UI, dashboards, player stats, and leaderboards.
  - Creating character entities (`entity_character`) during registration and enforcing the 3x reroll constraint (`mech_character_reroll`).
- Integration Constraint: Once matchmaking is confirmed, the BattleUI acts as a client bridge submitting user moves to the `module_backend` JSON API.

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[module_frontend]]`
- **Test Names:** `TestFrontendMatchmakingOrchestration`, `TestVueComponentRenders`

## EXPECTATION (For Testing)
- User signs up -> Laravel handles auth -> Vue renders Board -> Laravel passes moves to Go API.
