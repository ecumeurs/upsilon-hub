# Issue: Security Risk: Lack of Match Participant Access Control

**ID:** `20260312_match_participant_access_control`
**Ref:** `ISS-018`
**Date:** 2026-03-12
**Severity:** Critical
**Status:** Open
**Component:** `app/Http/Controllers/API/GameController.php`
**Affects:** All active matches

---

## Summary

Currently, any authenticated user can attempt to act or view the state of ANY active match by providing its ID. There is no authorization check to ensure the user is an actual participant of the match they are interacting with.

---

## Technical Description

### Background
Endpoints under `/api/v1/game/{id}/...` should be restricted to the participants defined in the `GameMatch` record or the engine's internal state.

### The Problem Scenario
An authenticated User A (not in match X) calls `POST /api/v1/game/X/action`. The controller forwards the request to the engine without checking if User A is a player in Match X.

### Where This Pattern Exists Today
- `app/Http/Controllers/API/GameController.php:18-51`

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | Low |
| Impact if triggered | Critical (Privacy leak, game interference) |
| Detectability | Low |
| Current mitigant | Match IDs are UUIDs (security through obscurity) |

---

## Recommended Fix

**Short term:** Implement a check in `GameController` that queries the `GameMatch` (or its cache) to ensure `auth()->id()` is listed among the players.

**Medium term:** Use Laravel Policies (`MatchPolicy`) to authorize all actions on the `GameMatch` resource.

---

## References

- [GameController.php](file:///workspace/battleui/app/Http/Controllers/API/GameController.php)
- [entity_game_match.atom.md](file:///workspace/docs/entity_game_match.atom.md)
