# Issue: Actor Stop Panic: Close of Closed Channel

**ID:** `20260511_actor_stop_panic_race`
**Ref:** `ISS-097`
**Date:** 2026-05-11
**Severity:** Critical
**Status:** Resolved
**Component:** `upsilontools/tools/actor`
**Affects:** `upsilonapi`, `upsilonbattle`

---

## Summary

The `Actor.Stop()` method was not idempotent, causing a `panic: close of closed channel` when an actor was signaled to stop multiple times. This occurred during arena destruction because a shared `HTTPController` was registered under multiple player IDs in the `Ruler`'s controller map, receiving an `ActorStop` message for each ID.

---

## Technical Description

### Background
The `Actor` model in Upsilon uses a `stopper` channel to signal background loops to exit. The `Stop()` method is responsible for closing this channel.

### The Problem Scenario
1.  In `upsilonapi`, multiple human players are managed by a single `HTTPController` instance to consolidate tactical events.
2.  The `ArenaBridge` registers this shared controller with the `Ruler` for each player ID.
3.  When a match ends, the `Ruler` iterates over its `Controllers` map and sends an `ActorStop` notification to each.
4.  If the map contains the same `HTTPController` multiple times (one per player), the controller receives multiple `ActorStop` messages.
5.  Each `ActorStop` message triggers a `defer a.Stop()` in the actor's dispatch loop.
6.  The first call to `Stop()` succeeds and closes the `stopper` channel.
7.  The second call to `Stop()` attempts to close the same channel again, triggering a Go runtime panic.
8.  Since this panic happens in a background goroutine, it crashes the entire engine process, leading to `ENGINE_UNREACHABLE` errors for subsequent client requests.

### Where This Pattern Exists Today
- `upsilontools/tools/actor/actor.go:283` (Original location of `close(a.stopper)`)
- `upsilonbattle/battlearena/ruler/ruler_lifecycle.go:120` (Cascading stop logic)

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | High (whenever multiple human players are in a match) |
| Impact if triggered | Critical (engine crash, service unavailability) |
| Detectability | High (explicit panic in logs) |
| Current mitigant | None prior to fix |

---

## Recommended Fix

**Short term:** Make `Actor.Stop()` idempotent using `sync.Once`.
**Medium term:** Deduplicate controllers in the `Ruler` before sending cascading stop signals.
**Long term:** Refactor `HTTPController` to only be registered once in the `Ruler`, or handle multi-player registration more robustly at the architecture level.

---

## Extra Data

The panic was observed in `engine.log` during the `e2e_progression_post_win_with_2` test suite execution.

---

## References

- [actor.go](file:///workspace/upsilontools/tools/actor/actor.go)
- [ruler_lifecycle.go](file:///workspace/upsilontools/tools/battlearena/ruler/ruler_lifecycle.go)
- [bridge_start.go](file:///workspace/upsilonapi/bridge/bridge_start.go)
