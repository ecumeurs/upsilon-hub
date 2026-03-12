# Issue: Missing Database Models in BattleUI

**ID:** `20260312_battleui_db_models`
**Ref:** `ISS-014`
**Date:** 2026-03-12
**Severity:** High
**Status:** Resolved
**Component:** `battleui/database`
**Affects:** `battleui` frontend, `upsilonapi` integration

---

## Summary

The battleui database is currently empty or lacking essential models needed for communication with the upsilon API. We specifically need a `characters` table linked to users (enforcing exactly 3 characters per user) to feed the API at battle start, and a Postgres `matches` table to cache game state, grid, and turn information for the frontend to avoid polling the upsilon API unnecessarily. 

---

## Technical Description

### Background
The `battleui` acts as the interface layer between the client frontend and the Upsilon game engine (`upsilonapi`). To start and manage an active arena (as described by `api_go_battle_action`), `battleui` must provide valid player entities and maintain an up-to-date representation of the match state to serve frontend requests efficiently.

### The Problem Scenario
1. A user attempts to initiate a battle via the frontend.
2. `battleui` must construct a payload containing the user's characters and pass it to `upsilonapi` to start the arena.
3. Because the `characters` table is empty or missing, `battleui` has no entities to send, blocking match creation.
4. During an active match, the frontend frequently requests the current game state (turn, grid geometry, positions).
5. Without a local cache, `battleui` must proxy every request to `upsilonapi`.
6. Since a Redis cache is explicitly rejected, `battleui` must use a true Postgres `matches` table to store this match state, but this table does not exist.

### Where This Pattern Exists Today
This gap exists in the `battleui` schema initialization and data models. The architectural concept is vaguely present in `/workspace/db.md` (e.g., `characters` and `match_history`), but the specific battle state cache (`matches` table) and strict enforcement are missing in the actual `battleui` implementation. Userstories and Usecases specs in ATDs have been split too much, leading to poorly defined data requirements at this stage.

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | High |
| Impact if triggered | High |
| Detectability | High — Match initiation straight-up fails; frontend state queries overwhelm API. |
| Current mitigant | None |

---

## Recommended Fix

**Short term:** Write the database migration for `characters` (linked to `players`, max 3) and `matches` (containing JSON or structured columns for game state, grid, turn).
**Medium term:** Implement the `battleui` repository layer to read/write to these tables, using `matches` as a cache during an active game.
**Long term:** Consolidate ATD Userstories and Usecases to provide clearer, holistic data requirements.

---

## References

- `/workspace/db.md` (TRPG Database Schema)
- `api_go_battle_action.atom.md`
- [ISS-013](file:///workspace/issues/ISS-013_20260312_battleui_api_service.md) (BattleUI API Communication Service Ownership)
