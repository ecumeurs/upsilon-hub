---
id: ui_registration_reroll_limit
human_name: Reroll Limit
type: UI
layer: ARCHITECTURE
version: 1.0
status: STABLE
priority: 5
tags: []
parents: 
  - [[ui_registration]]
dependents: []
---
# Reroll Limit

## INTENT
The reroll action is hard-capped at 3 total uses per account creation session.

## THE RULE / LOGIC
Character Generation Flow: The reroll action is hard-capped at 3

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[ui_registration_reroll_limit]]`
