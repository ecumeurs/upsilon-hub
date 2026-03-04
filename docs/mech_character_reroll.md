---
id: mech_character_reroll
human_name: Character Reroll Mechanic
type: MECHANIC
version: 1.0
status: REVIEW
priority: CORE
tags: [character, creation]
parents:
  - [[entity_player]]
  - [[entity_character]]
dependents:
  - [[ui_registration]]
---

# Character Reroll Mechanic

## INTENT
Provides the player a limited chance to re-randomize their starting roster stats before committing to the account creation.

## THE RULE / LOGIC
- Availability: Allowed strictly during account creation after the initial 3 characters are generated.
- Limit: The player may trigger a "Reroll" action up to a strict maximum of exactly 3 times per account creation flow.
- Effect: A Reroll completely discards the current set of 3 generated characters and mathematically rolls 3 brand new ones (adhering to `entity_character` stats).

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[mech_character_reroll]]`
- **Test Names:** `TestRerollStatChanges`, `TestRerollMaxLimit`

## EXPECTATION (For Testing)
- Player clicks reroll -> All 3 characters have new stats.
- Player clicks reroll 3 times -> Reroll counter hits limit -> 4th attempt is rejected.
