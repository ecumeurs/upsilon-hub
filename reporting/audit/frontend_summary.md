# Frontend (battleui Vue) — Summary Report

**System:** Vue 3 + Inertia + Tailwind SPA, Reverb WebSocket client, Playwright suite
**Date:** 2026-06-16 · **Full detail:** [frontend_detailed.md](frontend_detailed.md)

## Snapshot
| Dimension | Reading |
|---|---|
| UI→atom tagging | **Strong** — 90 `@spec-link` across components |
| UI test traceability | **None** — 0 `@test-link` |
| Playwright in CI | **Effectively absent** — runs locally (1 test), "No HTML Report found" in CI |
| Component health | `BattleArena.vue` **846 LOC** (> 600 error), `CharacterPanel.vue` 658 |
| Privacy exposure | Consumes leaked raw `player_id` UUID (ISS-098) |

## Top findings
1. **Playwright is a paper tiger:** the suite exists and passes locally, but only
   one test runs and CI captures no report. ISS-082 is marked Resolved while the
   real state is a single smoke test.
2. **ISS-084 status integrity failure:** "Refactor BattleArena into components"
   is Resolved, yet `BattleArena.vue` is still the largest file at 846 LOC — over
   the ATD 600-LOC error limit.
3. **Zero UI test traceability:** 90 spec-links, 0 test-links — UI atoms connect to
   components but to no proof.
4. **ISS-098 reaches the browser:** raw UUID `player_id` is used as entity id and
   color key in `BattleArena.vue` (with a nickname fallback). Fixed for free once
   the engine masks at `output.go`.

## Architect's commentary

**What is correct / sound.** The component taxonomy is clean and conventional —
Layouts / Pages / Shared primitives / domain folders (Arena, Character,
Dashboard) — and the reusable Tailwind primitives are exactly what you want for a
growing V2 surface. Crucially, the team built dedicated sandbox and test pages
(`BattleArenaSandbox`, `TestComponent*`), which are genuine isolation seams for
visual regression. The raw materials for a strong UI test story are present.

**What surprised me — good & bad.** *Good:* 90 UI spec-links is far more
documentation discipline than typical frontends carry. *Bad:* the gap between
claimed and actual testing is the widest of any system here. A passing
`playwright_last_run.log` with exactly one test, a CI line reading "No HTML Report
found," and an ISS-082 marked Resolved together paint a picture of testing
theatre — the infrastructure exists but isn't doing the job, and the issue tracker
says otherwise.

**What is not appropriate / not scalable.** `BattleArena.vue` at 846 LOC is the
clearest "must not survive to V2" item — it concentrates WebSocket wiring, board
projection, identity/color logic and action dispatch in one file, breaches the
project's own length-error guard, and is falsely marked as already-refactored
(ISS-084). V2 will pile skill modals, channeling indicators, and backstab feedback
onto exactly this component. It must be decomposed and given real UI tests before
that weight lands — otherwise the most player-visible surface becomes the least
verifiable.
