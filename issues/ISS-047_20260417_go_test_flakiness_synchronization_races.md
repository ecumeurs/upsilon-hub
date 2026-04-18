# Issue: Go Unit Test Flakiness caused by Initialization Races and Timeout Sensitivity

**ID:** `20260417_go_test_flakiness_synchronization_races`
**Ref:** `ISS-047`
**Date:** 2026-04-17
**Severity:** High
**Status:** Resolved
**Component:** `upsilonbattle/battlearena/ruler`
**Affects:** `upsilonapi` (Unit tests), `github.com/ecumeurs/upsilonapi/bridge`

---

## Summary

The Go unit test suite (notably `TestShotClockCancellation` and `TestArenaStartEndpoint`) exhibits intermittent failures in CI due to initialization races and extreme sensitivity to runner latency. This is a technical manifestation of the architectural issues documented in ISS-009 and ISS-010.

---

## Technical Description

### Background
The `Ruler` actor manages turn-based logic. As per **ISS-009**, its state is frequently manipulated directly by the `ArenaBridge` during initialization.

### The Problem Scenario
1.  **Ruler-Bridge Race**: `ArenaBridge.StartArena` creates a `Ruler`, which immediately starts its message loop via `init()`. The bridge then updates `Ruler.GameState.Grid` and calls `AddEntity` directly on the `Ruler` pointer. If the `Ruler`'s message loop processes a message while the bridge is writing, a data race occurs.
2.  **ShotClock Goroutine Race**: The `ShotClock` timer uses `time.AfterFunc` which executes in a separate goroutine. It reads `r.GameState.GetTurn()` to guard against stale timeouts. This read is unsynchronized with the main actor loop's `IncTurn()` calls.
3.  **Timeout Sensitivity**: GitHub Actions runners are often slow. The strict 2-second timeout in `ExpectMessage` (ruler tests) and `waitForWebhook` (api tests) is frequently exceeded, leading to false negatives in CI.

### Where This Pattern Exists Today
- `upsilonbattle/battlearena/ruler/ruler.go:646-648` (unprotected state read in timeout goroutine)
- `upsilonapi/bridge/bridge.go:62-63`, `89` (ownership bypass during init)
- `upsilonbattle/battlearena/ruler/ruler_shotclock_test.go:71` (short 2s timeout)
- `upsilonapi/main_test.go:229` (short 5s timeout for API/WS roundtrip)

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | High |
| Impact if triggered | High (CI blockage, false negative CI results) |
| Detectability | High — manifests as `Timeout waiting for message` or `-race` detector failures |
| Current mitigant | Rerunning CI jobs (manual) |

---

## Recommended Fix

**Short term:**
- Implement a synchronized `Configure` message for `Ruler` or delay `Start()` until config is finished.
- Increase `ExpectMessage` and webhook timeouts to 10s.
- Wrap all `ShotClock` logic within the actor's thread-safe queue.

**Medium term:** Complete **ISS-009** and **ISS-010** by making `GameState` private and strictly using the message queue for all interactions.

---

## References

- [ISS-009: Ruler ownership bypass](ISS-009_20260311_ruler_ownership_bypass.md)
- [ISS-010: Ruler readiness logic](ISS-010_20260311_ruler_readiness_logic.md)
- [ruler.go](file:///home/bastien/work/upsilon/projbackend/upsilonbattle/battlearena/ruler/ruler.go)
- [bridge.go](file:///home/bastien/work/upsilon/projbackend/upsilonapi/bridge/bridge.go)
