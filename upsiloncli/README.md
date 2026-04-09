# UpsilonCLI — API Journey Explorer & Tester

**UpsilonCLI** is an interactive command-line tool for exploring and testing the Upsilon Battle API ecosystem. It provides transparent access to every Laravel Gateway endpoint, real-time WebSocket monitoring, and tactical board visualization — all from the terminal.

**Tracking Issue:** [ISS-026](../issues/ISS-026_20260409_api_journey_tester_cli.md)

## Quick Start

```bash
cd /workspace/upsiloncli
go run ./cmd/upsiloncli
```

The CLI defaults to `http://localhost:8000` as the Laravel API base URL. Override with:
```bash
go run ./cmd/upsiloncli --base-url http://custom-host:8000
```

## Commands

| Command | Description |
|---|---|
| `routes` | List all available API endpoints with their `route_name` identifiers. |
| `call <route_name>` | Execute an endpoint interactively. Prompts for each input parameter with smart defaults from session context. |
| `jwt` | Display the current JWT token. |
| `jwt <token>` | Manually override the active JWT (for testing invalid/expired tokens). |
| `session` | Display current session context (user_id, match_id, characters, etc.). |
| `redraw` | Re-render the last known tactical board state. |
| `help` | Show available commands. |
| `exit` | Quit the CLI. |

### Auto Mode

Run a full end-to-end journey (Register → Login → Profile → Matchmaking → Combat → Cleanup) non-interactively:
```bash
go run ./cmd/upsiloncli --auto
```

## Architecture

```
cmd/upsiloncli/       Entry point (main.go)
internal/
  cli/                REPL loop, command dispatcher
  session/            JWT management, context store (match_id, user_id, etc.)
  api/                HTTP client, curl logger, response parser
  endpoint/           Endpoint registry and individual route implementations
  display/            Terminal output formatting, board renderer
```

## Transparency

Every API call made by the CLI displays:
1. The **full curl command** equivalent (copy-paste ready).
2. The **raw JSON response** (pretty-printed).
3. A **human-readable summary** of the result.

## JWT Lifecycle

- **Auto-capture**: Tokens from `login` / `register` responses are cached automatically.
- **Renewal**: If a response contains `meta.token`, the CLI transparently rotates its JWT per `[[mech_sanctum_token_renewal]]`.
- **Clearance**: Tokens are wiped on `logout` or `auth_delete` actions.
- **Override**: Use `jwt <token>` to inject an arbitrary token for testing.

## Session Context & Smart Defaults

The CLI tracks named values from API responses (e.g., `user_id`, `match_id`, `character_id`). When calling an endpoint that requires one of these parameters, the CLI pre-fills the default from context. The user can accept or override.

## Dependencies

- Go 1.25+
- Laravel API running on `localhost:8000`
- Upsilon Engine running on `localhost:8081`

## Related Documentation

- [Communication Reference](../communication.md)
- [API Gateway ATD](../docs/api_laravel_gateway.atom.md)
- [Matchmaking Flow](../docs/usecase_api_flow_matchmaking.atom.md)
- [Game Turn Flow](../docs/usecase_api_flow_game_turn.atom.md)
- [Token Renewal](../docs/mech_sanctum_token_renewal.atom.md)
