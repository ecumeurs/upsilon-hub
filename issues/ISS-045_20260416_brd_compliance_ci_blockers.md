# Issue: BRD Compliance CI Test Suite Blockers

**ID:** `20260416_brd_compliance_ci_blockers`
**Ref:** `ISS-045`
**Date:** 2026-04-16
**Severity:** High
**Status:** Open
**Component:** `upsiloncli`, `battleui`
**Affects:** CI Pipeline, BRD Compliance Matrix

---

## Summary

The implementation of automated BRD compliance tests via specialized CLI bot scripts is partially complete but currently blocked by environment constraints and synchronization issues. Specifically, the lack of a local `php` binary in the host environment (outside the devcontainer) hinders direct backend orchestration, and WebSocket event synchronization fails to trigger match start events reliably in the CLI agents.

---

## Technical Description

### Background
The project aims to achieve 100% automated validation of Business Requirements Document (BRD) rules through the `upsiloncli` farm. This involves spinning up multiple JS-controlled agents that simulate real player behavior and verify server-side constraints (e.g., password policy, progression locks, match resolution).

### The Problem Scenario
1. **Environment Lock**: Orchestration scripts (`test_match_resolution.sh`) and certain cleanup tasks expect a `php` binary for local artisan/database interactions, which is missing on the host.
2. **WebSocket Sync**: Even with the fix for `REVERB_PORT` (8080), agents frequently time out waiting for the `match.found` event. This prevents the "Match Resolution" tests from progressing to the tactical phase.
3. **Ghost Queue Entries**: Deleted bot accounts leave entries in `matchmaking_queues`, causing subsequent runs to fail with `500` errors (partially mitigated by a recent patch in `MatchMakingController.php`).

### Where This Pattern Exists Today
- `CI_INTEGRATION.md`: Contains the incomplete compliance matrix.
- `upsiloncli/tests/test_match_resolution.sh`: Fails due to timeouts.
- `upsiloncli/samples/progression_check.js`: Partially implemented; needs full verification drive.

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | High |
| Impact if triggered | Medium (Blocks CI) |
| Detectability | High — scripts explicitly fail with timeouts or 500s |
| Current mitigant | Manual verification of some backend logic |

---

## Recommended Fix

**Short term:** 
- Update `MatchMakingController` to more aggressively purge stale queue entries.
- Add retry logic to `upsiloncli` WebSocket listener for authorization failures.
- Use `npx` or container-wrapped commands for any PHP/Artisan tasks in the host.

**Medium term:** 
- Complete the JS implementation for GDPR (Data Portability) and Leaderboard testing.
- Integrate the `python3 upsilon_log_parser.py` analysis into the main CI YAML.

**Long term:** 
- Align the host environment with the devcontainer to ensure `php` and other dependencies are universally available for test runners.

---

## References

- [CI_INTEGRATION.md](CI_INTEGRATION.md)
- [test_match_resolution.sh](upsiloncli/tests/test_match_resolution.sh)
- [MatchMakingController.php](battleui/app/Http/Controllers/API/MatchMakingController.php)
- [ISS-043]: For Friendly Fire dependency
