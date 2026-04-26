# 3D Board Display Summary - BattleUI

This document provides a comprehensive summary of the 3D board display implemented in `battleui`, specifically focusing on the `ThreeGrid.vue` component and its associated ATD requirements.

## 1. Technical Layout (Code Implementation)

The 3D board is implemented using **Three.js** via the **TresJS** ecosystem (`@tresjs/core`, `@tresjs/cientos`).

### Core Component: `ThreeGrid.vue`
- **Coordinate System:**
  - **X:** Grid column (`grid.x`).
  - **Y:** Terrain elevation (`cell.height * TILE_HEIGHT`).
  - **Z:** Grid row (`grid.y`).
- **Constants:**
  - `TILE_SIZE = 1.0`
  - `TILE_HEIGHT = 0.25`
- **Camera & Controls:**
  - Uses `TresPerspectiveCamera` (FOV: 45) focused on the grid center.
  - `OrbitControls` for manual inspection, with damping and optional auto-rotation.
- **Scene Composition:**
  - **Terrain:** Iterates through `grid.cells` to render `Tile3D` components.
  - **Obstacles:** Renders `Obstacle3D` components for cells marked as obstacles.
  - **Pawns:** Renders `Pawn3D` for entities, positioned based on their logical grid coordinates and the surface height.
  - **Highlights:** Renders `GridHighlight` for movement and attack ranges.
  - **Facing Indicators:** Renders `FacingIndicator3D` to show character orientation.
- **Lighting & Environment:**
  - `TresAmbientLight`: Base illumination.
  - `TresPointLight`: Warm overhead key light (`#ffecb3`).
  - `TresSpotLight`: Neon accent lights (Cyan `#00f2ff` and Magenta `#ff00ff`) when effects are enabled.
  - `TresFogExp2`: Density-based fog for depth.
  - `PostProcess`: Custom shader-based post-processing.

---

## 2. ATD Requirements (Look and Feel)

### [[req_ui_look_and_feel]] - UI Look and Feel Aesthetic
> **INTENT:** To define the core visual identity and aesthetic philosophy of the Upsilon Battle project.
> 
> **LOGIC:**
> - Aesthetic: "Neon in the Dust" (Sci-fi Post-Apocalyptic).
> - Key Contrast: High-tech vibrancy vs. Gritty industrial decay.
> - UI Directives: 
>   * Use sharp, geometric shapes for tech elements.
>   * Apply texture overlays (dust, noise, rust) to backgrounds.
>   * Glow effects must be used sparingly for primary feedback.
>   * Motion should be linear and 'robotic'.
> - We favor the usage of sci-fi / computer-like terminology (Link terminated, Connection lost, etc.) along with a bit of mad max-like terminology (Scavenged, jury-rigged, etc.) to describe the game state and events.
> 
> **EXPECTATION:**
> - All UI elements must strictly follow the "Neon in the Dust" aesthetic.
> - High-contrast neon elements must be paired with low-contrast industrial textures.
> - The interface must feel alive through subtle kinetic feedback.

### [[ui_theme]] - UI Theme Specification
> **INTENT:** To provide a centralized specification for colors, typography, and styling tokens used across the application.
> 
> **EXPECTATION:**
> - The primary font must be Orbitron for all display/heading elements.
> - The Tailwind config must reflect the specified color palette exactly.
> - All colors must have appropriate contrast ratios for readability against dark backgrounds.

### [[ui_iso_board]] - 3D Hexagonal Board Grid UI (Three.js)
> **INTENT:** The real-time 3D board renderer displaying a dynamic NxM tile grid with terrain elevation, obstacles, character pawns, movement range highlights, and attack target highlights — built on Three.js via `@tresjs/core`.
> 
> **LOGIC:**
> - **Rendering:** Three.js real-time 3D scene via `<TresCanvas>` from `@tresjs/core`.
> - **Digital Hex Terrain:** The ground is rendered as a subdivided hexagonal grid. Each logical square tile is visually represented by a cluster of approximately 16 hexagonal columns (4x4 subdivision).
> - **Seamless Meshing:** Hexagonal columns align perfectly across logical tile boundaries by using a global hex coordinate system.
> - **Cell elevation:** Each hex column rises from $Y=0$ to $Y = height \times TILE\_HEIGHT$.
> - **Tile Types:**
>     - Normal: `TresCylinderGeometry` (6 segments), dark gunmetal (`#2a2a2a`).
>     - Obstacle: taller hexagonal columns in rusty brown (`#3d2b1f`).
> - **Pawns & Highlights:** Positioned relative to the logical grid centers $(gx, gy)$.
> 
> **EXPECTATION:**
> - Board renders as a visually correct 3D scene with terrain elevation visible.
> - Elevated tiles appear higher on the Y-axis than flat tiles; obstacles extrude above their tile surface.
> - Obstacle tiles are visually distinct from walkable tiles.
> - Character pawns appear at correct 3D positions (atop their tile's surface height).
> - Active-turn pawn glows brighter than idle pawns.
> - Board handles grid sizes from 5×5 to 15×15.

### [[ui_character_pawn]] - Character Pawn UI
> **INTENT:** The holographic 3D character pawn rendered on the isometric board — a faceted rotating cone body with a directional sphere head, floating name label, and HP bar.
> 
> **LOGIC:**
> - **Body:** A solid 3D hexagonal pyramid composed of 6 faceted planes that meet at a shared top vertex.
> - **Rotation Logic:** 
>   - **Inactive:** Static at an idle angle (30deg) to provide a consistent silhouette.
>   - **Active:** Rotates continuously on the Y-axis to indicate it is the acting character.
> - **Head:** A sphere with a slight indent indicating character facing direction.
> - **Holographic Effect:** Semi-transparent gradient facets with scanline overlay and intermittent static/glitch animation.
> - **Color Coding:** Blue (current player), Green (ally), Red (enemy player 1), Purple (enemy player 2). Distinct shade variations via HSL shift.
> 
> **EXPECTATION:**
> - Pawn renders as a visible holographic figure on the board.
> - Team colors are clearly distinguishable.
> - Active character has enhanced visual feedback.
> - Holographic glitch effect fires intermittently.

### [[ui_selection_highlight]] - Tactical Selection Highlight
> **INTENT:** To provide clear, high-fidelity visual feedback for tile-based tactical actions (movement, attack targeting) using a dynamic pulsing effect.
> 
> **LOGIC:**
> - **Geometry:** A circular disc (clipping a square plane via shader `discard`).
> - **Pulsing Animation:**
>   - **Move:** Static or slow pulse (`0.005` amplitude).
>   - **Attack:** High-energy pulse (`0.02` amplitude) to indicate threat.
> - **Color Coding:**
>   - **Cyan (#00f2ff):** Valid movement range.
>   - **Magenta (#ff00ff):** Valid attack target selection.
> - **Edge Glow:** A smoothstep-based border gradient (`dist 0.45 to 0.5`) to create a soft "neon" ring.

### [[ui_holo_obstacle]] - Holographic Obstacle UI
> **INTENT:** To provide a 3D holographic visual representation for blocked/obstacle tiles on the tactical board, consistent with the "Neon in the Dust" aesthetic and the character pawn design.
> 
> **LOGIC:**
> - **Geometry:** A precision-engineered 3D isometric box perfectly matched to the 64x32 tactical grid tile diamond.
> - **Visuals:** Industrial scanlines, randomized rural/rust gradients, and glitch overlays are maintained on the perfect solid volume.
> - **Grounding:** Base glow matches tile dimensions (64x32).
> 
> **EXPECTATION:**
> - Obstacle tiles are capped by a visible 3D monolith.
> - The visual style (scanlines, glitches) is consistent with the `CharacterPawn`.
> - The amber color is clearly distinguishable from the team colors used for pawns.

