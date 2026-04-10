# Issue: AI enemies (for PVE) require database user records

**ID:** `20260409_ai_enemies_require_user`
**Ref:** `ISS-028`
**Date:** 2026-04-09
**Severity:** Medium
**Status:** Resolved
**Component:** `battleui/app/Http/Controllers/API`
**Affects:** `MatchMakingController.php`, `MatchParticipant` model

---

## Summary

Currently, the `match_participants` table has a foreign key constraint requiring every participant to have a valid record in the `users` table. This forces the system to create "dummy" AI users or use a system AI user to satisfy the database schema when starting 1v1 PVE matches.

---

## Technical Description

### Background
In a 1v1 PVE match, a human player is matched against an AI opponent. The AI opponent is a virtual entity managed by the Go engine.

### The Problem Scenario
When attempting to register an AI player as a participant in a match, the `MatchMakingController` must provide a `player_id`. If this ID does not exist in the `users` table, a `SQLSTATE[23503]: Foreign key violation` occurs.

### Where This Pattern Exists Today
- `MatchMakingController.php` (joinMatch logic)
- `database/migrations/..._create_match_participants_table.php`

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | High (every PVE match) |
| Impact if triggered | Medium (requires database clutter) |
| Detectability | High (SQL error) |
| Current mitigant | Hardcoded AI user in `users` table created via `firstOrCreate`. |

---

## Recommended Fix

**Short term:** Use a stable `AI_PLAYER_ID` and ensure it exists in the `users` table.
**Medium term:** Relax the foreign key constraint on `match_participants.player_id` to allow null or a reserved range for AI participants.
**Long term:** Refactor the participant system to separate "Account Owners" (Users) from "Tactical Controllers" (which can be AI or Users).

---

## References

- [MatchMakingController.php](file:///workspace/battleui/app/Http/Controllers/API/MatchMakingController.php)
- [MatchParticipant.php](file:///workspace/battleui/app/Models/MatchParticipant.php)
