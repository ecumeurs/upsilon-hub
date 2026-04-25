# Issue: Cross-stack error handling harmonization

**ID:** `20260425_cross_stack_error_handling`
**Ref:** `ISS-081`
**Date:** 2026-04-25
**Severity:** Medium
**Status:** Open
**Component:** `battleui/app/Http/Controllers/API/*`, `battleui/app/Http/Requests/**`, `battleui/app/Exceptions/Handler.php`, `upsilonapi/handler/*.go`, `upsiloncli/internal/api/*.go`
**Affects:** Every external client that wants typed error handling beyond the `action`/`forfeit` happy paths.

---

## Summary

`error_key` is currently propagated only on the engine action paths (`POST /game/{id}/action`, `POST /game/{id}/forfeit`). All other error paths — Laravel `FormRequest` rejections, auth (`/auth/*`), matchmaking (`/matchmaking/*`), profile/character (`/profile/*`), admin (`/admin/*`), and Laravel framework exceptions (404, 419, validation failures) — return envelopes without an `error_key`. Tests therefore have to substring-match `e.message` for every non-engine failure, which is fragile across i18n changes and Laravel version upgrades.

This issue tracks the cross-stack harmonization needed to give every typed failure an `error_key`. It pairs with [ISS-080](ISS-080_20260425_error_key_atd_and_envelope.md) (which defines the taxonomy and asks whether `error_key` should be promoted to a top-level envelope field).

---

## Technical Description

### Background

- Today, `error_key` only exists on engine rejections (move/attack/skill rules + arena lookup + a handful of bridge-level request validation cases).
- Laravel `FormRequest` validation produces a 422 with `{ errors: { field: [...] } }` shape; not envelope-conformant and not key-tagged.
- Auth controllers throw `ValidationException` or hand-roll `error("Invalid credentials", 401)`; no typed key.
- Matchmaking controllers similarly use `error()` with prose.
- The CLI scenario suite has been writing tests like:
  ```js
  upsilon.assertEquals(e.error_key, "rule.password_policy", ...)
  ```
  That assertion does not hold today and would only hold once Laravel side propagates a key.

### Where This Pattern Exists Today

Search will turn up a long list, but representative offenders:

- `battleui/app/Http/Controllers/API/AuthController.php` — login/register/update/password rejections
- `battleui/app/Http/Controllers/API/MatchMakingController.php` — `Already in queue`, `In a match`, etc.
- `battleui/app/Http/Controllers/API/ProfileController.php` — reroll restrictions, attribute cap rejections
- `battleui/app/Exceptions/Handler.php` — global rendering of `ValidationException`, `AuthenticationException`, `NotFoundHttpException`
- `battleui/app/Http/Requests/**/*.php` — every FormRequest's `messages()` produces strings, no keys

### The Problem Scenario

```
CLI test → POST /matchmaking/join with mode="garbage"
Laravel FormRequest → 422 {"errors":{"game_mode":["Invalid mode"]}}
                          (no envelope, no error_key)
CLI test → catches → e.error_key === undefined
                  → fallback substring match on e.message
                  → flakes when Laravel locale changes
```

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | High — every new validation rule perpetuates the asymmetry |
| Impact if triggered | Medium — only test reliability and frontend ergonomics; no production exploit |
| Detectability | Low — tests pass on first attempt then flake; no lint catches the omission |
| Current mitigant | Tests substring-match `e.message` on non-engine paths |

---

## Recommended Fix

**Short term:**
- Decide the taxonomy in [ISS-080](ISS-080_20260425_error_key_atd_and_envelope.md). Until that lands, do nothing here.
- Add a TODO comment block in `ApiResponder::error()` referencing this issue and ISS-080 so future contributors don't reinvent.

**Medium term:**
- Subclass `Illuminate\Validation\ValidationException` with an `error_key` property, register a custom renderer in `Handler::render()` that includes it on the envelope.
- Refactor `FormRequest` subclasses to set `protected string $errorKey = "request.validation.<resource>.<field>"`. Either as a per-rule key (pricey, more accurate) or one per request (cheap, less accurate).
- Make all hand-rolled `error()` calls in controllers pass a key as a 4th argument: `$this->error($message, 401, [], ['error_key' => 'auth.credentials.invalid'])`.
- Update `[[api_auth_*]]`, `[[api_matchmaking]]`, `[[api_profile_*]]` atoms to enumerate the keys they emit.

**Long term:**
- Single `ErrorKey` enum shared between Laravel, Go, TypeScript via an atom-driven generator (see ISS-080 long-term).
- Lint that fails any `ApiResponder::error()` call lacking a key.

---

## References

- [ISS-080](ISS-080_20260425_error_key_atd_and_envelope.md) — sister issue for taxonomy + envelope promotion
- `[[api_standard_envelope]]` — envelope contract
- `battleui/app/Traits/ApiResponder.php` — current entry point for envelope construction
- `battleui/app/Http/Controllers/API/AuthController.php` — example of hand-rolled error responses
- `upsiloncli/tests/scenarios/edge_auth_password_policy_full.js` — example of a test that needs typed-key handling
