---
id: mech_actor_pattern
human_name: "Actor Design Pattern"
type: MECHANIC
layer: IMPLEMENTATION
version: 1.0
status: DRAFT
priority: 5
tags: [actor, encapsulation, messaging, concurrency]
parents:
  - [[mech_message_queue]]
dependents: []
---

# Actor Design Pattern

## INTENT
Encapsulate internal state and behavioral logic behind a message-driven communication interface to prevent race conditions and simplify concurrency.

## THE RULE / LOGIC
- **Message-Driven**: All interactions with an Actor must occur via messages sent to its `inputChan`.
- **Single Threaded Processing**: An Actor processes exactly one message at a time using its internal `MessageQueue`.
- **Dispatcher Pattern**: The Actor must map incoming message types to specific handler functions.
- **Protocol Compliance**:
  - `Call`: Must produce exactly one Reply (payload or ACK) to the caller AND an ACK to the queue.
  - `Notification`: Must produce exactly one ACK to the queue.
- **Error Isolation**: Failures in one Actor should be logged and acknowledged to the queue to prevent blocking, but should not crash the Actor unless explicitly requested (e.g., protocol violation).

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag**: `@spec-link [[mech_actor_pattern]]`
- **Primary Interface**:
  - `NotifyActor(msg)`: Fire-and-forget notification.
  - `SendActor(msg, callbackChan)`: Request-reply call.
  - `AddCallHandler(type, handler, validator)`: Register a synchronous handler.
  - `AddNotificationHandler(type, handler, validator)`: Register an asynchronous handler.

## EXPECTATION (For Testing)
- **Atomicity**: Actor state is only modified during message processing.
- **Liveness**: Every message dispatched to a handler MUST eventually result in an ACK to the queue, even if validation fails or a panic occurs in the handler (if recoverable).
- **Isolation**: Concurrent calls to `NotifyActor` or `SendActor` from different goroutines are safe.
