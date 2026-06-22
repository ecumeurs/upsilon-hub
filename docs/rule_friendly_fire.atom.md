---
id: rule_friendly_fire
human_name: Friendly Immunity Rule
type: RULE
layer: ARCHITECTURE
version: 1.0
status: STABLE
priority: 5
tags: []
parents:
  - [[spec_match_format]]
dependents: []
---
# Friendly Immunity Rule

## INTENT
Same-team entities cannot apply destructive behavior to one another.

## THE RULE / LOGIC
Prohibits self-inflicted harm between entities on the same side. Acceptance criteria:
- Entities identified as belonging to the same team cannot apply destructive behavior (attacks, negative modifiers) to one another.
- This protection covers characters controlled by the same player and, in a 2v2 match, characters controlled by an allied player.

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[rule_friendly_fire]]`
- **Implementation:** `upsilonbattle/battlearena/ruler/rules/attack.go`, `upsilonbattle/battlearena/ruler/rules/skill.go`
