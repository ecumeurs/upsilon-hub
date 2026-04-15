---
id: mech_actor_dispatch_loop
human_name: "Actor Dispatch Loop"
type: MECHANIC
layer: IMPLEMENTATION
version: 1.0
status: DRAFT
priority: 5
tags: [actor, dispatch, loop, logic]
parents:
  - [[mech_actor_pattern]]
dependents: []
---

# Actor Dispatch Loop

## INTENT
Centrally orchestrate the processing of incoming messages and callbacks, ensuring that the Actor remains responsive to both internal and external stimuli.

## THE RULE / LOGIC
- **Dual-Channel Selection**: The dispatcher must concurrently listen to:
  1. The `MessageQueue` executor channel (incoming requests).
  2. The `CallbackChan` (replies from other actors).
- **Hierarchical Handler Lookup**: For any incoming message, the dispatcher must search in order:
  1. **Typed Handlers**: Modern `callHandlers` or `notificationHandlers` maps.
  2. **Internal Signals**: Special handling for `ActorStarted`, `ActorStop`.
  3. **Legacy Handlers**: The `methods` map for backward compatibility.
- **Protocol Enforcement**:
  - `Call` messages must be processed via `CallContext` and checked for replies.
  - `Notification` messages must be acknowledged to the queue after handler execution.

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag**: `@spec-link [[mech_actor_dispatch_loop]]`
- **Key Methods**:
  - `Actor.Start()`: Spawns the main `select` loop goroutine.
  - `Actor.processMessage(msg)`: Entry point for logic dispatching.
  - `Actor.processReply(msg)`: Entry point for handling callback responses.

## EXPECTATION (For Testing)
- **Sequentiality**: Multiple messages on different channels must be processed one at a time via the shared logic in `processMessage`.
- **Transparency**: Unhandled messages should be logged or cause a panic based on the `CrashOnUnhandled` flag.
