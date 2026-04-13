# Issue: Internal ID Exposure in Public APIs

**ID:** `20260413_id_exposure`
**Ref:** `ISS-034`
**Date:** 2026-04-13
**Severity:** Medium
**Status:** Open
**Component:** `app/Http/Controllers/API`
**Affects:** `Leaderboard`, `Matchmaking`, `Profile`

---

## Summary

Internal database UUIDs are currently being emitted directly to front-end and CLI consumers. This exposes implementation details, increases payload size, and presents a reconnaissance risk for attackers targeting specific accounts or matches.

---

## Technical Description

### Background
The application uses UUIDs as primary keys for Users, matches, and characters. These are intended for internal orchestration.

### The Problem Scenario
Public APIs like `GET /api/v1/leaderboard` previously included the raw `id` (UUID) for every combatant. While useful for the front-end to highlight the "current user", it exposes the UUID of every other player.

### Where This Pattern Exists Today
- **API Responses:**
    - [LeaderboardController.php](file:///workspace/battleui/app/Http/Controllers/API/LeaderboardController.php) (Fixed locally, but pattern exists in others)
    - [MatchmakingQueueResource.php](file:///workspace/battleui/app/Http/Resources/MatchmakingQueueResource.php) -> `user_id`
    - [MatchParticipantResource.php](file:///workspace/battleui/app/Http/Resources/MatchParticipantResource.php) -> `player_id`
    - [CharacterResource.php](file:///workspace/battleui/app/Http/Resources/CharacterResource.php) -> `id`
    - [GameController.php](file:///workspace/battleui/app/Http/Controllers/API/GameController.php) -> `state()` returns `player_id` for all participants.
- **WebSocket Events:**
    - `MatchFound` -> `match_id`
    - `board.updated` -> Complete `BoardState` containing `player_id`, `entity_id`, `winner_id` for all actors.

### Relevance Investigation
- **Battle Engine (Go):** **REQUIRED.** The engine uses `player_id` (User UUID) and `entity_id` (Character UUID) as primary keys for state management and action validation. Laravel must continue sending these to the engine.
- **Frontend (Vue):** **PARTIALLY REQUIRED.**
    - Needs to highlight "Self" (can be replaced by `is_self: boolean`).
    - Needs to target entities for actions (requires a unique ID, but could be a masked `public_id` or `slug`).
    - Needs to subscribe to private channels (uses `ws_channel_key`, which is a pseudonym UUID).

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | High |
| Impact if triggered | Medium |
| Detectability | High — evident in any JSON response |
| Current mitigant | None (manual stripping in individual controllers) |

---

## Recommended Fix

**Short term:** Surgically unset `id` fields in Controller `map()` functions and replace with semantic identifiers (e.g., `is_self: boolean`).  
**Medium term:** Implement **Resources/DTOs** (Laravel API Resources) to strictly define the outgoing contract and automatically filter internal metadata.  
**Long term:** Use **HashIDs** or **Slugs** for public identification if a unique identifier is absolutely required by the client.

---

## References

- [LeaderboardController.php](file:///workspace/battleui/app/Http/Controllers/API/LeaderboardController.php)
- [battleui_api_dtos.atom.md](file:///workspace/docs/battleui_api_dtos.atom.md)
