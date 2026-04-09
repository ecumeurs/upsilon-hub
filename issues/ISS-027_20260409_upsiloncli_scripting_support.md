# Issue: UpsilonCLI Scripting & Automated Scenario Support

**ID:** `20260409_upsiloncli_scripting_support`
**Ref:** `ISS-027`
**Date:** 2026-04-09
**Severity:** Medium
**Status:** Open
**Component:** `upsiloncli/internal/script`
**Affects:** Developers, QA, CI/CD Pipelines

---

## Summary

Currently, UpsilonCLI is primarily an interactive explorer. To facilitate complex integration testing and multi-user simulations (e.g., verifying a 2v2 PVE match end-to-end), the tool needs a non-interactive "Scripting Mode." This would allow feeding the CLI a definition file (JSON or YAML) containing a sequence of commands, conditional logic based on response data, and event-driven triggers.

---

## Technical Description

### Background
UpsilonCLI already manages a `Session` context and an `Endpoint` registry. However, moving from interactive prompts to scripted execution requires a dedicated interpreter.

### The Problem Scenario
A developer wants to test the full lifecycle of a 2v2 match. Currently, they must open 4 terminals and manually type commands. A scripting engine would allow:
1. **Command Sequences**: `auth_login` -> `matchmaking_join` -> `game_state`.
2. **Dynamic Variable Capture**: Capture `match_id` from a `MatchFound` WebSocket event and use it in subsequent `game_action` calls.
3. **Event Listeners**: "When `BoardUpdated` is received AND `current_entity_id == my_id`, execute `MOVE` then `ATTACK`."
4. **Assertions**: "Verify that `HP` decreased after `ATTACK`."

### Where This Pattern Exists Today
- [upsiloncli/internal/cli/repl.go](file:///workspace/upsiloncli/internal/cli/repl.go): The current interactive loop.
- [upsiloncli/internal/session/session.go](file:///workspace/upsiloncli/internal/session/session.go): The context store that needs to be expanded for script variables.
- [ISS-026](file:///workspace/issues/ISS-026_20260409_api_journey_tester_cli.md): The original CLI requirement which mentioned an `--auto` mode, of which this is the logical evolution.

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | High (Required for robust CI testing) |
| Impact if triggered | Low (Mainly a feature gap, not a bug) |
| Detectability | High |
| Current mitigant | Manual interactive testing or brittle bash scripts |

---

## Recommended Fix

**Short term:** Add a `--script <file.yaml>` flag to `upsiloncli`. Implement a basic YAML parser that iterates through a list of `actions`. Each action can be a `call` or a `wait_for_event`.

**Medium term:** Implement a simple DSL or use a Go-embedded scripting language (like `otto` for JS or a Lua VM) to provide full control over logic and variables.

**Long term:** Support concurrent execution of multiple scripts (multi-user simulation) from a single CLI coordinator.

---

## References

- [upsiloncli/README.md](file:///workspace/upsiloncli/README.md)
- [usecase_api_flow_game_turn.atom.md](file:///workspace/docs/usecase_api_flow_game_turn.atom.md)
- [usecase_api_flow_matchmaking.atom.md](file:///workspace/docs/usecase_api_flow_matchmaking.atom.md)
