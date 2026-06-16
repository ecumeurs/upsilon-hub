# Go Backend — Summary Report

**System:** `upsilonbattle` (engine) + `upsilonapi` (bridge) + shared Go libs
**Date:** 2026-06-16 · **Full detail:** [go_backend_detailed.md](go_backend_detailed.md)

## Snapshot
| Dimension | Reading |
|---|---|
| ATD tag density | **Strong** — 125 spec-links / 101 test-links in the engine; real bidirectional traceability |
| ATD structural health | **Poor** — 158 atom-level lint defects (battle 102, api 36, types 20) |
| Business-req coverage | Engine binds to mechanics/rules, not business atoms; only `uc_combat_turn` has executing proof |
| Test honesty | **Concern** — security-relevant rules (wrong-controller, friendly-fire, out-of-turn) are SKIP/FAIL, not really covered |
| Open issues verified | ISS-098 (High), ISS-099 (Med), ISS-096 (Med) all **confirmed unfixed in source** |

## Top findings
1. **ISS-098 (High, OPEN, confirmed):** raw User UUID leaked as `player_id` in
   `upsilonapi/api/output.go:208/469/474`. Privacy-requirement violation; blocker
   for public board-state exposure.
2. **Test-quality gap (user's core concern, confirmed):**
   `edge_attack_wrong_controller_with_2` FAILs while the engine guard is correct —
   the *test* is broken, compounded by a shipped error-key typo
   `entity.controller.missmatch`. The matching EC rows are SKIP. The rule works;
   the proof doesn't exist.
3. **ISS-099 (Med, OPEN, confirmed):** `ZoneProperty.Set` silently degrades any
   non-`Single/Neighbours` zone to `Single` — V2 AoE skills cannot work.
4. **ISS-096 (Med, OPEN, confirmed):** traps without `TriggerType` fail silently.
5. **Issue tracking decayed:** ISS-091 & ISS-097 are Resolved but the README
   implies otherwise; **ISS-046 (engine turn-hang, High) is a phantom — the file
   no longer exists** though the README links it.

## Architect's commentary

**What is correct / sound.** The engine is the best-engineered system in the
ecosystem. Package decomposition (`battlearena → ruler → rules/turner/gamestate/
controller`) is clean, the Ruler is split by concern under the length guards, and
the V2 AI archetype controllers (fighter/ranger/support/sneak) extend a shared
base exactly as the milestone intended. ATD tag density here is the real thing —
this team clearly *did* practice traceability where it mattered.

**What surprised me — good & bad.** *Good:* genuine, dense `@spec-link`/
`@test-link` usage and deliberate `gamestate` versioning that gives ISS-054
resurrection a real substrate. *Bad:* the engine advertises "Crash Early" yet the
two newest open issues (ISS-096, ISS-099) are both **silent-fallback** bugs — the
exact opposite of the stated doctrine. And a misspelled error key
(`missmatch`) is shipped on a security path, which simultaneously explains a
failing E2E test and proves the test was never truly asserting the rule.

**What is not appropriate / not scalable.** `output.go` (512 LOC) doing DTO
definition *and* engine→wire mapping *and* harboring the ISS-098 privacy leak is
the structural weak point — it needs splitting with a single identity-masking
seam (ISS-084/085 already anticipate this). More broadly, **the documentation
graph has rotted faster than the code**: 158 lint defects and a phantom
High-severity issue mean the ATD layer can no longer be trusted as ground truth
without the remediation pass — a serious problem for a project whose entire
thesis is "absolute traceability."
