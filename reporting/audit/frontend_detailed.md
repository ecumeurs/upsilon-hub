# Frontend (battleui Vue) — Detailed Investigation Report

**Auditor:** Principal Systems Architect · **Date:** 2026-06-16
**Scope:** `battleui/resources/js/**` (Vue 3 + Inertia + Tailwind), Reverb WS
client, `tests/playwright/**`, UI atoms in `battleui/docs/`
**Method:** static analysis + `atd` CLI; CI status from `reporting/CI_report.md`.

---

## 1. Documentation coverage

| Metric | Value | Note |
|---|---:|---|
| JS/Vue source files | 98 | Pages, Layouts, Components/{Shared,Character,Arena,Dashboard} |
| `@spec-link` in `resources/js` | **90** | strong UI→atom tagging |
| `@test-link` in `resources/js` | **0** | **no test traceability whatsoever** |
| UI-type atoms (workspace) | 72 (47 `ui_*` files) | |
| Playwright spec files | 6 | `battle_arena`, `components`, `user_flows`, `visual_smoke_test`, `battle_arena_sandbox`, `battle_debug` |

The frontend mirrors the backend pattern but more extreme: solid implementation
tagging (90 spec-links), **zero** test linkage. The Vue components are
ATD-annotated to UI atoms, but nothing connects those atoms to the Playwright
specs that supposedly verify them.

## 2. Business-requirement & test reality

- **Playwright works locally but is invisible to CI.** `playwright_last_run.log`
  shows a green run — but only **1 test executed** ("1 passed"), and
  `reporting/CI_report.md` ends with **"❌ No HTML Report found."** So the suite
  is runnable yet (a) barely exercised and (b) not captured by CI. The six spec
  files are not being driven as a gate.
- **ISS-082 (Playwright test seams) is marked Resolved**, but the realised state
  is a single passing smoke test, not the component-isolation visual-regression
  suite the issue scopes. The resolution is aspirational relative to what runs.
- UI use-case atoms (`uc_player_login`, `uc_matchmaking`, `us_character_reroll`,
  session timeout) have spec-links into components but **no executing UI proof** —
  customer journeys are asserted only indirectly via the `upsiloncli` E2E layer.

## 3. ISS-098 privacy leak — frontend blast radius

The engine's raw-UUID `player_id` leak (engine report §3) is **consumed directly
in the client**:
- `BattleArena.vue:39` `String(myPlayer.value.player_id || myPlayer.value.nickname)`
- `BattleArena.vue:212` entity id derived from `p.player_id || p.nickname`
- `BattleArena.vue:237` foe identity compared on `player_id`
- `BattleArena.vue:244` `colors[String(p.player_id)] = color` — UUID used as a map key

There is a `|| nickname` fallback (defensive), but whenever the engine sends a
real `player_id`, the **raw database UUID reaches the browser DOM/state**. Fixing
ISS-098 at the engine seam (single mask) also closes the frontend exposure with no
client change required — another argument for the `output.go` masking seam.

## 4. Architecture & scalability

**Sound**
- Sensible component taxonomy: `Layouts/`, `Pages/`, `Components/Shared`,
  `Components/Character`, `Components/Arena`, `Components/Dashboard/panels`.
- Reusable primitives (`PrimaryButton`, `Modal`, `TextInput`, `InputError`) —
  conventional, themeable, Tailwind-driven.
- Dedicated sandbox/test pages (`BattleArenaSandbox`, `TestComponent*`) give real
  isolation seams for visual testing — the right foundation for ISS-082.

**Risk / not scalable**
- **`BattleArena.vue` = 846 LOC — over the 600 ATD *error* threshold** (and
  `Components/Dashboard/panels/CharacterPanel.vue` = 658). These are god-components
  mixing WS wiring, board projection, identity/color logic, and action dispatch.
- **ISS-084 ("Refactor BattleArena into components") is marked Resolved while
  `BattleArena.vue` remains the largest file at 846 LOC.** Either the split
  regressed or the issue was closed prematurely — a status-integrity defect.
- **Zero UI test traceability** means any V2 UI work (skill modals, equipment
  slots, shop) ships without a regression net visible to the graph.
- `ActionPanel.vue` (544) and `IsoBoardGrid.vue` (421) are over the 400 warn line —
  expected for an isometric board, but they will absorb the V2 skill/channeling/
  backstab UI and should be split pre-emptively.

## 5. Recommendations (no code changed here)
1. Wire Playwright into CI with report capture; drive all 6 specs, not 1.
2. Re-open or re-scope ISS-084 — `BattleArena.vue` is not actually decomposed.
3. Split `BattleArena.vue` along WS / board / identity / action seams before V2.
4. Add `@test-link` tags from UI atoms to the Playwright specs so UI coverage
   becomes measurable.
5. Treat the frontend `player_id` consumption as resolved-by-engine (ISS-098).
