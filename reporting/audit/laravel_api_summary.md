# Laravel API (battleui) — Summary Report

**System:** battleui PHP — controllers, routes, envelope, auth, matchmaking,
admin, shop/inventory, profile/GDPR
**Date:** 2026-06-16 · **Full detail:** [laravel_api_detailed.md](laravel_api_detailed.md)

## Snapshot
| Dimension | Reading |
|---|---|
| ATD impl tagging | **Strong** — 179 `@spec-link` across `app/` |
| ATD test tagging | **Broken** — only 11 `@test-link`; 18 PHPUnit files invisible to the graph |
| Envelope/contract | **Good** — centralised in `StandardEnvelope` middleware |
| BRD compliance reporting | **Untrustworthy** — CR-05/08/10/11 marked ❌ while their E2E tests pass ✅ |
| Critical open issue | **ISS-093** admin self-anonymize still reachable |

## Top findings
1. **BRD compliance table is decoupled from test execution (major):** four CR
   rows are reported FAIL while the matching E2E scenarios PASS in the *same*
   report. The compliance gate isn't reading the results it summarises.
2. **ISS-093 (Critical, partially open, confirmed):** `delete()` self-guards but
   `anonymize()` only blocks the *last* admin — a non-last admin can self-anonymize
   and overwrite its own credentials. The issue itself notes CI cascade failures.
3. **PHP test traceability collapse:** 179 spec-links vs 11 test-links. Tests
   exist but aren't ATD-linked, so `atd check` shows `NO_TESTS` on real,
   tested requirements (matchmaking, progression, privacy, account).
4. **ISS-088 contradiction:** the credit-economy E2E is a known *test* defect
   (reads attacker, expects target) yet shows ✅ — assertion is likely hollow.

## Architect's commentary

**What is correct / sound.** The envelope discipline is the highlight: a single
`StandardEnvelope` middleware enforces the API contract rather than scattering it
across controllers, which is exactly the right place to later harmonise
`error_key` (ISS-080/081). Domain controller segmentation and FormRequest
validation are clean, conventional Laravel — easy to onboard and extend for V2's
shop/inventory/skill surface.

**What surprised me — good & bad.** *Good:* implementation-side ATD density (179
spec-links) is far higher than I expected for a management layer — someone took
traceability seriously in the controllers. *Bad:* the BRD compliance table
actively contradicts its own E2E results. That is worse than missing coverage —
it is a green/red signal that doesn't track reality, which will let real
regressions through while crying wolf on passing features. I was also surprised a
Critical admin-lockout issue (ISS-093) is only half-fixed and still openly blamed
for CI flakiness.

**What is not appropriate / not scalable.** Two things should not survive to V2:
(1) **per-method admin self-protection** — destructive admin actions must share
one policy, not divergent inline checks; and (2) **runtime self-documentation by
parsing `.atom.md` files** in `HelpController` (ISS-083), which couples the
request path to documentation layout. Underlying both is the same systemic gap as
elsewhere: the *proof* layer (tests, compliance reporting) has drifted away from
the implementation, so the dashboards can no longer be trusted without manual
verification.
