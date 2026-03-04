---
id: req_matchmaking
human_name: Matchmaking Flow Requirement
type: REQUIREMENT
version: 1.1
status: REVIEW
priority: CORE
tags: [matchmaking, network]
parents: []
dependents:
  - [[spec_match_format]]
---

# Matchmaking Flow Requirement

## INTENT
Provides simple avenues for players to find opponents or play against the system.

## THE RULE / LOGIC
- The system must offer four explicit matchmaking queue options:
  1. 1v1 PVE: Player vs Computer.
  2. 1v1 PVP: Player vs Player.
  3. 2V2 PVE: Two Players vs Computer.
  4. 2V2 PVP: Two Players vs Two Players.
- Transition: Players join the Waiting Room until the required human count is met, then instantly spawn onto the Board.

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[req_matchmaking]]`
- **Test Names:** `TestMatchmakingQueues`

## EXPECTATION (For Testing)
- Player selects 1v1 PVE -> Bypass Wait -> Starts 1v1 Match vs AI.
- Player selects 2v2 PVP -> Waits for 3 additional players to join.
