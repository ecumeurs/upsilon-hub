---
id: mech_initiative
human_name: Initiative & Delay Mechanic
type: MECHANIC
version: 1.0
status: REVIEW
priority: CORE
tags: [combat, initiative]
parents:
  - [[module_game]]
dependents:
  - [[mech_action_economy]]
---

# Initiative & Delay Mechanic

## INTENT
Determines turn order mathematically based on action weight and randomly rolled starting values.

## THE RULE / LOGIC
- Game Startup Roll: At the very beginning of the match, every character rolls an initial initiative value ranging from `1` to `1000`.
- Active State: A character receives their active turn only when their evaluated initiative ticker reaches `0`.
- Delay Costs: Actions performed during an active turn (Pass, Move, Attack) incur a cumulative numeric Delay Cost. 
- Requeue Calculation: At the end of the turn, the character's required delay until their next turn is calculated using the summed Delay Cost of their performed actions (plus any penalties, see `mech_action_economy`).

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[mech_initiative]]`
- **Test Names:** `TestStartupInitiativeBounds`, `TestDelayCostAccumulation`

## EXPECTATION (For Testing)
- Turn ends after a Move and Attack -> Character Next Turn timer is set to `Base Delay + Delay(Move) + Delay(Attack)`.
