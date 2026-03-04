---
id: spec_match_format
human_name: Match Format Specification
type: SPECIFICATION
version: 1.0
status: REVIEW
priority: CORE
tags: [game, format]
parents:
  - [[module_game]]
dependents: []
---

# Match Format Specification

## INTENT
Defines the structural sizes and win conditions for a game session.

## THE RULE / LOGIC
- Supported Modes: Matches must be strictly formatted as either 1v1 (One Player vs One Player/AI) or 2v2 (Two Players vs Two Players/AI).
- Team Composition: Each active player independently controls exactly 3 characters on the board.
- Win Condition: Victory is assigned to a team when the opposing team's characters have all been overcome (HP reduced to 0).

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[spec_match_format]]`
- **Test Names:** `TestMatchWinCondition`, `TestTeamCompositionLimit`

## EXPECTATION (For Testing)
- 1v1 match starts -> Each side has 3 characters.
- 2v2 match starts -> Each side has 6 characters (3 per player).
