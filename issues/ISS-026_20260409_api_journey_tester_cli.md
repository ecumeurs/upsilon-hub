# Issue: Upsilon API Journey Explorer & Tester CLI

**ID:** `20260409_api_journey_tester_cli`
**Ref:** `ISS-026`
**Date:** 2026-04-09
**Severity:** Medium
**Status:** Open
**Component:** `tools/journey-tester`
**Affects:** Developers, QA, API Users, Integrators

---

## Summary

The system requires a versatile CLI tool to facilitate development and testing of the API-first ecosystem. This tool must go beyond basic sequence testing, providing an "Interactive Explorer" mode that facilitates access to all API methods and real-time monitoring of WebSocket events.

---

## Technical Description

### Background

The Upsilon ecosystem relies on a complex interplay between Laravel (Gateway/WebSockets) and Go (Battle Engine). Manually tracking these interactions across REST and WebSockets is high-friction for developers.

### The Problem Scenario

Developers need a tool that:
1. **Facilitates API Access**: 
    - `routes`: List all available endpoints with unique `route_name` identifiers.
    - `call <route_name>`: Interactive command that prompts for each required input parameter, supporting smart defaults from session context (e.g., auto-filling the last received `match_id` or `character_id`).
2. **WebSocket Awareness**: Listen to and display push notifications from Laravel Reverb (e.g., `MatchFound`, `BoardUpdated`) in real-time.
3. **Advanced Session Management**:
    - **JWT Lifecycle**: Automatic capture of tokens during Login/Register; cache clearance on Logout/Account termination.
    - **Token Renewal**: Per [[mech_sanctum_token_renewal]], the CLI must monitor `meta.token` in response envelopes and transparently rotate the active JWT.
    - **Manual Override**: `jwt <new_token>` command to allow intentional session hijacking or testing with invalid/expired keys.
4. **Tactical Visualization**:
    - **Dynamic Rendering**: Upon receiving `battle.start` or `board.updated` events, the CLI must draw an ASCII/text-based representation of the tactical board.
    - **Telemetry**: Display entity positions, HP, and active status effects alongside the board.
    - **Manual Refresh**: `redraw` command to force a re-render of the current tactical state.
5. **Transparency**: Always display the equivalent `curl` command and pretty-printed response for every action.
6. **Autopilot Mode**: Maintain the ability to run a full "Journey" (Register to Delete) via a single `--auto` flag.

```
[CLI] session: {user_id: "...", match_id: "match-123"}
[CLI] > call game_action
[PROMPT] entity_id [default: char-456]: 
[PROMPT] type [MOVE|ATTACK|PASS]: ATTACK
[PROMPT] target_coords (x,y): 5,2
[CURL] curl -X POST ... /api/v1/game/match-123/action ...
[REPLY] { "success": true, "meta": {"token": "renewed-jwt-..."} }
[SYSTEM] Token renewed and cached.
[WS] BoardUpdated event received.
[VISUAL]
  . . . E .
  . P . . .
  . . . . .
Entities: [P] Player (HP:10) [E] Enemy (HP:7)
```

### Where This Pattern Exists Today

The use cases and documentation to follow:
- [usecase_api_flow_matchmaking.atom.md](file:///workspace/docs/usecase_api_flow_matchmaking.atom.md)
- [usecase_api_flow_game_turn.atom.md](file:///workspace/docs/usecase_api_flow_game_turn.atom.md)
- [mech_sanctum_token_renewal.atom.md](file:///workspace/docs/mech_sanctum_token_renewal.atom.md)
- [communication.md](file:///workspace/communication.md)

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | High |
| Impact if triggered | Medium (Development friction, regression risk) |
| Detectability | High |
| Current mitigant | Manual browser/Postman/Tinker testing |

---

## Recommended Fix

**Short term:** Implement a Go CLI in `tools/journey-tester` using `cobra` or `urfave/cli` for command management, `github.com/gorilla/websocket` for Reverb integration, and a terminal-based UI library (e.g., `github.com/pterm/pterm`) for board rendering.

**Medium term:** Add auto-discovery of endpoints via `/api/v1/help` to populate the `route_name` registry dynamically.

**Long term:** Extend to support Multi-User simulation (concurrent CLI instances).

---

## References

- [usecase_api_flow_matchmaking.atom.md](file:///workspace/docs/usecase_api_flow_matchmaking.atom.md)
- [usecase_api_flow_game_turn.atom.md](file:///workspace/docs/usecase_api_flow_game_turn.atom.md)
- [api_laravel_gateway.atom.md](file:///workspace/docs/api_laravel_gateway.atom.md)
- [communication.md](file:///workspace/communication.md)
