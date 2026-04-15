# Issue: Aggressive AI Controller Deadlocks on Match Termination

**ID:** `20260412_ai_deadlock_termination`
**Ref:** `ISS-032`
**Date:** 2026-04-12
**Severity:** Medium
**Status:** Resolved
**Component:** `upsilonbattle/battlearena/controller/controllers`
**Affects:** `upsilonapi/bridge`

---

## Summary

The `AggressiveController` (AI) implementation uses an unbuffered channel to communicate matching termination. When the `BattleEnd` notification is received, the controller attempts to send `true` to `BattleFinished`. However, if the starting component (typically the bridge or a test runner) is not actively listening to this channel, the AI actor blocks indefinitely, causing a resource leak and potential hangs during match teardown.

---

## Technical Description

### Background
Actors in the Upsilon system should remain responsive. The AI controller is designed to signal its completion via the `BattleFinished` channel.

### The Problem Scenario
1. A battle ends. `ruler.go` notifies all controllers via `rulermethods.BattleEnd`.
2. `AggressiveController` receives the message.
3. It executes `ctl.BattleFinished <- true`.
4. Since the channel is unbuffered (`make(chan bool)`) and often not drained, the `AggressiveController` actor blocks at this line.
5. All subsequent messages to this controller are ignored, and the goroutine remains alive.

### Where This Pattern Exists Today
- [/workspace/upsilonbattle/battlearena/controller/controllers/aggressive.go:29](/workspace/upsilonbattle/battlearena/controller/controllers/aggressive.go#L29) - channel definition.
- [/workspace/upsilonbattle/battlearena/controller/controllers/aggressive.go:226](/workspace/upsilonbattle/battlearena/controller/controllers/aggressive.go#L226) - blocking send.

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | High (every AI match) |
| Impact if triggered | Medium — Resource leak / Hang during cleanup |
| Detectability | Medium — Requires monitoring goroutine count or timeouts |
| Current mitigant | None |

---

## Recommended Fix

**Short term:** Change the channel to be buffered: `make(chan bool, 1)`.
**Medium term:** Implement a proper lifecycle management for the AI actors that doesn't rely on blocking channels.

---

## References

- [aggressive.go](/workspace/upsilonbattle/battlearena/controller/controllers/aggressive.go)
- [bridge.go](/workspace/upsilonapi/bridge/bridge.go)
