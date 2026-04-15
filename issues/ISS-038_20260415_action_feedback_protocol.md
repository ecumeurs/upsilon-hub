# Issue: Lack of Standardized Action Feedback in Protocol

**ID:** `20260415_action_feedback_protocol`
**Ref:** `ISS-038`
**Date:** 2026-04-15
**Severity:** Medium
**Status:** Open
**Component:** `upsilonapi`, `upsilonbattle`, `communication.md`
**Affects:** `battleui` (Frontend), `upsiloncli`

---

## Summary

Currently, when a player (AI or human) takes an action, the system broadcasts the new game state, but does not provide explicit feedback about the action itself (e.g., damage dealt, path taken, who attacked who). Clients (Frontend/CLI) must diff the game state to understand what happened, which leads to poor UX (teleporting figures, lack of visual feedback).

---

## Technical Description

### Background
The Upsilon Engine (Go) handles tactical actions and broadcasts state updates via webhooks. The `board.updated` event contains the full `BoardState`.

### The Problem Scenario
1. Player A attacks Player B.
2. Engine computes damage and updates Player B's HP.
3. Engine broadcasts `board.updated` with the new HP.
4. Frontend receives the state, see Player B has less HP, but doesn't know *why* or *by whom* it was hit without complex diffing.
5. Frontend cannot easily trigger a "Hit" animation or show floating damage numbers at the right moment.

### Where This Pattern Exists Today
- `upsilonapi/api/output.go`: `BoardState` struct.
- `upsilonapi/bridge/http_controller.go`: Webhook dispatch logic.
- `upsilonbattle/battlearena/ruler/rulermethods/rulermethods.go`: Protocol messages.

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | High |
| Impact if triggered | Medium |
| Detectability | High |
| Current mitigant | None (manual state diffing required) |

---

## Recommended Fix

**Short term:** Add an `ActionFeedback` field to the `BoardState` and populate it during action processing in the bridge.  
**Medium term:** Enrich `upsilonbattle` broadcast messages (`ControllerAttacked`, `ControllerMoved`) with specific result data (damage, path).  
**Long term:** Implement a dedicated event log/history service for replayability and audit.

---

## References

- [communication.md](file:///workspace/communication.md)
- [upsilonapi/README.md](file:///workspace/upsilonapi/README.md)
- [rulermethods.go](file:///workspace/upsilonbattle/battlearena/ruler/rulermethods/rulermethods.go)
