---
id: ui_holo_obstacle
status: STABLE
type: UI
layer: ARCHITECTURE
version: 1.0
tags: [ui, combat, board, obstacle, holographic]
parents:
  - [[ui_iso_board]]
dependents: []
human_name: Holographic Obstacle UI
priority: 3
---

# Holographic Obstacle UI

## INTENT
To provide a 3D holographic visual representation for blocked/obstacle tiles on the tactical board, consistent with the "Neon in the Dust" aesthetic and the character pawn design.

## THE RULE / LOGIC
- **Geometry:** A precision-engineered 3D isometric box perfectly matched to the 64x32 tactical grid tile diamond.
- **Coordinate System:** Component $(0,0)$ is anchored to the tile center. Base vertices are located at $(\pm32, 0)$ and $(0, 16)$ in screen space.
- **Side Facets:** 
    - Width: 32px (horizontal projection).
    - Skew: $\pm 26.565^\circ$ (matching 2:1 slope).
    - Alignment: Perfectly anchored to base vertices $(\pm32, 0)$ and $(0, 16)$.
- **Top Facet:** Square of side 45.25px rotated $60^\circ(X)/45^\circ(Z)$ to form a 64x32 diamond. Centered exactly at $Y = -H_{monolith}$ to seal the apex.
- **Visuals:** Industrial scanlines, randomized rural/rust gradients, and glitch overlays are maintained on the perfect solid volume.
- **Grounding:** Base glow matches tile dimensions (64x32) at $Y=+16$ relative to center.

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[ui_holo_obstacle]]`
- **Component:** `HoloObstacle.vue`
- **Props:** `seed` (used for height/timing variance).

## EXPECTATION
- Obstacle tiles are capped by a visible 3D monolith.
- The visual style (scanlines, glitches) is consistent with the `CharacterPawn`.
- The amber color is clearly distinguishable from the team colors used for pawns.
