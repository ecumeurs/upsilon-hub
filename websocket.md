# Upsilon Battle: WebSocket Protocol Specification

> [!TIP]
> **Live Event Registry:** A machine-readable list of all WebSocket events, channel patterns, and subscription triggers is available via the [Live API Help Endpoint](http://localhost:8000/api/v1/help).

This document details the real-time communication protocol used by the Upsilon Battle Arena, based on the **Pusher v7 Protocol** implemented via **Laravel Reverb**. This implementation follows the `[[api_websocket]]` specification.

## 1. Connection Initiation

The WebSocket server listens on port `8080` (default for Reverb). To connect, you must provide your `REVERB_APP_KEY`.

### Via wscat (Diagnostic Tool)
```bash
wscat -c "ws://127.0.0.1:8080/app/qtjp54myattne9euwedu?protocol=7&client=js&version=8.4.0-rc2&flash=false"
```

---

## 2. Handshake Procedure

Upon connection, the server immediately sends a connection establishment event. You **must** capture the `socket_id` from this message; it is required for all private channel authorizations.

### Server Response
```json
{
  "event": "pusher:connection_established",
  "data": "{\"socket_id\":\"1234.5678\",\"activity_timeout\":30}"
}
```
*Note: The `data` field is often returned as a double-encoded JSON string.*

---

Channel authorization ensures that only the rightful owners can access private streams, as defined in `[[api_websocket]]`.
Authentication happens at the **channel level**, not the connection level. Sensitive data (match updates, notifications) is sent over `private-` channels.

### Auth Request (HTTP)
To subscribe to a private channel, you must obtain an authorization signature from the Laravel API.

**Endpoint:** `POST /broadcasting/auth`

**Example CURL:**
```bash
curl -X POST http://localhost:8000/broadcasting/auth \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -H "Accept: application/json" \
  -d "socket_id=1234.5678" \
  -d "channel_name=private-arena.match-uuid"
```

### API Reply
```json
{
  "auth": "qtjp54myattne9euwedu:5f9e8a7b6c5d4e3f2a1b"
}
```

---

## 4. Subscription Protocol (In wscat)

Once you have the `auth` signature from the CURL command above, **switch back to your wscat terminal** and send a `pusher:subscribe` event.

### Client Message
```json
{
  "event": "pusher:subscribe",
  "data": {
    "channel": "private-arena.123-uuid",
    "auth": "qtjp54myattne9euwedu:5f9e8a7b6c5d4e3f2a1b"
  }
}
```

### Server Acknowledgment
If successful, the server replies with:
```json
{
  "event": "pusher_internal:subscription_succeeded",
  "channel": "private-arena.123-uuid"
}
```

---

## 5. Game Events

### MatchFound
Sent on the `private-user.{ws_channel_key}` channel when the matchmaking engine pairs you with an opponent.
- **Specification:** `[[api_websocket_user_notifications]]`
- **Pseudonym:** Uses the `ws_channel_key` pseudonym retrieved from the `UserResource` to avoid revealing the raw User UUID.
- **Event Name:** `match.found`
- **Payload:** `{"match_id": "uuid"}`

### Board Updated (Tactical State)
Sent on the `private-user.{ws_channel_key}` channel whenever the tactical state changes.
- **Specification:** `[[api_websocket_arena_updates]]`
- **Surgical Masking:** Unlike global broadcasts, this event is customized for each recipient using the `BoardStateResource`. The `is_self` and `current_player_is_self` flags are pre-computed based on the recipient's identity. Additionally, characters that have been eliminated are preserved in the roster with `dead: true` and `hp: 0` to maintain state consistency.
- **Event Name:** `board.updated`
- **Payload:** Strictly follows the `[[api_standard_envelope]]` format. The tactical state is located in the `data` field of the envelope.
- **Envelope Example:** `{"request_id": "uuid", "success": true, "data": {"match_id": "uuid", ...BoardState...}}`

---

## 6. Connection Maintenance (Ping/Pong)

To prevent the Reverb server from closing the connection due to inactivity (typically after 60 seconds), the client should send a heartbeat.

### Client Heartbeat (Every 30s)
Paste this into `wscat`:
```json
{"event":"pusher:ping"}
```

### Server Response
```json
{"event":"pusher:pong"}
```

---

## 7. Complete Walkthrough Example

1. **Connect**: `wscat -c "ws://127.0.0.1:8080/app/qtjp...?"`
2. **Receive SocketID**: `{"socket_id": "888.999"}`
3. **Get Auth**: 
   ```bash
   curl -X POST http://localhost:8000/broadcasting/auth \
     -H "Authorization: Bearer <TOKEN>" \
     -d "socket_id=888.999&channel_name=private-user.my-ws-key"
   ```
4. **Subscribe**: Paste `{"event":"pusher:subscribe","data":{"channel":"private-user.my-ws-key","auth":"key:sig"}}` into wscat.
5. **Listen**: Wait for `match.found` and `board.updated` events on the same stream.

---

## 8. Traceability

This protocol is governed by the following ATD atoms:
- **Master Protocol:** `[[api_websocket]]`
- **User Notifications:** `[[api_websocket_user_notifications]]`
- **Arena Updates:** `[[api_websocket_arena_updates]]`
- **Data Structures:** `[[battleui_api_dtos]]`

@spec-link [[api_websocket]]
@spec-link [[api_websocket_user_notifications]]
@spec-link [[api_websocket_arena_updates]]
