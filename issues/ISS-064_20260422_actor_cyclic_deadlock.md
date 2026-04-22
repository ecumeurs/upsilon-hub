# Issue: Cyclic Deadlock Risk in Actor Request-Reply

**ID:** `20260422_actor_cyclic_deadlock`
**Ref:** `ISS-064`
**Date:** 2026-04-22
**Severity:** High
**Status:** Open
**Component:** `upsilontools/tools/actor`
**Affects:** `upsilontools/tools/actor`, `upsilonbattle`

---

## Summary

The current Actor implementation is prone to deadlocks when two actors perform cyclic request-reply calls. This is caused by unbuffered `CallbackChan` and synchronous `ctx.Reply` calls inside message handlers.

---

## Technical Description

### Background
Actors process messages and replies sequentially in a single dispatch loop. When an actor calls another actor via `SendActor`, it provides a `CallbackChan` for the reply.

### The Problem Scenario
1. **Actor A** is in a handler for a message.
2. **Actor B** is in a handler for a message.
3. **Actor A** calls `ctx.Reply` to send a response to **Actor B**. This blocks on `b.CallbackChan <- msg`.
4. **Actor B** calls `ctx.Reply` to send a response to **Actor A**. This blocks on `a.CallbackChan <- msg`.

Because both `CallbackChan` are unbuffered and both actors are currently busy in their handlers, neither actor can read from their channel, resulting in a permanent deadlock.

This manifests as flaking in `TestActorDeadlock_CyclicCall` because it depends on the `select` order when multiple messages are ready in the actor loop.

### Where This Pattern Exists Today
- `upsilontools/tools/actor/actor.go:70` (`ctx.Reply` sends to `ReplyChan`)
- `upsilontools/tools/actor/actor.go:190` (`CallbackChan` initialization)

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | Medium (triggered by specific cyclic patterns) |
| Impact if triggered | High (Actors hang, system becomes unresponsive) |
| Detectability | High (Test timeouts, stack dumps) |
| Current mitigant | None |

---

## Recommended Fix

**Short term:** Buffer `CallbackChan` in `Actor` struct to allow non-blocking replies for most common scenarios.
**Medium term:** Buffer `inputChan` in `MessageQueue` to ensure `SendActor` and `NotifyActor` are truly asynchronous.
**Long term:** Redesign the dispatcher to handle replies even when a handler is executing, or enforce a strict non-cyclic call graph.

---

## References

- [actor.go](file:///home/bastien/work/upsilon/upsilon-hub/upsilontools/tools/actor/actor.go)
- [actor_deadlock_repro_test.go](file:///home/bastien/work/upsilon/upsilon-hub/upsilontools/tools/actor/actor_deadlock_repro_test.go)
