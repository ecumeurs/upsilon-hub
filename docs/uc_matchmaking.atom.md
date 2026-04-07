---
id: uc_matchmaking
human_name: Matchmaking & Queue Use Case
type: MODULE
layer: ARCHITECTURE
version: 1.0
status: STABLE
priority: 5
tags: []
parents: []
dependents:
  - [[uc_matchmaking_match_start]]
  - [[uc_matchmaking_matchmaking]]
  - [[uc_matchmaking_pve]]
  - [[uc_matchmaking_pvp]]
  - [[uc_matchmaking_redirect_to_board]]
---
# Matchmaking & Queue Use Case

## INTENT
To aggregate the constituent rules of Matchmaking & Queue Use Case.

## THE RULE / LOGIC
End-to-end narrative of a logged-in player selecting a game mode and transitioning to the board.

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[uc_matchmaking]]`
