# Issue: Upgradable Pawn Appearance & Model System

**ID:** `20260415_pawn_appearance_system`
**Ref:** `ISS-040`
**Date:** 2026-04-15
**Severity:** Medium
**Status:** Open
**Component:** `battleui` (Frontend), `battleui` (Backend/User Profile)
**Affects:** `CharacterPawn.vue`, `BattleArena.vue`, `User` Model

---

## Summary

Implement an upgradable "Pawn Appearance System" that allows players to customize their gladiator's physical representation in the arena. Similar to the Holo-Emote system (ISS-039), this system provides tiers of holographic models, selectable color palettes, and ambient visual effects that evolve as the player progresses.

---

## Technical Description

### 1. Visual Customization Layers
The pawn's appearance is divided into three primary layers, all managed via CSS/SVG and potentially lightweight GL components in the future.
*   **Model Tier:** The base geometric shape of the pawn (e.g., from a basic cone/sphere to complex faceted cyber-gladiators).
*   **Color Palette:** A set of CSS variables defining the holographic glow, facet colors, and scanline tints.
*   **Environmental Effects:** Persistent "aura" or "glitch" animations surrounding the pawn (e.g., drift particles, data-leak shadows).

### 2. The Appearance Hierarchy (Progression)

| Tier | Name | Visual Complexity | Key Features |
| :--- | :--- | :--- | :--- |
| **Tier 1** | **Fragment** | Wireframe / Basic Geoms | Minimal facets, high transparency, monochrome. |
| **Tier 2** | **Drifter** | Low-poly faceted | Basic solid colors, static scanlines. |
| **Tier 3** | **Enforcer** | Structured Armor Plates | Multi-color palettes, pulse animations. |
| **Tier 4** | **Centurion** | Dense mesh / Refractions | Reflective surfaces, chromatic aberration. |
| **Tier 5** | **Ascendant** | Particle-integrated | Aura effects, custom model geometry, stable/high-refresh. |

### 3. Data Structure (Laravel Backend)
The `User` model (or a related `Customization` model) needs to track availability and current selections:
*   `active_pawn_model` (String/ID)
*   `active_pawn_palette` (String/ID)
*   `active_pawn_effect` (String/ID)
*   `unlocked_appearance_items` (JSON): A collection of IDs for models, palettes, and effects.

### 4. Administrative Requirements
*   **Availability Management:** Admins must be able to flag certain models/colors as "Seasonal", "Ranked Reward", or "Store Item".
*   **Preview Mode:** A way to test CSS-driven pawn configurations before releasing them to the general population.

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | Medium |
| Impact if triggered | Medium (Visual clutter, performance overhead on many pawns) |
| Detectability | High |
| Current mitigant | Existing `ui_character_pawn` atom provides a baseline for expansion. |

---

## Recommended Fix

**Short term:** Update `User` model to support appearance fields. Expand `CharacterPawn.vue` to accept `model_tier` and `palette` props.
**Medium term:** Implement a "Customization" tab in the user profile for live preview and selection.
**Long term:** Build an Admin dashboard for managing the "Appearance Store" and creating new tiered content without engine code changes.

---

## References

- [ISS-039_20260415_holo_emote_system.md](file:///workspace/issues/ISS-039_20260415_holo_emote_system.md)
- [ui_character_pawn.atom.md](file:///workspace/docs/ui_character_pawn.atom.md)
- [CharacterPawn.vue](file:///workspace/battleui/resources/js/Components/Arena/CharacterPawn.vue)
- [User.php](file:///workspace/battleui/app/Models/User.php)
