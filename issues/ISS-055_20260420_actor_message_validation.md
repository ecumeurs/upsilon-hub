# Issue: Actor Message Type Validation

**ID:** `20260420_actor_message_validation`
**Ref:** `ISS-055`
**Date:** 2026-04-20
**Severity:** Low
**Status:** Open
**Component:** `upsilontools/tools/actor`
**Affects:** `upsilonapi`, `upsiloncli`

---

## Summary

The `Actor` implementation should validate if the target message is of the correct type (Call vs. Notification) at the beginning of `SendActor` and `NotifyActor`.

---

## Technical Description

### Background
`SendActor` is intended for Request-Response interactions (Calls), while `NotifyActor` is for fire-and-forget (Notifications).

### The Problem Scenario
Calling `SendActor` with a notification-type message or vice-versa can lead to protocol violations or hanging channels if not handled correctly.

### Where This Pattern Exists Today
`upsilontools/tools/actor/actor.go`

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | Low |
| Impact if triggered | Medium |
| Detectability | High |
| Current mitigant | Some panics exist, but aren't strictly enforced based on message metadata. |

---

## Recommended Fix

**Short term:** Add checks at the entry of `SendActor` and `NotifyActor`.  
**Medium term:** Update `message.Message` to include an explicit `Kind` field (Call/Notification).

---

## References

- [ui_investigation.md](file:///home/bastien/work/upsilon/projbackend/ui_investigation.md)
