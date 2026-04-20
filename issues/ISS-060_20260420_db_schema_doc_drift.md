# Issue: Database Schema and Documentation Drift

**ID:** `20260420_db_schema_doc_drift`
**Ref:** `ISS-060`
**Date:** 2026-04-20
**Severity:** Medium
**Status:** Open
**Component:** `docs/db.md`
**Affects:** `data_persistence`, `entity_game_match`

---

## Summary

There is a significant delta between the centralized database documentation (`db.md`), the ATD atoms (`entity_game_match`), and the actual PostgreSQL implementation. This drift creates confusion for developers and risks introducing bugs in the BattleUI and Go Engine integration.

---

## Technical Description

### Background
The project uses `db.md` as the primary human-readable reference for the database schema. Architectural requirements are captured in ATD atoms (e.g., `data_persistence`, `entity_game_match`). The implementation is managed via Laravel migrations and a PostgreSQL instance.

### The Problem Scenario
A comparison between `db.md` and the actual Postgres schema revealed the following discrepancies:
1.  **Table Naming**: `players` (doc) vs `users` (actual). `matchmaking_queue` (doc) vs `matchmaking_queues` (actual).
2.  **Column Naming**: `password_hash` (doc) vs `password` (actual). `winner_team_id` (doc) vs `winning_team_id` (actual).
3.  **State Schema**: `db.md` refers to a single `board_state` JSONB field, while the implementation uses `game_state_cache` and `grid_cache` (matching the `entity_game_match` atom, but not `db.md`).
4.  **Metadata Differences**: Missing Laravel-specific fields (`remember_token`, `email_verified_at`) and discrepancies in matchmaking queue attributes.

### Where This Pattern Exists Today
- `/workspace/db.md` (Multiple locations)
- `PostgreSQL: public.users`, `public.game_matches`, `public.matchmaking_queues`

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | High |
| Impact if triggered | Medium |
| Detectability | Medium — manifests as mismatched API payloads or failed database queries in development |
| Current mitigant | Defensive coding in the Go/PHP bridge, but documentation remains misleading |

---

## Recommended Fix

**Short term:** Update `db.md` and its Mermaid ERD to match the actual PostgreSQL schema exactly.
**Medium term:** Perform an ATD audit to ensure all `ENTITY` and `DATA` atoms use the implementation-proven naming conventions.
**Long term:** Automate `db.md` synthesis directly from the database schema or migration files to prevent future drift.

---

## References

- [db.md](file:///workspace/db.md)
- [entity_game_match.atom.md](file:///workspace/docs/entity_game_match.atom.md)
- [data_persistence.atom.md](file:///workspace/docs/data_persistence.atom.md)
