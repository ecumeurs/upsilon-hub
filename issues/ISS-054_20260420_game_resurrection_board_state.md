# Issue: Game Resurrection from Board State

**ID:** `20260420_game_resurrection_board_state`
**Ref:** `ISS-054`
**Date:** 2026-04-20
**Severity:** Medium
**Status:** Open
**Component:** `upsilonapi/bridge`
**Affects:** `battleui`, `upsilonapi`

---

## Summary

The frontend needs a mechanism to attempt "game resurrection" from a persisted board state when the Go API crashes or when a player re-logins during an active match.

---

## Technical Description

### Background
Currently, if the Go API crashes, the game state in memory is lost. While matches are tracked in the database, the active "ruler" and "controller" actors are gone.

### The Problem Scenario
1. A match is ongoing.
2. The Go API process crashes or is restarted.
3. The player refreshes the frontend.
4. The frontend sees an active match in the DB but cannot connect to the Go API bridge for that match.

### Random Seeding Conflict
The current system calls `tools.Seed()` (global `rand.Seed`) every time a `Ruler` is created (`NewRuler`). In a "resurrection" scenario where multiple actors are re-instantiated simultaneously after a crash, this leads to:
1. **Global Contention**: Matches re-seeding each other's random stream.
2. **Determinism Loss**: Inability to exactly replay a match if the seed isn't isolated and persisted.
3. **Sequence Duplication**: Simultaneous re-seeding results in identical random outcomes for different matches.

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | Medium |
| Impact if triggered | High |
| Detectability | High |
| Current mitigant | None |

---

## Recommended Fix

**Short term:** Add a heartbeat from frontend to Go API to detect crashes early.  
**Medium term:** Implement a "resurrection" endpoint that can re-instantiate Ruler/Controller actors from a saved DB snapshot. Refactor `Ruler` to use an isolated `*rand.Rand` source seeded by `MatchID`.
**Long term:** Move toward a more resilient actor model where state is periodically checkpointed.

---

## References

- [ui_investigation.md](file:///home/bastien/work/upsilon/projbackend/ui_investigation.md)
- [communication.md](file:///home/bastien/work/upsilon/projbackend/communication.md)
