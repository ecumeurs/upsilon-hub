# Issue: Recursive Mutex Lock Deadlock in CLI Listener

**ID:** `Ref_20260505_cli_listener_recursive_lock_deadlock`
**Ref:** `ISS-094`
**Date:** 2026-05-05
**Severity:** Critical
**Status:** Open
**Component:** `upsiloncli/internal/ws/listener.go`
**Affects:** `upsiloncli` E2E test execution, WebSocket event processing

---

## Summary

The `notifyWaiters` function in `listener.go` contains a double-lock on `l.waitMu`. Since `sync.Mutex` in Go is not recursive, this causes a permanent deadlock of the WebSocket `listenLoop` goroutine as soon as any message is received from the server. This prevents any further WebSocket event processing, causing E2E tests that depend on these events (like `match.found` or `turn.started`) to hang indefinitely.

---

## Technical Description

### Background

The `Listener` struct uses `waitMu` to synchronize access to `waiters` and `buffer` maps. `notifyWaiters` is called by the `listenLoop` for every incoming WebSocket message to dispatch the event to any pending `WaitForData` calls or buffer it.

### The Problem Scenario

1.  A message is received in `listenLoop`.
2.  `listenLoop` calls `l.notifyWaiters(envelope.Event, envelope.Data)`.
3.  `notifyWaiters` calls `l.waitMu.Lock()` at line 378.
4.  `notifyWaiters` performs some unmarshaling.
5.  `notifyWaiters` calls `l.waitMu.Lock()` **again** at line 399.
6.  The goroutine hangs forever at line 399, holding the lock from line 378.

```go
377: func (l *Listener) notifyWaiters(eventName string, data json.RawMessage) {
378: 	l.waitMu.Lock()
379: 	defer l.waitMu.Unlock()
...
399: 	l.waitMu.Lock()
400: 	defer l.waitMu.Unlock()
...
```

### Where This Pattern Exists Today

- `upsiloncli/internal/ws/listener.go` lines 377-424.

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | High |
| Impact if triggered | High |
| Detectability | High — manifests as indefinite hang in E2E tests |
| Current mitigant | None |

---

## Recommended Fix

**Short term:** Remove the redundant second lock in `notifyWaiters`.

**Medium term:** Run a project-wide audit for similar redundant locks.

**Long term:** Consider using a more robust event dispatching pattern or ensuring that mutexes are only locked once per function scope.

---

## References

- [listener.go](file:///workspace/upsiloncli/internal/ws/listener.go)
- [e2e_combat_turn_management.js](file:///workspace/upsiloncli/tests/scenarios/e2e_combat_turn_management.js)
