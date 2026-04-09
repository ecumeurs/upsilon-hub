---
id: api_websocket_user_notifications
human_name: "WebSocket User Notifications (Private)"
type: API
layer: ARCHITECTURE
version: 1.0
status: STABLE
priority: 3
tags: [websocket, matchmaking, notifications]
parents:
  - [[api_websocket]]
dependents: []
---

# WebSocket User Notifications (Private)

## INTENT
To provide authenticated, user-specific notifications such as matchmaking results and account-level alerts.

## THE RULE / LOGIC
1. **Channel Name**: `private-user.{user_id}`
   - `{user_id}` can be the User UUID or the account nickname (if authorized).
2. **Authorization**: Only the owner of the user account can subscribe.
3. **Core Events**:
   - `match.found`: Triggered when a match is found.
     - **Payload**:
       - `match_id`: `string (UUID)`
       - `user_id`: `string (UUID)`
       - `data`: `array` (optional match metadata)

## TECHNICAL INTERFACE (The Bridge)
- **Channel Pattern:** `private-user.*`
- **Code Tag:** `@spec-link [[api_websocket_user_notifications]]`
- **Laravel Event:** `App\Events\MatchFound`

## EXPECTATION (For Testing)
- User logs in -> Subscribes to `private-user.{id}` -> Signature valid.
- Matchmaking pairs player -> Event `match.found` received by client.
