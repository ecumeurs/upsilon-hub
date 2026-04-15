# Issue: Holo-Emote Procedural Reaction System

**ID:** `20260415_holo_emote_system`
**Ref:** `ISS-039`
**Date:** 2026-04-15
**Severity:** Medium
**Status:** Open
**Component:** `battleui` (Frontend), `battleui` (Backend/User Profile)
**Affects:** `CharacterPawn.vue`, `BattleArena.vue`, `User` Model

---

## Summary

Implement a "Holo-Emote System" that triggers procedural reactions (emojis/text) above gladiator pawns during battle. The system uses a "pure CSS" holographic aesthetic that scales in visual quality (Tiers 1-5) based on the player's profile progression.

---

## Technical Description

### 1. Visual Foundation & Pawn Attachment
The emotes are spatial HUD overlays emitted by the gladiators' holographic projectors.
* **Attachment:** Anchored to the top of the "pawn" figure in `CharacterPawn.vue`.
* **Animations:** "Float" animation (out of sync with pawn), "broken" jitter using `steps()`, and "leak" neon glow via `text-shadow`.
* **Shapes/Textures:** Aggressive polygons (`clip-path`), scanlines (`linear-gradients`).

### 2. The Tech-Hierarchy (Level Progression)

| Tier | Name | Visual Quality | CSS Key Features |
| :--- | :--- | :--- | :--- |
| **Tier 1** | **Scrap** | Monochrome, unstable, heavy flickering. | `opacity` jitter, `grayscale(1)`, 1-bit color. |
| **Tier 2** | **Wasteland** | Multi-state animations (like neon). | `steps(2)` animations, 2-color palette. |
| **Tier 3** | **Syndicate** | High-energy, "bleeding" colors. | Chromatic aberration (`text-shadow` offsets). |
| **Tier 4** | **Elite** | Structured HUD, scanline masks. | `mask-image`, complex `clip-path` shapes. |
| **Tier 5** | **Legendary** | Stable, high-refresh, ghosting. | `mix-blend-mode: screen`, smooth `blur`. |

### 3. Data Structure (Laravel Backend)
The `User` model requires new fields to drive the logic:
* `emitter_tier` (Int): Max CSS complexity level.
* `unlocked_palettes` (Array/JSON): Neon hex codes.
* `equipped_emote_set` (JSON): Mappings for reaction triggers.

### 4. Procedural Reaction Logic
Emotes trigger based on a "Proc Rate" (random chance) within a **3-tile radius** ("Brawl Zone").

| Event | Emote | Proc Chance | Context / Vibe |
| :--- | :--- | :--- | :--- |
| **Takes Damage** | `SHOCKED` | 30% | High-contrast flicker; yellow/orange neon. |
| **Ally Hurt** | `COMPASSIONATE` | 20% | cyan "glitch heart". |
| **Enemy Hurt** | `JEERING` | 40% | toxic green "haha" or "L". |
| **Ally Deals DMG** | `CHEERING` | 20% | Bright gold "LFG" or "W". |
| **Ally Passes** | `BOOING` | 20% | Red desaturated "COWARD". |
| **Kill (Finisher)** | `WASTED` | 60% | Forces Tier +2 (when tier 1: Chromatic Aberration). |

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | Medium |
| Impact if triggered | Low (Visual/Flavor) |
| Detectability | High |
| Current mitigant | None |

---

## Recommended Fix

**Short term:** Update `User` model and migrations. Implement basic CSS for Tiers 1-2.  Add dedicated profile section to manage this system (equiping emotes, palettes, etc).
**Medium term:** Implement the proc logic in the Frontend, hooked into the `ActionFeedback` protocol (ISS-038). Add a dedicated admin page to manage this system.
**Long term:** Add randomized text/symbol variations and more complex "Legendary" effects.

---

## References

- [battleui/README.md](file:///workspace/battleui/README.md)
- [ISS-038_20260415_action_feedback_protocol.md](file:///workspace/issues/ISS-038_20260415_action_feedback_protocol.md)
- [req_ui_look_and_feel.atom.md](file:///workspace/docs/req_ui_look_and_feel.atom.md)
- [CharacterPawn.vue](file:///workspace/battleui/resources/js/Components/Arena/CharacterPawn.vue)
