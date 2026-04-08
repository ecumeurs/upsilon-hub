---
id: rule_team_mechanics
human_name: "Battle Team Mechanics"
type: RULE
layer: ARCHITECTURE
version: 1.0
status: DRAFT
priority: 5
tags: [combat, team, logic]
parents:
  - [[us_take_combat_turn]]
dependents:
  - [[rule_forfeit_battle]]
---

# Battle Team Mechanics

## INTENT
To allow grouping entities into teams, enabling cooperative gameplay and team-based win conditions.

## THE RULE / LOGIC
- Each entity must possess a `TeamID` (Integer) property.
- All entities sharing the same `TeamID` are considered allies.
- All entities with different `TeamID`s are considered enemies.
- **Targeting Restriction:** Defensive or supportive actions target allies. Offensive actions must target enemies.
- **Victory Condition:** A team wins when all entities belonging to all other teams are defeated or have forfeited.

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[rule_team_mechanics]]`
- **Property Name:** `TeamID`
- **Related Issue:** `#ISS-003`
- **Related Tests:** `TestTeamBasedTargeting`, `TestTeamVictoryCondition`

## EXPECTATION (For Testing)
- Entity A (Team 1) attempts to attack Entity B (Team 1) -> Action rejected.
- Entity A (Team 1) attempts to attack Entity C (Team 2) -> Action accepted.
- All Team 2 entities defeated -> Team 1 is declared the winner.
