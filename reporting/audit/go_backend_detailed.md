# Go Backend — Detailed Investigation Report

**Auditor:** Principal Systems Architect
**Date:** 2026-06-16
**Scope:** `upsilonbattle` (core engine), `upsilonapi` (bridge), shared Go libs
(`upsilontypes`, `upsilonmapdata`, `upsilonmapmaker`, `upsilonserializer`,
`upsilontools`)
**Method:** static analysis + `atd` CLI (graph + semantic index rebuilt this
session; MCP server was down — see INDEX.md). No suites executed; pass/fail taken
from `reporting/CI_report.md` (commit `d29fd1b`, 2026-05-12) and reasoned against
current working-tree source.

---

## 1. Documentation coverage

| Project | Atoms | `@spec-link` (src) | `@test-link` | Lint defects | Non-test LOC | Test files |
|---|---:|---:|---:|---:|---:|---:|
| upsilonbattle | 130 | 125 | 101 | 102 | 6,250 | 44 |
| upsilonapi | 63 | 42 | 24 | 36 | 2,747 | 11 |
| upsilontypes | 26 | 18 | 13 | 20 | — | — |

**Observations**
- **Traceability density is genuinely high** in the engine: 125 `@spec-link` and
  101 `@test-link` tags across `upsilonbattle` is real, bidirectional ATD usage,
  not decoration. This is the healthiest system in the repo on tag density.
- **Structural lint is broken at scale**: 102 + 36 + 20 = **158 atom-level lint
  defects** in the Go tree alone. Recurring classes:
  - Missing mandatory sections (`## EXPECTATION`, `## TECHNICAL INTERFACE`).
  - `Unresolved parent link` to a short alias that is not an atom id — e.g.
    `[[upsilon_vision]]` (real id `vision_upsilon_vision`),
    `[[shared:req_tech_debt_backlog]]` (no such atom),
    `[[upsilonbattle:entity_player]]` cross-project miss.
  - `Traceability Gap: @spec-link present but 0 @test-link (Missing Proof)` on
    `mechanic_mech_ai_termination`, `mechanic_randomization_helpers`,
    `mechanic_math_core_utils`, `mechanic_mech_battle_engine_stress_testing`.
- **`atd check` cross-project**: the engine implements almost none of the shared
  *business* atoms directly (expected — those bind to Laravel/API). Only
  `uc_combat_turn` shows test coverage (7 test links) from the engine side.
- **Orphaned STABLE atoms (`crawl --gaps`)** touching the engine:
  `upsilonbattle:rule_credit_action_communication_layer`,
  `upsilonbattle:mechanic_mech_battle_startup_handshake` — declared STABLE but no
  implementation link resolves.

## 2. Business-requirement & test reality

### CI baseline failures attributable to the engine boundary
| CI item | Status | Finding |
|---|---|---|
| `edge_attack_wrong_controller_with_2` | ❌ FAIL | **Engine logic is correct** — `ruler/rules/attack_checks.go:25` calls `CheckControllerForEntity` and returns `entity.controller.missmatch`. The failure is in the E2E scenario, not the rule. The error key is misspelled (`missmatch`); any assertion expecting `mismatch` would fail. This is the user's "tests don't cover their own case" concern made concrete. |
| EC-11 "Attack Wrong Controller" | ⏭️ SKIP (edge report) | Same rule, asserted nowhere. The engine guard exists but has **no executing proof** — SKIP in the edge matrix and FAIL in E2E means zero real coverage of a security-relevant rule. |
| EC-10/EC-12/EC-19 (out-of-turn, friendly-fire, targeting) | ⏭️ SKIP | Rules exist in `ruler/rules/` but the EC harness skips them. Claimed "covered" in narrative, unproven in fact. |

### Engine guards verified present in source
- **Movement (ISS-091, now Resolved)**: `rules/move.go:51` captures the
  `Grid.MoveEntity` error (previously ignored) and `:174` returns
  `entity.path.obstacle`. Fix is in the working tree; README still mislabels it.
- **Controller identity**: present (above).

## 3. Issue reconciliation (engine-owned)

| Ref | File status | Reality in source | True state |
|---|---|---|---|
| **ISS-098** internal user-id leak | Open (High) | **Confirmed unfixed.** `upsilonapi/api/output.go:208` `PlayerID: ent.ControllerID.String()`, plus `:469`/`:474` map raw `ControllerID` (User UUID) into `player_id` / `current_player_id`. Directly violates `[[requirement_customer_user_id_privacy]]`. | **Open — valid, High.** |
| **ISS-099** zone support gap | Open (Med) | **Confirmed unfixed.** `upsilontypes/property/def/skill.go:177` `ZoneProperty.Set` handles only `"Single"`/`"Neighbours"` and silently falls back to `Single` (`:187`) for anything else (`Circle:3`, `Square:2`). Blocks V2 AoE. | **Open — valid, Med.** |
| **ISS-096** trap trigger enforcement | Open (Med) | Matches source: `rules/positionaleffect.go` `processSinglePositionalEffect` returns early & silently when `TriggerType` property is absent. Violates "Crash Early". | **Open — valid, Med.** |
| **ISS-091** movement bypass | Resolved | Fix present (§2). | **Resolved — but README implies Open.** |
| **ISS-097** actor stop panic race | Resolved | Idempotent `Stop()` fix described; affects `upsilontools/tools/actor`. | **Resolved — README stale.** |
| **ISS-046** turner hands turn to dead entity | Listed Open (High) in README | **File does not exist** in `issues/`. Phantom row — either resolved-and-deleted without README update, or lost. Turner hang is a battle-stopping defect; its tracking is gone. | **Unknown — tracking lost.** |
| ISS-065/066/067/070/071/073/078 (V2 wave) | mixed | AI archetypes (`controller/archetype/*.go`: fighter/ranger/support/sneak) and behavior pipeline are implemented; credit/backstab logic present in `ruler/rules`. Status drift: feature atoms still DRAFT despite STABLE implementations (see Documentation Remediation). | Largely implemented. |

## 4. Architecture & scalability

**Sound**
- Clean package decomposition: `battlearena → ruler → {rules, turner, gamestate,
  controller}`. The Ruler is split by concern (`ruler_turn`, `ruler_victory`,
  `ruler_shotclock`, `ruler_lifecycle`) — good cohesion, mostly under length guards.
- Archetype controllers extend a shared `archetype.go` base cleanly (V2 AI goal met).
- `gamestate/gamestate_version.go` shows deliberate versioning of engine state —
  the substrate ISS-054 resurrection depends on.

**Risk / not scalable**
- **`upsilonapi/api/output.go` = 512 LOC** (> 400 ATD warn). It is both the DTO
  schema *and* the engine→wire mapping, and it is where the ISS-098 privacy leak
  lives. This file should be split (DTO vs mapper) and given a single masking
  seam — see ISS-084/085 intent.
- **`controller/controllers/aggressive.go` = 430 LOC** — over warn; AI decision
  logic concentrated in one file.
- **Error-key hygiene**: `entity.controller.missmatch` typo is shipped. Error keys
  are part of the cross-stack contract (ISS-080/081); a typo here silently breaks
  any consumer or test matching the correct spelling.
- **Silent-fallback anti-pattern** recurs (ISS-096 trap, ISS-099 zone): the engine
  claims "Crash Early" (`.agent/rules/COMMON.md`) but degrades silently in these
  paths — the opposite of the stated standard.

## 5. Recommendations (no code changed here)
1. Treat **ISS-098** as a release blocker for any public board-state exposure;
   introduce a single masking function in the `output.go` mapper.
2. Fix the `missmatch` error key and reconcile it with the E2E assertion; then
   un-SKIP EC-10/11/12/19 so the security rules have executing proof.
3. Implement `ZoneProperty` pattern parsing (ISS-099) before V2 AoE skills.
4. Replace silent fallbacks (ISS-096/099) with explicit errors per Crash-Early.
5. Resolve the lint defect classes (alias links, missing sections) — see remediation.
