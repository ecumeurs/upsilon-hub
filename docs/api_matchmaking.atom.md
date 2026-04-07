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
- `POST /api/v1/matchmaking/join`: Add player to queue.
- `GET /api/v1/matchmaking/status`: Poll current matchmaking/match status.
- `DELETE /api/v1/matchmaking/leave`: Remove player from queue.

### Request - Join (Wrapped in [[api_standard_envelope]])
- `game_mode`: `string` - The game mode the player wants to join (e.g., "1v1_PVP", "2v2_PVP").

### Response - Join / Status (Wrapped in [[api_standard_envelope]])
- `status`: `string` ("queued", "matched", "idle")
- `match_id`: `uuid` (optional, for "matched")
- `expected_participants`: `int` (Total expected for the mode: 1, 2, or 4)
- `empty_slots`: `int` (Remaining slots)
- `queued_at`: `datetime` (optional, for "queued" status)

## TECHNICAL INTERFACE (The Bridge)
- **API Endpoint:** `/api/v1/matchmaking/*`
- **Code Tag:** `@spec-link [[api_matchmaking]]`
- **Related Issue:** `ISS-007`
- **Test Names:** `TestJoinQueue`, `TestLeaveQueue`, `TestMatchFinding`

## EXPECTATION (For Testing)
- Join -> Player ID and characters stored in database `matchmaking_pool` (or equivalent persistent store).
- Leave -> Entry removed from database.
- Two compatible entries in pool -> Call Go `arena/start` -> Broadcast `game.started` to both.
