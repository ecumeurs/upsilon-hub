---
id: spec_match_format_win_condition_rule
human_name: Win Condition Rule
type: REQUIREMENT
layer: BUSINESS
version: 1.1
status: STABLE
priority: 5
tags: [combat, victory, team]
parents:
  - [[spec_match_format]]
dependents: []
---
# Win Condition Rule

## INTENT
Defines the win condition for a match

## THE RULE / LOGIC
- Victory is assigned to a Team via `WinnerTeamID` when all entities of all other teams have been defeated (HP <= 0) or have forfeited.
- Individual winners (`WinnerID`) are deprecated and must not be used to ensure player privacy and support team-based scoring.
- The `BattleEnd` message must explicitly carry the `WinnerTeamID`.

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[spec_match_format_win_condition_rule]]`
- **Test Names:** `TestVictoryStandardizationForfeit`
