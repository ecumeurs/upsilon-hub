---
id: rule_character_create_character
human_name: Character creation base allocation
type: RULE
version: 1.1
status: STABLE
priority: CORE
tags: [character, creation]
parents: 
  - [[entity_character]]
dependents: []
---

# Character creation base allocation

## INTENT
Define the base attributes and initial point distribution for new characters.

## THE RULE / LOGIC
- **Base Attributes:** Every character starts with:
  - HP: 3
  - Attack: 1
  - Defense: 1
  - Movement: 1
- **Initial Allocation:** Exactly **4 additional points** MUST be dispatched across these 4 categories during character creation.
- **Total Points:** A starting character has exactly 10 total attribute points (6 base + 4 dispatched).

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[rule_character_create_character]]`
- **Test Names:** `TestInitialStatConsistency`
