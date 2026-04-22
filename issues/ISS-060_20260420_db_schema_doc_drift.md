# Issue: Database Schema and Documentation Drift

**ID:** `20260420_db_schema_doc_drift`
**Ref:** `ISS-060`
**Date:** 2026-04-20
**Severity:** Medium
**Status:** Resolved
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

## Investigation Results (2026-04-22)

### Summary
The investigation revealed that this was primarily a **documentation drift** issue rather than an implementation problem. The actual PostgreSQL schema is correct and functional, but the `db.md` documentation contained outdated information.

### Critical Discrepancies Found & Fixed

#### 1. Users Table
- **Removed:** `email_verified_at` field (doesn't exist in actual schema)
- **Added:** Missing `role` field with default 'Player' value
- **Added:** `updated_at` index and role constraint documentation
- **Fixed:** `ws_channel_key` documented as nullable (was incorrectly shown as required)

#### 2. Match Participants Table
- **Added:** `player_id` is **nullable** (critical for AI/bot support in PvE modes)
- **Added:** CHECK constraint for status field (WIN/LOSS only)
- **Added:** Comprehensive notes about human vs AI participants

#### 3. General Improvements
- **Updated:** Mermaid ERD to accurately reflect actual schema
- **Added:** Documentation for CHECK constraints and indexes
- **Added:** Acknowledgment of Laravel system tables (cache, jobs, etc.)

### Documentation Updates Completed

1. **Updated `db.md`:**
   - Corrected users table structure
   - Fixed match_participants nullable field documentation
   - Updated Mermaid ERD with proper constraints

2. **Created New ATD Atoms:**
   - `docs/entity_users.atom.md` - Comprehensive user entity specification
   - `docs/entity_match_participants.atom.md` - Complete match participation specification

3. **Updated Existing ATD Atoms:**
   - `docs/data_persistence.atom.md` - Enhanced with complete entity listing
   - `docs/entity_game_match.atom.md` - Updated dependencies

### Action Plan
See `ISS-060_action_plan.md` for comprehensive remediation strategy including:
- Phase 1: Core documentation updates (✅ Completed)
- Phase 2: ATD atom enhancements (✅ Completed)
- Phase 3: Verification & testing (⏳ Pending)
- Phase 4: Automation implementation (⏳ Pending)

### Risk Assessment
**Current Risk:** 🟢 **LOW**
- No code changes required
- No breaking changes to existing functionality
- Schema is production-ready

**Business Impact:** 🟢 **POSITIVE**
- Improved developer onboarding experience
- Reduced confusion about schema structure
- Better support for AI/bot participants in PvE modes
- Enhanced ATD compliance

## References

- [db.md](file:///workspace/db.md) - ✅ Updated to match actual schema
- [entity_users.atom.md](file:///workspace/docs/entity_users.atom.md) - ✅ Created
- [entity_match_participants.atom.md](file:///workspace/docs/entity_match_participants.atom.md) - ✅ Created
- [entity_game_match.atom.md](file:///workspace/docs/entity_game_match.atom.md) - ✅ Updated
- [data_persistence.atom.md](file:///workspace/docs/data_persistence.atom.md) - ✅ Updated
- [ISS-060_action_plan.md](file:///workspace/issues/ISS-060_action_plan.md) - Comprehensive remediation plan
- [ISS-060_investigation_completed.md](file:///workspace/issues/ISS-060_investigation_completed.md) - Investigation summary
