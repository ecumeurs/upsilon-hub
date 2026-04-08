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
- `GET /api/v1/matchmaking/status`: Poll current matchmaking/match status or reconnect to an active match.
- `DELETE /api/v1/matchmaking/leave`: Remove player from any active queue.

### Request - Join (Wrapped in [[api_standard_envelope]])
- `game_mode`: `string` ("1v1_PVP", "2v2_PVP", "1v1_PVE", "2v2_PVE")

### Response - Join / Status (Wrapped in [[api_standard_envelope]])
- `status`: `string` ("queued", "matched", "idle")
- `match_id`: `string (UUID)|null`
- `expected_participants`: `int|null`
- `empty_slots`: `int|null`
- `queued_at`: `string (ISO8601)|null`

## TECHNICAL INTERFACE (The Bridge)
- **API Endpoint:** `/api/v1/matchmaking/*`
- **Code Tag:** `@spec-link [[api_matchmaking]]`
- **Related Issue:** `ISS-007`
- **Test Names:** `TestJoinQueue`, `TestLeaveQueue`, `TestMatchFinding`

## EXPECTATION (For Testing)
- Join -> Player ID and characters stored in database `matchmaking_pool` (or equivalent persistent store).
- Leave -> Entry removed from database.
- Two compatible entries in pool -> Call Go `arena/start` -> Broadcast `game.started` to both.
