---
id: entity_player
human_name: Player Account Entity
type: ENTITY
version: 1.0
status: REVIEW
priority: CORE
tags: [player, account]
parents: []
dependents:
  - [[entity_character]]
---

# Player Account Entity

## INTENT
Defines the required player identity and initial setup.

## THE RULE / LOGIC
- Registration: Every player must connect through a logged-in account to play.
- Initial Setup: Upon account creation, the player is automatically granted exactly 3 characters.
- These characters must have their attributes rolled according to the rules defined in `entity_character`.
- Statistics Tracking: The player entity must track the absolute number of game wins and losses. A win/loss ratio is mathematically derived from these figures.

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[entity_player]]`
- **Test Names:** `TestPlayerAccountCreation`

## EXPECTATION (For Testing)
- Create new player account -> Player immediately has 3 instantiated characters with rolled stats in their roster.
- Fetch Player Stats -> Wins: 0, Losses: 0.
