# Issue: forwardToWebhook Missing Default Handler

**ID:** `20260420_webhook_missing_default`
**Ref:** `ISS-057`
**Date:** 2026-04-20
**Severity:** Low
**Status:** In Progress
**Component:** `upsilonapi/bridge`
**Affects:** `upsilonapi`

---

## Summary

The `forwardToWebhook` function in `HTTPController` uses a switch statement on message types but lacks a `default` case, which leads to silent failures or unhandled events when new message types are added.

---

## Technical Description

### Background
`forwardToWebhook` is responsible for translating Go actor messages into API event payloads for the frontend/webhooks.

### The Problem Scenario
When an event like `ControllerForfeit` or `BattleEnd` occurs, if it's not explicitly handled in the switch, no action feedback is produced, potentially leaving the UI in an inconsistent state or missing important logs.

### Where This Pattern Exists Today
`upsilonapi/bridge/http_controller.go:60`

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | Medium |
| Impact if triggered | Low |
| Detectability | Low |
| Current mitigant | None |

---

## Recommended Fix

**Short term:** Add a `default` case that logs the unexpected message type.  
**Medium term:** Implement a more robust event-to-API mapping system.

---

## References

- [ui_investigation.md](file:///home/bastien/work/upsilon/projbackend/ui_investigation.md)
