---
id: ui_iso_board
status: STABLE
type: UI
layer: ARCHITECTURE
version: 1.0
tags: [ui, combat, board, isometric, 3d]
dependents:
  - [[ui_holo_obstacle]]
  - [[ui_character_pawn]]
human_name: Isometric Board Grid UI
priority: 5
parents:
  - [[ui_battle_arena]]
  - [[mech_board_generation_board_dimensions]]
---

# New Atom

## INTENT
The 3D isometric grid board renderer displaying a dynamic NxM tile grid with obstacles, character pawns, movement range highlights, and attack target highlights.

## THE RULE / LOGIC
- **Rendering:** Pure CSS 3D isometric transform (rotateX/rotateZ) applied to the grid container.
- **Grid Size:** Dynamic based on `width` and `height` props (5-15 tiles per axis per ATD spec).
- **Tile Types:** 
    - Normal (dark gunmetal)
    - Obstacle (rusty/blocked with texture)
    - Move-Highlighted (cyan glow)
    - Attack-Highlighted (magenta glow)
    - **Active Unit Highlight:** The tile occupied by the character currently holding the turn pulses with a neon green glow.
- **Character Placement:** 3D CharacterPawn components rendered at entity grid positions.
- **Path Display:** Movement path highlighted with connected cyan tiles.
- **Grid Lines:** Subtle neon-styled borders on each tile.

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[ui_iso_board]]`
- **Component:** `IsoBoardGrid.vue`
- **Props:** `grid` (width, height, cells[][]), `entities[]`, `highlights` (optional)

## EXPECTATION
- Board renders as a visually correct isometric grid.
- Obstacle tiles are visually distinct from walkable tiles.
- Character pawns appear at their correct grid positions.
- Board handles grid sizes from 5×5 to 15×15.
