# Issue: Cyclic Deadlock Risk in Actor Request-Reply

**ID:** `20260422_actor_cyclic_deadlock`
**Ref:** `ISS-064`
**Date:** 2026-04-22
**Severity:** High
**Status:** Resolved
**Component:** `upsilontools/tools/actor`
**Affects:** `upsilontools/tools/actor`, `upsilonbattle`

---

## Summary

The Actor implementation was prone to deadlocks when two actors perform cyclic request-reply calls. This was caused by the "side-channel" nature of `CallbackChan`, which bypassed the `MessageQueue` and blocked the sender if the recipient was busy.

---

## Technical Description

### Background
Actors process messages and replies sequentially in a single dispatch loop. When an actor calls another actor via `SendActor`, it provides a `CallbackChan` for the reply.

### The Problem Scenario
1. **Actor A** is in a handler for a message.
2. **Actor B** is in a handler for a message.
3. **Actor A** calls `ctx.Reply` to send a response to **Actor B**. This blocks on `b.CallbackChan <- msg`.
4. **Actor B** calls `ctx.Reply` to send a response to **Actor A**. This blocks on `a.CallbackChan <- msg`.

Because both `CallbackChan` were unbuffered and both actors were busy in their handlers, neither actor could read from their channel, resulting in a permanent deadlock.

### The Fix: Unified Queue Dispatch
The Actor dispatch loop was refactored to listen exclusively to the `MessageQueue`. A background redirector now pipes all `CallbackChan` stimuli into the `MessageQueue`. 

Benefits:
- **Non-blocking Replies**: The `MessageQueue` accepts the reply immediately into its internal slice, unblocking the sender.
- **Sequentiality**: Replies are processed one-by-one, interleaved with other messages in arrival order.
- **Traceability**: All stimuli now honor the `mq.currentMessage` execution lock.

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | Medium (triggered by specific cyclic patterns) |
| Impact if triggered | High (Actors hang, system becomes unresponsive) |
| Detectability | High (Test timeouts, stack dumps) |
| Current mitigant | None |

---

## Recommended Fix (Implemented)

**Unified Queue Dispatch:** All replies are routed through the `MessageQueue`. This eliminates the side-channel race condition and ensures the Actor is never blocked from receiving a response while executing a handler.

---

## References

- [actor.go](file:///home/bastien/work/upsilon/upsilon-hub/upsilontools/tools/actor/actor.go)
- [actor_deadlock_repro_test.go](file:///home/bastien/work/upsilon/upsilon-hub/upsilontools/tools/actor/actor_deadlock_repro_test.go)
