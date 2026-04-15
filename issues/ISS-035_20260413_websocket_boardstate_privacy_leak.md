# Issue: WebSocket BoardState Privacy Leak

**ID:** `20260413_websocket_boardstate_privacy_leak`
**Ref:** `ISS-035`
**Date:** 2026-04-13
**Severity:** High
**Status:** Resolved
**Component:** `battleui/app/Events/BoardUpdated.php`
**Affects:** `battleui/resources/js/services/game.js`, `upsiloncli/internal/ws/listener.go`

---

## Summary

The `board.updated` WebSocket event currently broadcasts the same state payload to all players subscribed to a match's private channel. This results in an information leak where players can see full details of opponent characters (AI or Human) that should be hidden (Fog of War / competitive integrity).

---

## Technical Description

### Background
The `BoardUpdated` event in Laravel is broadcast on `private-arena.{match_id}`. This is a shared channel for all participants of a match.

### The Problem Scenario
When the Go Battle Engine reports a state change via webhook, Laravel triggers `BoardUpdated`. The `broadcastWith()` method (see [BoardUpdated.php](file:///workspace/battleui/app/Events/BoardUpdated.php#L41-L70)) performs basic ID masking but does not filter character attributes based on ownership.

Current filtering only unsets IDs:
```php
if (isset($payload['entities']) && is_array($payload['entities'])) {
    foreach ($payload['entities'] as &$entity) {
        unset($entity['player_id']);
    }
}
```

### Where This Pattern Exists Today
- [BoardUpdated.php](file:///workspace/battleui/app/Events/BoardUpdated.php)
- [websocket.md](file:///workspace/websocket.md) (documentation lacks privacy constraints)

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | High |
| Impact if triggered | High (Competitive advantage, unintended bot behavior) |
| Detectability | High — manifest in browser network tab or bot logs |
| Current mitigant | Basic ID masking (UUID removal) |

---

## Recommended Fix

**Short term:**
- Implement per-user broadcasting for match updates instead of a shared channel.
- Or, use Laravel's internal logic to filter `broadcastWith` based on the recipient's `socket_id` if possible (though shared channels usually send one payload).

**Medium term:**
- Refactor match updates to use `private-match.{match_id}.player.{id}` channels.
- Filter the `entities` array:
    - If `is_self` (or owned by user): return all fields.
    - If not owned: return only `hp` and basic semantic identifiers.

**Long term:**
- Investigate if the Laravel API should maintain a strict mapping to allow the Go engine to send differential updates directly, rather than the Gateway doing all the filtering.

---

## References

- [BoardUpdated.php](file:///workspace/battleui/app/Events/BoardUpdated.php)
- [websocket.md](file:///workspace/websocket.md)
- [api_websocket_arena_updates.atom.md](file:///workspace/docs/api_websocket_arena_updates.atom.md)
