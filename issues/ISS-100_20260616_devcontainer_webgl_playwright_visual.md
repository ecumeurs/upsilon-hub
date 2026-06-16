# Issue: Devcontainer lost WebGL — Playwright 3D visual specs cannot render

**ID:** `20260616_devcontainer_webgl_playwright_visual`
**Ref:** `ISS-100`
**Date:** 2026-06-16
**Severity:** Medium
**Status:** Open
**Component:** `battleui/tests/playwright/`, devcontainer (Chromium/SwiftShader)
**Affects:** `battle_arena`, `visual_smoke_test`, `battle_arena_sandbox`, `components` specs; any TresJS/three.js render verification.

---

## Summary

Headless Chromium in the current devcontainer **cannot create a WebGL context** — a
direct probe returns `{"webgl":false}` (no SwiftShader software-GL fallback). Because
the arena UI renders through `@tresjs/core` (three.js), the 3D Playwright specs cannot
render or snapshot anything: `battle_arena` finds 0 pawn overlays, `visual_smoke_test`
times out on `__upsilonDebug.ready`, and `battle_arena_sandbox` / `components` hang
waiting for scene state that never populates. This used to work (committed PNG
baselines exist under `tests/playwright/__snapshots__/`, and `playwright.config.ts`
explicitly relies on **SwiftShader** for deterministic rendering) — so this is an
**environment regression**, the same class as the lost `db` host resolution noted
during the 2026-06-16 audit (devcontainer drift).

---

## Technical Description

- `playwright.config.ts` comment: "Three.js renders under SwiftShader (software
  rasterizer), which is deterministic across machines — the right property for
  committed PNG baselines." That assumption no longer holds in this container.
- Probe: launch Chromium via `@playwright/test`, `canvas.getContext('webgl'||'webgl2'||'experimental-webgl')` → all null.
- The functional, DOM-only specs (`user_flows`, `battle_debug`) pass fine — only the
  WebGL/3D-dependent assertions are blocked.

### Rejected workaround (do NOT reintroduce)
A WP-D1 attempt made the specs green by (a) self-skipping WebGL-dependent assertions
and (b) adding a **hidden `opacity:0; aria-hidden` `.pawn-overlay` DOM layer in
production `ThreeGrid.vue`** so `battle_arena` could find pawns without WebGL. This
was reverted: it puts test-only scaffolding in production and makes the spec verify
invisible dummy elements instead of real rendered pawns (hollow coverage).

---

## Recommended Fix

**Short term:** restore software WebGL in the devcontainer — e.g. launch Chromium with
`--use-gl=angle --use-angle=swiftshader` (or `--enable-unsafe-swiftshader` on newer
Chromium) via Playwright `launchOptions.args`, and/or install the SwiftShader libs the
image previously had. Re-probe until `getContext('webgl')` is non-null.
**Medium term:** add a CI/devcontainer smoke check asserting WebGL availability so this
regression is caught early; restore the snapshot baseline run.
**Long term:** the host-GPU-over-SSH escape hatch documented in
[ISS-082](ISS-082_20260425_frontend_playwright_test_seams.md) if SwiftShader
determinism ever diverges.

---

## References
- [ISS-082](ISS-082_20260425_frontend_playwright_test_seams.md) — Playwright suite + seams (reopened; visual coverage blocked by this).
- `battleui/playwright.config.ts` — SwiftShader assumption.
- `reporting/audit/frontend_detailed.md` — frontend audit context.
- Related environment drift: lost `db` host resolution (audit 2026-06-16).
