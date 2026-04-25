
### **Phase 1: Environmental Lighting & Atmosphere**

**Objective:** Transform the lighting from a "sunny day" to an "underground neon arena."
**Target File:** `ThreeGrid.vue`

**Instructions for Flash:**
> "Refactor the lighting and environment in the provided `ThreeGrid.vue` file using TresJS.
> 1. **Ambient Light:** Change the `<TresAmbientLight>` intensity to `0.15` and set its color to a dark, cool hex (e.g., `#0a0a1a`).
> 2. **Fog:** Add a `<TresFogExp2>` to the `<TresCanvas>` with a color matching the clear-color (`#05050a`) and a density of `0.025` to create depth drop-off.
> 3. **Directional to Spot/Point:** Remove the existing `<TresDirectionalLight>`. Replace it with two `<TresSpotLight>` components positioned at opposite diagonal corners of the grid (`gridCenter.x/z`). Make one cyan (`#00f2ff`) and the other magenta (`#ff00ff`), with high intensity (e.g., `5.0`), casting shadows, and angled towards the grid center."

---

### **Phase 2: Architectural Refactoring & Overlays**

**Objective:** Break the monolith into maintainable components and restore the missing 2D UI elements.
**Target Files:** `ThreeGrid.vue` (parent), `Pawn3D.vue` (new), `Obstacle3D.vue` (new), `Tile3D.vue` (new).

**Instructions for Flash:**
> "Refactor the TresJS scene in `ThreeGrid.vue` by extracting the repeating meshes into separate Vue components. 
> 1. **Create `Tile3D.vue`:** Extract the `<TresMesh>` for the floor tiles. It should accept a `tile` object prop.
> 2. **Create `Obstacle3D.vue`:** Extract the obstacle mesh.
> 3. **Create `Pawn3D.vue`:** Extract the `<TresConeGeometry>` pawn mesh. It must accept an `entity` object prop.
> 4. **Add UI Overlays:** Inside the new `Pawn3D.vue`, import the `<Html>` component from `@tresjs/cientos`. Position it slightly above the cone (`Y + 1.5`). Inside the `<Html>` tags, render a standard HTML `<div>` containing the pawn's name and a basic CSS health bar based on `entity.hp`. Ensure the `<Html>` component uses `transform` so it scales with the camera."

---

### **Phase 3: Materials & Textures**

**Objective:** Add the "battered and worn" aesthetic to the clean geometry.
**Target Files:** `Tile3D.vue`, `Obstacle3D.vue`

**Instructions for Flash:**
> "Upgrade the materials in the TresJS components to simulate battered, physical tech. 
> 1. **Obstacles (`Obstacle3D.vue`):** Update the `<TresMeshStandardMaterial>`. Since we don't have external texture files yet, simulate a dark, rough metallic surface. Set `color="#1a1a1c"`, `roughness={0.9}`, and `metalness={0.8}`. 
> 2. **Highlights:** For the highlighted cells (movement/attack), change the transparent `<TresBoxGeometry>` to a `<TresPlaneGeometry>` rotated `-Math.PI / 2` to lie flat on the tile. Use a `<TresMeshBasicMaterial>` with `wireframe={true}`, `transparent={true}`, and an emissive neon color (cyan or magenta) to look like a projected grid instead of a solid block."

---

### **Phase 4: Holograms & Post-Processing (The Glitch Aspect)**

**Objective:** Recreate the signature cyberpunk holographic scanlines and add a global neon bloom.
**Target Files:** `Pawn3D.vue`, `ThreeGrid.vue`

**Instructions for Flash:**
> "Implement cyberpunk visual effects using TresJS and Three.js.
> 1. **Hologram Shader (`Pawn3D.vue`):** Replace the `<TresMeshStandardMaterial>` on the pawn with a `<TresShaderMaterial>`. Write a simple vertex and fragment shader. The fragment shader should output the entity's team color, but multiply the alpha by `sin(vUv.y * 50.0 + uTime * 5.0)` to create scrolling horizontal scanlines. Enable `transparent={true}` and `additiveBlending`. Pass a `uTime` uniform that updates via `useRenderLoop()`.
> 2. **Neon Bloom (`ThreeGrid.vue`):** Add post-processing to the main canvas. Import `EffectComposer`, `RenderPass`, and `UnrealBloomPass` from `three/examples/jsm/postprocessing`. Wrap the scene to apply a subtle bloom effect (threshold: 0.2, strength: 1.5, radius: 0.4) so the emissive materials and hologram shaders physically glow against the dark fog."