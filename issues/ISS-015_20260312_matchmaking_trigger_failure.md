# Issue: Matchmaking Trigger Failures in 2v2 and PVE Modes

**ID:** `20260312_matchmaking_trigger_failure`
**Ref:** `ISS-015`
**Date:** 2026-03-12
**Severity:** High
**Status:** Resolved
**Component:** `app/Http/Controllers/API/MatchMakingController.php`
**Affects:** `MatchmakingQueue`, `UpsilonApiService`

---

## Summary

Matchmaking for 2v2 PVP and 2v2 PVE modes fails to trigger the `startArena` engine call correctly. While the queueing logic appears to work, the "sparking" phase fails to transition to a "matched" state in automated tests, resulting in `startArena` never being called.

---

## Technical Description

### Background
The `MatchMakingController` is responsible for aggregating players into queues based on `game_mode`. When the required player count is met (1 for 1v1_PVE, 2 for 2v2_PVE, 4 for 2v2_PVP), it should create a `GameMatch` and call the Upsilon API to start the arena.

### The Problem Scenario
In `ExtraMatchmakingTest`, when 4 players join a `2v2_PVP` queue:
1. P1, P2, P3 join and are correctly queued.
2. P4 joins, meeting the `neededOpponents` count (3 others + self).
3. The logic enters the matched block but either fails to prepare the DTOs correctly or fails the engine call validation.

```php
// Fails to reach this or fails within
if ($opponents->count() >= $neededOpponents) {
    $matchedPlayers = $opponents->take($neededOpponents)->push($queueEntry);
    return $this->sparkPVEMatch($request, $matchedPlayers, $mode);
}
```

### Where This Pattern Exists Today
- `app/Http/Controllers/API/MatchMakingController.php:77-84`
- `tests/Feature/API/ExtraMatchmakingTest.php`

---

## Risk Assessment

| Factor              | Value                                         |
| ------------------- | --------------------------------------------- |
| Likelihood          | High                                          |
| Impact if triggered | High                                          |
| Detectability       | High — manifests as 100% test failure for 2v2 |
| Current mitigant    | 1v1 PVP and 1v1 PVE are working correctly     |

---

## Recommended Fix

**Short term:** Fix the DTO preparation logic in `sparkPVEMatch` to ensure it correctly handles both Eloquent Models and Arrays in the `$players` collection.

**Medium term:** Refactor matchmaking into a dedicated Service class to separate queue management from controller logic.

**Long term:** Implement a centralized matchmaking worker/orchestrator to handle multi-party matches.

---

## References

- [MatchMakingController.php](file:///workspace/battleui/app/Http/Controllers/API/MatchMakingController.php)
- [ExtraMatchmakingTest.php](file:///workspace/battleui/tests/Feature/API/ExtraMatchmakingTest.php)
- [mech_matchmaking.atom.md](file:///workspace/docs/mech_matchmaking.atom.md)
