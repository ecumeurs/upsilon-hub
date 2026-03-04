---
id: mech_action_economy
human_name: Turn Action Economy Mechanic
type: MECHANIC
version: 1.0
status: REVIEW
priority: CORE
tags: [combat, turn]
parents:
  - [[module_game]]
dependents:
  - [[mech_initiative]]
---

# Turn Action Economy Mechanic

## INTENT
Defines the allowable actions and temporal constraints for a character's active turn.

## THE RULE / LOGIC
- Action Economy Costs:
  - Move: `+20` delay cost per tile moved.
  - Attack: `+100` delay cost.
  - Pass: `+300` delay cost.
- Time Constraint: Turn duration is strictly capped at 30 seconds.
- Timeout Penalty: If a turn lasts exactly 30 seconds without completion, an automatic "Pass" is triggered, and a strict penalty of `+100` delay cost is added on top of the base Pass cost (Total `+400`).

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[mech_action_economy]]`
- **Test Names:** `TestActionLimits`, `TestTurnTimeoutPenalty`

## EXPECTATION (For Testing)
- Turn elapses 30 seconds -> Turn auto-ends -> Character receives +100 delay cost penalty.
