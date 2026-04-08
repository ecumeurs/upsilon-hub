---
id: api_matchmaking
human_name: Matchmaking API
type: API
layer: ARCHITECTURE
version: 1.0
status: DRAFT
priority: 5
tags: [matchmaking, queue, api]
parents:
  - [[api_laravel_gateway]]
  - [[api_standard_envelope]]
dependents: []
---
# Matchmaking API

## INTENT
To allow players to enter a queue and be matched with opponents.

## THE RULE / LOGIC
**Endpoints:**
- `POST /api/v1/matchmaking/join`: Add player to queue for a specific game mode.
- `GET /api/v1/matchmaking/status`: Poll current matchmaking/match status.
- `DELETE /api/v1/matchmaking/leave`: Remove player from any active queue.
- `GET /api/v1/match/stats/waiting`: Get number of players currently in queue.
- `GET /api/v1/match/stats/active`: Get number of active matches from the engine.

## TECHNICAL INTERFACE (The Bridge)
- **API Endpoint:** `/api/v1/matchmaking/*`
- **Code Tag:** `@spec-link [[api_matchmaking]]`
- **Related Issue:** `ISS-007`
- **Test Names:** `TestJoinQueue`, `TestLeaveQueue`, `TestMatchFinding`

## EXPECTATION (For Testing)
- Join -> Player ID and characters stored in database `matchmaking_pool` (or equivalent persistent store).
- Leave -> Entry removed from database.
- Two compatible entries in pool -> Call Go `arena/start` -> Broadcast `game.started` to both.
