---
id: rule_friendly_fire
human_name: Friendly Immunity Rule
type: RULE
version: 1.0
status: REVIEW
priority: CORE
tags: [combat, rules]
parents:
  - [[module_game]]
dependents: []
---

# Friendly Immunity Rule

## INTENT
Prevents players and characters on the same side from inflicting harm upon one another.

## THE RULE / LOGIC
- Target Validation: Characters identified as belonging to the same team cannot apply destructive behavior (e.g., attacks, negative modifiers) to one another.
- This applies to characters controlled by the identical player and characters controlled by an allied player in a 2v2 match.

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[rule_friendly_fire]]`
- **Test Names:** `TestFriendlyFirePrevention`

## EXPECTATION (For Testing)
- Player attempts to attack an allied character -> Action is blocked/invalid.
