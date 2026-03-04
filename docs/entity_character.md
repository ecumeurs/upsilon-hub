---
id: entity_character
human_name: Character Entity
type: ENTITY
version: 1.0
status: REVIEW
priority: CORE
tags: [character, entities]
parents:
  - [[module_game]]
dependents:
  - [[rule_progression]]
---

# Character Entity

## INTENT
Defines the baseline stat block and attributes of a playable character unit.

## THE RULE / LOGIC
- Every character has four primary attributes:
  - HP (Health Points): Survival capacity.
  - Movement: Distance in squares the character can travel in one action.
  - Attack: Numerical value for offensive capability.
  - Defense: Numerical value for defensive capability.
- New Character Roll:
  - A new character is allocated exactly 10 initial attribute points.
  - A minimum of 3 points MUST be allocated to HP.
  - The remaining 7 points are distributed randomly among HP, Movement, Attack, and Defense, ensuring a minimum of 1 point is allocated to each of these attributes.

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[entity_character]]`
- **Test Names:** `TestCharacterCreationStats`, `TestCharacterMinHP`

## EXPECTATION (For Testing)
- Generate a new character -> HP is at least 3 -> Sum of HP, Movement, Attack, Defense is exactly 10.
