# Issue: Refactor BattleArena into components and restore visual effects

**ID:** `20260425_component_split_effects_plan`
**Ref:** `ISS-084`
**Date:** 2026-04-25
**Severity:** Medium
**Status:** Open
**Component:** `battleui/resources/js/Pages/BattleArena.vue`
**Affects:** `battleui/resources/js/Components/Arena/*`

---

## Summary

Refactor the arena UI by splitting `ThreeGrid.vue` into dedicated components (Pawn, Obstacle, Cell, Grid) and re‑introduce the hologram, glitch, and CRT scanline effects that were lost during the TresJS migration. Provide a global `effects=true` query‑parameter toggle to enable all effects for visual‑review testing.

---

## Technical Description

### Background

The current `BattleArena.vue` page uses a monolithic `ThreeGrid.vue` implementation that renders tiles, obstacles, and pawns inline. After moving to TresJS, several visual effects (hologram transparency on pawns/obstacles and CRT‑style scanlines) were removed, breaking visual‑review tests.

### The Problem Scenario

1. `ThreeGrid.vue` contains large loops that directly create `<TresMesh>` elements.
2. No per‑component effect handling exists; the hologram shader and scanline overlay are missing.
3. Playwright visual‑review tests fail because the page layout differs from the original design.

```text
[Current]  -> Single ThreeGrid with inline meshes (no effects)
[Desired]  -> ArenaGrid component hosting ArenaCell, ArenaPawn, ArenaObstacle
            + HologramMaterial for pawns/obstacles
            + CRT overlay component
            + Global `effects` toggle
```

### Where This Pattern Exists Today

- `/workspace/battleui/resources/js/Components/Arena/ThreeGrid.vue`
- `/workspace/battleui/resources/js/Pages/BattleArena.vue`

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | Medium |
| Impact if triggered | Medium |
| Detectability | High – console errors or missing effects will be obvious |
| Current mitigant | None – effects are currently absent |

---

## Recommended Fix

**Short term:**
- Create new component files (`ArenaGrid.vue`, `ArenaCell.vue`, `ArenaPawn.vue`, `ArenaObstacle.vue`).
- Move rendering logic from `ThreeGrid.vue` into these components.
- Add a global `effects` query‑parameter toggle and pass it down.
- Implement a simple hologram shader using `@tresjs/shader-material` and a CRT overlay component with CSS.

**Medium term:**
- Add ATD atoms for each new component and effect for traceability.
- Update Playwright tests to include a visual‑review mode (`?effects=true`).

**Long term:**
- Evaluate performance and add fallback materials for low‑end devices.

---
## Technical Lessons Learned (Post-Mortem)

### 1. Library Versioning & Hook Naming
*   **Issue:** `useRenderLoop` (often cited in docs) was unavailable in `@tresjs/core` v5.8.0.
*   **Lesson:** Use **`useLoop()`** instead. The specific update hook is **`onRender`**, not `onLoop`.

### 2. Post-Processing Configuration
*   **Issue:** Manual `EffectComposer` implementation is fragile and conflicts with the standard TresJS render loop.
*   **Lesson:** Use the **`@tresjs/post-processing`** package. Note that component names differ from Three.js classes: use **`<UnrealBloom />`** (not `Bloom`) and **`<EffectComposer />`**.

### 3. Stability of HTML Overlays
*   **Issue:** `<Html>` from `@tresjs/cientos` causes a `MutationObserver` crash if initialized before the canvas parent is stable in the DOM.
*   **Lesson:** Implement a **Double-Guard**:
    - Parent (`ThreeGrid`) must signal `ready` after its own `onMounted`.
    - Child (`Pawn3D`) must wait for both `onMounted` and the parent's `ready` signal before rendering `<Html>`.

### 4. Materials & Shaders
*   **Shader Material:** `@tresjs/shader-material` was not found/available. Manual **`<TresShaderMaterial />`** is the correct path but requires careful uniform management via `useLoop`.
*   **Shadows:** `PCFSoftShadowMap` is deprecated in Three.js r184; use **`THREE.PCFShadowMap`** to clear console warnings.

### 5. Architectural Improvements
*   The split of `ThreeGrid.vue` into specialized sub-components (`Tile3D`, `Obstacle3D`, `Pawn3D`) significantly reduced cognitive load and is recommended for any future attempts.

---
## Added informations

### Effects Details

| Effect | Applied To | Description |
|---|---|---|
| Hologram Shader | `Pawn3D.vue` | Custom shader material with scrolling scanlines and 20Hz flicker. Uses Additive Blending. |
| Neon Bloom | `PostProcess.vue` | UnrealBloomPass via `@tresjs/post-processing` with intensity 1.5. |
| Wireframe Highlights | `ThreeGrid.vue` | Flat planes with `wireframe: true` for tactical cells. |
| Camera-Facing Labels | `Pawn3D.vue` | Non-transformed `<Html>` overlays for Name/HP. |
---

## References

- `/workspace/battleui/resources/js/Components/Arena/ThreeGrid.vue`
- `/workspace/battleui/resources/js/Pages/BattleArena.vue`
- ATD documentation for UI atoms.
