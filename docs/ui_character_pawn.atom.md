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
- **Body:** A solid 3D hexagonal pyramid composed of 6 faceted planes that meet at a shared top vertex.
- **Rotation Logic:** 
  - **Inactive:** Static at an idle angle (30deg) to provide a consistent silhouette.
  - **Active:** Rotates continuously on the Y-axis to indicate it is the acting character.
- **Head:** A sphere with a slight indent indicating character facing direction.
- **Holographic Effect:** Semi-transparent gradient facets with scanline overlay and intermittent static/glitch animation.
- **Color Coding:** Blue (current player), Green (ally), Red (enemy player 1), Purple (enemy player 2). Distinct shade variations via HSL shift.
- **Pawn Interactivity:** Pointer events pass through to underlying tiles.

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[ui_character_pawn]]`
- **Component:** `CharacterPawn.vue`
- **Props:** `entity`, `teamColor`, `isActive` (bool), `shadeOffset` (int)

## EXPECTATION
- Pawn renders as a visible holographic figure on the board.
- Team colors are clearly distinguishable.
- Active character has enhanced visual feedback.
- Holographic glitch effect fires intermittently.
- **Pointer Events:** Pawn does not capture mouse clicks, allowing players to select the tile underneath.
