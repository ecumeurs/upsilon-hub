# Issue: Security Risk: Battle Action player_id Usurpation

**ID:** `20260312_action_player_id_usurpation`
**Ref:** `ISS-017`
**Date:** 2026-03-12
**Severity:** Critical
**Status:** Resolved
**Component:** `app/Http/Controllers/API/GameController.php`
**Affects:** Battle engine integrity

---

## Summary

The battle action proxy endpoint accepts a `player_id` directly from the request payload. This allows any authenticated user to send actions on behalf of another player (enemy or teammate) if they know the UUIDs involved, effectively allowing them to usurp control of the game.

---

## Technical Description

### Background
`GameController::action` should ensure that the action being sent is strictly for the authenticated user.

### The Problem Scenario
1. Attacker joins match.
2. Attacker observes Player B's `player_id` and `entity_id` via WebSocket or API.
3. Attacker sends `POST /api/v1/game/{match}/action` with Player B's `player_id`.
4. Laravel forwards this to the Go engine without validating that `player_id == auth()->id()`.

### Where This Pattern Exists Today
- `app/Http/Controllers/API/GameController.php:32-37`

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | Medium |
| Impact if triggered | Critical (Griefing, rank manipulation) |
| Detectability | Low — looks like a valid engine action |
| Current mitigant | None |

---

## Recommended Fix

**Short term:** Remove `player_id` from the API request validation. In the controller, automatically inject `$request->user()->id` as the `player_id` before forwarding to the `UpsilonApiService`.

**Medium term:** Implement a `BattleAuthorizationMiddleware` to verify the user's role in the match.

---

## References

- [GameController.php](file:///workspace/battleui/app/Http/Controllers/API/GameController.php)
- [api_battle_proxy.atom.md](file:///workspace/docs/api_battle_proxy.atom.md)
