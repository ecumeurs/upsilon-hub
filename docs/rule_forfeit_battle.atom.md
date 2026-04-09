---
id: rule_forfeit_battle
human_name: "Forfeit Battle Rule"
type: RULE
layer: ARCHITECTURE
version: 1.0
status: STABLE
priority: 5
tags: [combat, forfeit, resolution]
parents:
  - [[us_take_combat_turn]]
  - [[rule_team_mechanics]]
dependents: []
---

# Forfeit Battle Rule

## INTENT
To allow a player to concede a match, resulting in an immediate victory for the opposing side(s) and proper arena closure.

## THE RULE / LOGIC
- A player may declare "FORFEIT" at any time during their character's turn.
- **PvE Resolution (Single Player vs AI):**
  - If the human player forfeits, the battle arena is closed immediately.
  - The human player is marked as "DEFEATED".
- **PvP Resolution (Multi-Player):**
  - If a player forfeits, all entities belonging to that player's `TeamID` are considered to have surrendered.
  - The forfeiting team is marked as "DEFEATED".
  - Victory is handed to the remaining team(s) with active entities.
  - In a 2v2 scenario, the forfeit of any single player on a team covers the entire team.

## TECHNICAL INTERFACE (The Bridge)
- **API Endpoint:** `POST /internal/arena/:id/action` (Type: "FORFEIT")
- **Code Tag:** `@spec-link [[rule_forfeit_battle]]`
- **Related Issue:** `#ISS-003`

## EXPECTATION (For Testing)
- Player A (Team 1) forfeits -> System broadcasts `BattleEnd` with `WinnerControllerID` from Team 2 -> Arena closed.
- In PvE, Player A forfeits -> System broadcasts `BattleEnd` where Player A is not the winner.
