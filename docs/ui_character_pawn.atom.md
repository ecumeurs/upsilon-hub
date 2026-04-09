---
id: ui_character_pawn
status: STABLE
priority: 5
tags: [ui, combat, pawn, holographic, 3d]
parents:
  - [[ui_iso_board]]
dependents: []
type: UI
layer: ARCHITECTURE
version: 1.0
human_name: Character Pawn UI
---

# New Atom

## INTENT
The holographic 3D character pawn rendered on the isometric board — a faceted rotating cone body with a directional sphere head, floating name label, and HP bar.

## THE RULE / LOGIC
- **Body:** A 5-6 sided inverted cone shape (CSS polygon), slowly rotating.
- **Head:** A sphere with a slight indent indicating character facing direction.
- **Holographic Effect:** Semi-transparent with scanline overlay, intermittent static/glitch animation.
- **Color Coding:** Blue (current player), Green (ally), Red (enemy player 1), Purple (enemy player 2). Distinct shade variations per character via HSL shift.
- **Overlay:** Floating name label above the pawn + filled HP bar (current/max HP).
- **Active State:** Brighter glow and pulsing animation when it's this character's turn.

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[ui_character_pawn]]`
- **Component:** `CharacterPawn.vue`
- **Props:** `entity`, `teamColor`, `isActive` (bool), `shadeOffset` (int)

## EXPECTATION
- Pawn renders as a visible holographic figure on the board.
- Team colors are clearly distinguishable.
- Active character has enhanced visual feedback.
- Holographic glitch effect fires intermittently.
