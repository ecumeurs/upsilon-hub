---
id: ui_character_roster
status: STABLE
type: UI
layer: ARCHITECTURE
priority: 5
tags: [ui, character, roster]
parents:
  - [[ui_dashboard]]
dependents: []
human_name: Character Roster Component
version: 1.0
---

# New Atom

## INTENT
To provide a standalone component for managing character rosters, stats, and progression.

## THE RULE / LOGIC
- Encapsulates character stat display and mutation logic.
- Implements frontend validation for `rule_progression`.
- Handles asynchronous state updates for character attributes.

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[ui_character_roster]]`
- **File:** [CharacterRoster.vue](file:///workspace/battleui/resources/js/Components/CharacterRoster.vue)
- **API:** `GET /v1/profile/characters`, `POST /v1/profile/character/{id}/reroll`, `POST /v1/profile/character/{id}/upgrade`

## EXPECTATION
- Must display all characters in the player's roster.
- Must provide buttons for Reroll (if wins == 0) and Progression (if points available).
- Must enforce movement allocation rules (1 per 5 wins).
- Must sync stat updates with the backend and refresh the UI.
