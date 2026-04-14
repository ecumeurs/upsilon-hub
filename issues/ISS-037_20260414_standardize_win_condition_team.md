# Issue: Standardize Win Condition: winner_team over winner_is_self

**ID:** `20260414_standardize_win_condition_team`
**Ref:** `ISS-037`
**Date:** 2026-04-14
**Severity:** Medium
**Status:** Open
**Component:** `battleui/app/Http/Resources`, `battleui/resources/js`, `upsiloncli`
**Affects:** `battleui/app/Http/Resources/BoardStateResource.php`, `battleui/resources/js/services/tactical.js`, `upsiloncli/upsilon_log_parser.py`

---

## Summary

The current win condition detection in the frontend and CLI relies on a server-computed `winner_is_self` boolean. This is brittle and hides the underlying team-based victory state. We should standardize on using a `winner_team` (or `winner_team_id`) field and ensure that the frontend components and the CLI log parser are aware of the user's own team ID to determine victory status locally.

---

## Technical Description

### Background
In `BoardStateResource.php`, the backend calculates if the winning player is the current user. While convenient, this prevents the UI from showing which team won if the user is a spectator or in multi-user team scenarios.

### The Problem Scenario
1. A 2v2 match concludes.
2. The engine reports a winning player ID.
3. `BoardStateResource.php` masks this ID into `winner_is_self: true/false`.
4. The frontend doesn't actually know which team won, only if the user won.
5. In the CLI log parser, we try to detect the `winner_team` from the raw board, but it's inconsistent because the field naming isn't standardized between WebSocket events and REST API responses (often `winner_id` vs `winner_team_id`).

### Where This Pattern Exists Today
- **Laravel Resource:** `BoardStateResource.php` (Lines 41-47: computes `winner_is_self` but unsets `winner_id`).
- **Vue.js Service:** `tactical.js` (Lacks a `myTeam()` helper that consistently resolves the team ID from the state).
- **CLI Parser:** `upsilon_log_parser.py` (Line 136-138: attempts to get `winner_team_id` from the payload).

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | Medium |
| Impact if triggered | Medium — UI/CLI inconsistency in victory reporting. |
| Detectability | High — Bots may fail to recognize a team-mate's win as their own without `winner_is_self`. |
| Current mitigant | `winner_is_self` boolean helper. |

---

## Recommended Fix

**Short term:** 
- Update `BoardStateResource.php` to include `winner_team_id` in the masked data.
- Ensure every player object in the board state has a `team` (already partially done).

**Medium term:** 
- Add a `myTeam(gameState)` helper to `tactical.js` and equivalent in the CLI.
- Change win detection logic from `if (winner_is_self)` to `if (state.winner_team === myTeam)`.

**Long term:** 
- Remove `winner_is_self` entirely once all clients are transitioned to team-based detection.

---

## References

- [BoardStateResource.php](file:///workspace/battleui/app/Http/Resources/BoardStateResource.php)
- [tactical.js](file:///workspace/battleui/resources/js/services/tactical.js)
- [upsilon_log_parser.py](file:///workspace/upsiloncli/upsilon_log_parser.py)
- [ISS-031](file:///workspace/issues/ISS-031_20260412_winner_id_missing.md)
