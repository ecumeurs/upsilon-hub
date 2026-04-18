# Issue: Turn Start Webhook Missing When AI Goes First

**ID:** `20260418_turn_start_webhook_unicast`
**Ref:** `ISS-048`
**Date:** 2026-04-18
**Severity:** High
**Status:** Open
**Component:** `upsilonbattle/battlearena/ruler`
**Affects:** `upsilonapi/bridge/http_controller.go`

---

## Summary

The `turn.started` webhook event is intermittently missing from CI test results and production matches. This occurs because the `Ruler` sends the `ControllerNextTurn` notification as a unicast message to the active controller only. If a match involves an AI player and it is chosen to go first, any `HTTPController` (acting as a webhook proxy) will not receive the notification and thus fail to emit the `turn.started` event.

---

## Technical Description

### Background
The `upsilonapi` uses `HTTPController` to bridge Ruler notifications to webhook endpoints. One of the critical events is `turn.started`, which is triggered by the `ControllerNextTurn` message from the `Ruler`.

### The Problem Scenario
1. A match starts with two players: Player 1 (Human, using `HTTPController`) and Player 2 (AI, using `AggressiveController`).
2. The `Ruler` decides Player 2 goes first.
3. The `Ruler` sends `ControllerNextTurn` ONLY to Player 2's controller.
4. Player 1's `HTTPController` never receives the message.
5. `HTTPController` never calls `forwardToWebhook` for `turn.started`.
6. The frontend or test suite fails to receive the expected event.

### Where This Pattern Exists Today
- [ruler.go](file:///home/bastien/work/upsilon/projbackend/upsilonbattle/battlearena/ruler/ruler.go#L282) and [ruler.go](file:///home/bastien/work/upsilon/projbackend/upsilonbattle/battlearena/ruler/ruler.go#L558) use unicast `NotifyActor`.
- [aggressive.go](file:///home/bastien/work/upsilon/projbackend/upsilonbattle/battlearena/controller/controllers/aggressive.go#L98) lacks a guard to verify if the entity turn belongs to it, which would be necessary if we switch to broadcast.

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | High (Intermittent, depends on turn order) |
| Impact if triggered | High (Frontend desync, broken test assertions) |
| Detectability | Medium (Visible in CI failed tests) |
| Current mitigant | None |

---

## Recommended Fix

**Short term:** Broadcast `ControllerNextTurn` to all registered controllers in the match.  
**Medium term:** Add guards to all controller implementations of `ControllerNextTurn` to ignore notifications for entities they do not control.  
**Long term:** Standardize on a broadcast `TurnStarted` event that is separate from the `ControllerNextTurn` action-request.

---

## References

- [ruler.go](file:///home/bastien/work/upsilon/projbackend/upsilonbattle/battlearena/ruler/ruler.go)
- [http_controller.go](file:///home/bastien/work/upsilon/projbackend/upsilonapi/bridge/http_controller.go)
- [aggressive.go](file:///home/bastien/work/upsilon/projbackend/upsilonbattle/battlearena/controller/controllers/aggressive.go)
