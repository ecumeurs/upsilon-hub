# Issue: ATD for `error_key` taxonomy and possible envelope promotion

**ID:** `20260425_error_key_atd_and_envelope`
**Ref:** `ISS-080`
**Date:** 2026-04-25
**Severity:** Medium
**Status:** Open
**Component:** `upsilonapi/api/output.go`, `upsilonapi/bridge/bridge.go`, `upsilonbattle/battlearena/ruler/rules/*`, `battleui/app/Traits/ApiResponder.php`, `upsiloncli/internal/script/bridge.go`
**Affects:** Every test, frontend, or CLI consumer that wants to discriminate on a typed failure.

---

## Summary

`error_key` is now plumbed end-to-end (engine ruler → upsilonapi handler → Laravel proxy → CLI/JS) via `meta.error_key` in the standard envelope. There is no atom that defines the format, the namespacing rules, or the canonical list of keys, and there is no rule preventing two contributors from inventing competing keys for the same condition. Per maintainer guidance, `error_key` should:

1. Be governed by an ATD atom (or atom family) that defines its grammar and lists every key the API can surface.
2. Possibly be promoted from `meta.error_key` to a top-level `error_key` field on the envelope, since it is now part of the *effective* contract for typed error handling and "fishing it out of meta" is awkward in PHP/JS clients.

This issue tracks both decisions.

---

## Technical Description

### Background

- The standard envelope is `{ request_id, message, success, data, meta }` — see `[[api_standard_envelope]]`.
- `meta` is documented as "Side information for debugging or testing (optional)" and is the slot we currently use to carry `error_key`.
- The engine emits keys like `entity.path.obstacle`, `entity.turn.missmatch`, `entity.controller.missmatch`, `rule.friendly_fire`, `entity.movement.already`, `entity.path.too.long`, `entity.path.notvalid`, `entity.path.notadjacent`, `entity.path.occupied`, `entity.movement.nocredits`, `entity.movement.credits`. The shape is loosely "domain.subdomain.kind" with three components, but the rule is unwritten.
- New keys introduced 2026-04-24 by the bridge layer (request validation, arena lookup): `request.target_coords.missing`, `request.player_id.invalid`, `request.entity_id.invalid`, `arena.notfound`. These are mine; they have no atom backing.

### The Problem Scenario

A frontend developer wants to react to `entity.path.obstacle` differently from `entity.path.notadjacent`. They try to read `error.error_key` and it works for `action`/`forfeit` but does not exist for matchmaking, auth, or any Laravel `FormRequest` rejection (those are out of scope of this issue — see ISS-081). They then write a defensive fallback. The next contributor adds a new ruler error and uses `entity.attack.range` instead of `entity.skill.range` because there is no canonical list. Drift accumulates.

### Where This Pattern Exists Today

- Engine emitters: `upsilonbattle/battlearena/ruler/rules/move.go` (lines 109-184), `attack.go` (165-170), `skill.go` (135-140).
- API surface: `upsilonapi/api/output.go:NewErrorWithKey`, `upsilonapi/bridge/bridge.go` (action + forfeit returns).
- Laravel propagation: `battleui/app/Traits/ApiResponder.php:error()`, `battleui/app/Http/Controllers/API/GameController.php`.
- CLI consumption: `upsiloncli/internal/script/bridge.go:jsCall` (lifts `meta.error_key` to `e.error_key`).

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | High — every new ruler condition or validation rule risks adding a divergent key. |
| Impact if triggered | Medium — tests pass for the wrong reason; clients write brittle string compares. |
| Detectability | Low — no compile-time check, no lint, no canonical list to diff against. |
| Current mitigant | None beyond convention. |

---

## Recommended Fix

**Short term:**
- File a new atom `[[api_error_keys]]` (REQUIREMENT, ARCHITECTURE) that:
  - Locks the grammar to `<domain>.<subdomain>.<kind>` (3 dotted segments).
  - Lists every key currently surfaced (mining the engine + bridge + Laravel sources).
  - Marks ownership: each key links back to the atom whose rule produces it.
- Add a Go test in `upsilonbattle` that walks `ReplyWithError(...)` call sites and asserts every emitted key is present in the atom's enumeration (parser regex over the .atom.md file is enough — no runtime registry needed).

**Medium term:**
- Promote `error_key` to a first-class envelope field. Update `[[api_standard_envelope]]`:
  ```
  { request_id, message, success, data, meta, error_key? }
  ```
  `error_key` only present when `success == false`. This is a communication-layer change and per `CLAUDE.md` requires the maintainer's explicit go-ahead before merging.
- Update `stdmessage.StandardMessage` (Go), `ApiResponder` (PHP), and `client.Response` (CLI) in lockstep. Keep `meta.error_key` as a transitional alias for one release.

**Long term:**
- Generate a typed enum (`ErrorKey`) in Go and a TypeScript discriminated union for the frontend, both produced from the atom file. This makes drift impossible.

---

## References

- `[[api_standard_envelope]]` — current envelope contract
- `upsilonapi/bridge/bridge.go:257-352` — handler-facing returns
- `upsilonapi/api/output.go:107-138` — `NewError` / `NewErrorWithKey`
- `upsilonbattle/battlearena/ruler/rules/move.go`, `attack.go`, `skill.go` — emitter sites
- `upsiloncli/internal/script/bridge.go:jsCall` — consumer-side
- Sister issue [ISS-081](ISS-081_20260425_cross_stack_error_handling.md) — harmonization across non-engine paths
