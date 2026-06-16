# Laravel API (battleui) — Detailed Investigation Report

**Auditor:** Principal Systems Architect · **Date:** 2026-06-16
**Scope:** `battleui/app/Http/**`, `battleui/routes/**`, `battleui/tests/**`,
`battleui/docs/*.atom.md` (the API/UI atom home)
**Method:** static analysis + `atd` CLI. No suites run; CI status from
`reporting/CI_report.md` (stale, 2026-05-12) reasoned against current source.

---

## 1. Documentation coverage

| Metric | Value | Note |
|---|---:|---|
| battleui atoms | 63 | API + UI atoms live in `battleui/docs/`, **not** root `docs/` |
| `@spec-link` in PHP (`app/`) | **179** | very dense implementation tagging |
| `@test-link` in PHP (`tests/`,`app/`) | **11** | **near-absent test traceability** |
| PHPUnit test files | 18 | Feature + Unit suites genuinely exist |
| Lint defects (battleui) | 47 | missing `## EXPECTATION`, alias links, etc. |

**The core documentation defect here is the spec/test asymmetry.** 179 spec-links
vs 11 test-links is why `atd check` reports `NO_TESTS` on
`requirement_customer_user_account`, `requirement_customer_user_id_privacy`,
`rule_matchmaking_single_queue`, `rule_progression`, `uc_admin_login` — **the
tests exist** (`PVEMatchmakingTest`, `ExtraMatchmakingTest`, `MatchVerificationTest`,
`GdprTest`, `LeaderboardTest`, `ErrorHandlingTest`, …) but are not ATD-linked, so
the traceability graph cannot see them. Coverage tooling under-reports PHP truth.

## 2. Business-requirement & test reality

### Reporting-integrity failure in the BRD Compliance table (major)
The CI report's "BRD Compliance Checks" marks four rows ❌ **while the matching
E2E scenarios in the very same report pass ✅**:

| BRD row | BRD status | Same-report E2E scenario | E2E status |
|---|:--:|---|:--:|
| CR-05 Matchmaking PvP (Queue) | ❌ | `e2e_matchmaking_pvp_queue_with_2` | ✅ |
| CR-08 Match Resolution (Standard) | ❌ | `e2e_match_resolution_standard_with_2` | ✅ |
| CR-10 Progression (Post-Win) | ❌ | `e2e_progression_post_win_with_2` | ✅ |
| CR-11 Progression Constraints | ❌ | `e2e_progression_constraints_with_2` | ✅ |

The BRD compliance gate is **decoupled from actual test execution** — it is not
reading the E2E results it claims to summarise. This makes the compliance section
untrustworthy in both directions (false negatives here; potentially false
positives elsewhere). This is the reporting-side twin of the engine's
"tests don't assert their case" problem.

### Genuinely failing E2E (Laravel-owned)
- `edge_match_queue_while_in_match_with_2` ❌ — the single-queue guard
  (`rule_matchmaking_single_queue`) round-trip. `MatchMakingController` has the
  pieces (`joinMatch`, zombie-queue cleanup at `:130`) but the EC equivalent
  (EC-32) is SKIP, so this guard also lacks reliable proof.

## 3. Issue reconciliation (Laravel-owned)

| Ref | Severity | Reality in source | True state |
|---|---|---|---|
| **ISS-093** admin self-destruction | Critical | **Partially open.** `AdminController::delete` self-guards (`:162` `if ($user->id === auth()->id())`). `AdminController::anonymize` (`:147`) guards **only** "last remaining administrator" (`:150`) — a non-last admin can still anonymize **itself**, overwriting its own credentials. The destructive path the issue describes is still reachable. The issue notes this is "causing cascading failures in the CI suite." | **Open — valid, Critical** for `anonymize`. |
| **ISS-092** skill property sync | High | Cross-service `PropertyDTO` lacks fields for complex skill props (`Range`, `TargetType`, `TargetingMechanics`); engine `Set()` empty. Pairs with engine ISS-099. | Open — valid. |
| **ISS-088** credit economy payload mismatch | Med | `e2e_credit_economy.js` reads attacker state where it expects target state. **Test defect, not API defect** — yet `e2e_credit_economy` shows ✅ in CI, so either the test was patched without closing the issue, or the assertion is hollow. | Open — needs reconciliation. |
| **ISS-083** automate API help | Med | `HelpController` still hybrid reflection + `.atom.md` parsing. | Open — valid. |
| **ISS-090** action endpoint segregation | High | (summary body sparse in file.) | Needs scoping. |

## 4. Architecture & scalability

**Sound**
- **Centralised envelope**: `app/Http/Middleware/StandardEnvelope.php` enforces
  the `[[api_standard_envelope]]` contract in one place rather than per-controller
  — correct, and the right seam for ISS-080/081 `error_key` harmonisation.
- Controllers are cleanly segmented by domain (`API/*Controller.php`,
  `API/Admin/*`), with FormRequest validation (`JoinMatchRequest`).
- Auth uses the standard Laravel Auth scaffolding (`Auth/*Controller.php`) plus
  JWT/session bridge — conventional and maintainable.

**Risk / not scalable**
- **Test traceability is effectively absent** (11 links). For a project whose
  thesis is end-to-end traceability, the PHP layer's *proof* edge is missing; CI's
  "PHPUnit … Unknown" summary is a direct symptom.
- **Admin destructive actions lack a uniform self-protection policy** — `delete`
  and `anonymize` guard differently. This should be one policy/middleware, not
  per-method ad-hoc checks (ISS-093).
- **`HelpController` self-documentation by parsing atom files at runtime**
  (ISS-083) couples the request path to doc-file layout — fragile and won't scale
  with the V2 surface.

## 5. Recommendations (no code changed here)
1. **Fix the BRD compliance generator** to read real E2E results (or delete the
   section) — it is currently misinforming stakeholders.
2. Close ISS-093 by routing all admin destructive actions through one
   self-and-last-admin guard.
3. Backfill `@test-link` tags onto the 18 existing PHPUnit files so PHP coverage
   becomes visible (cheap, high leverage).
4. Reconcile ISS-088: make `e2e_credit_economy` assert target state, then either
   close the issue or keep it open with the real defect.
