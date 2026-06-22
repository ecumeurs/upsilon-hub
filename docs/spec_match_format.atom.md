---
id: spec_match_format
human_name: Match Format Specification
type: REQUIREMENT
layer: BUSINESS
version: 1.0
status: STABLE
priority: 5
tags: []
parents: []
dependents:
  - [[rule_friendly_fire]]
---
# Match Format Specification

## INTENT
Defines the structural format of a match: its mode, per-player team size, and the team-based win condition.

## THE RULE / LOGIC
Defines the structural sizes and win condition for a game session. Acceptance criteria:
- **Mode:** A match is strictly either 1v1 (one player vs one player/AI) or 2v2 (two players vs two players/AI).
- **Team composition:** Each active player independently controls exactly 3 characters on the board.
- **Win condition:** Victory is assigned to a team via `WinnerTeamID` when all entities of every other team are defeated (HP <= 0) or have forfeited.
- Individual winners (`WinnerID`) are deprecated and must not be used, to preserve player privacy and support team-based scoring.
- The `BattleEnd` message must explicitly carry the `WinnerTeamID`.

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[spec_match_format]]`
