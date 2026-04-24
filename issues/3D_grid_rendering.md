# 3D Grid Rendering Implementation Plan

**Issue:** ISS-3D-RENDERING
**Status:** DESIGN
**Priority:** HIGH
**Assigned:** —
**Created:** 2026-04-24

---

## Executive Summary

This document outlines the implementation of 3D grid rendering for UpsilonBattle's Vue.js frontend (battleui) using TresJS. The CLI (upsiloncli) will remain functionally intact as the grid structure is backward compatible.

**Timeline:** 6-8 weeks estimated
**Scope:** 5 Phases
**Risk:** MODERATE (requires API changes but maintains CLI compatibility)

---

## Technical Considerations & Corrections

### 1. Coordinate System Mapping (Critical Fix)

**Issue Identified:** Three.js uses Y-axis for vertical height, Z-axis for depth. Current pseudo-code incorrectly reuses `position.y` for both.

**Three.js Coordinate System:**
```
     Y (up)
     |
     |
     +---- X (right)
    /
   Z (depth)
```

**API/Grid Coordinate System:**
```
(grid.x, grid.y) → Three.js position:
  X = grid.x * TILE_SIZE
  Y = grid.height * TILE_HEIGHT  ← Use grid.height, NOT grid.y
  Z = grid.y * TILE_SIZE
```

**Requirement:** Entity position from API must be mapped as:
- `props.entity.position.x` → Three.js X
- `grid.cell.height` (from terrain) → Three.js Y
- `props.entity.position.y` → Three.js Z

### 2. InstancedMesh Reactivity Overhead

**Issue:** Vue's computed properties with large matrices will recalculate on every grid mutation, potentially causing full InstancedMesh buffer rebuilds.

**Solution:** Separate static terrain from dynamic state updates.

```vue
<script setup>
// Static terrain matrices - computed ONCE on mount or when grid changes shape
const terrainInstances = ref(null)
const staticMatrices = computed(() => {
  const matrices = []
  const { width, height, cells } = props.grid

  for (let x = 0; x < width; x++) {
    for (let y = 0; y < height; y++) {
      const cell = cells[x]?.[y]
      if (!cell) continue
      const matrix = new THREE.Matrix4()
      matrix.setPosition(x * TILE_SIZE, cell.height * TILE_HEIGHT, y * TILE_SIZE)
      matrices.push(matrix)
    }
  }
  return matrices
})

// Dynamic highlights - updated separately via setMatrixAt
const highlightIndices = ref([])

function setMatrixAt(instanceIndex, matrix) {
  if (terrainInstances.value) {
    // Direct matrix update without full rebuild
    const dummy = new THREE.Object3D()
    dummy.updateMatrix()
    terrainInstances.value.setMatrixAt(instanceIndex, dummy.matrix)
  }
}

// Only recompute when grid dimensions change (not cell content changes)
watch(() => [props.grid.width, props.grid.height], () => {
  // Full rebuild only on resize
})
</script>
```

### 3. Post-Processing Performance Warning

**Issue:** UnrealBloomPass + custom CRT shader + FilmPass can drop FPS below 30 on low-end integrated graphics.

**Solution:** Implement graphics quality toggle in settings.

```vue
<!-- battleui/resources/js/Components/Arena/SettingsPanel.vue -->
<template>
  <div class="quality-toggle">
    <label>
      <input type="checkbox" v-model="highQuality" />
      High Quality (Neon Effects)
    </label>
    <label>
      <input type="checkbox" v-model="performanceMode" />
      Performance Mode
    </label>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const highQuality = ref(localStorage.getItem('graphics_quality') !== 'performance')
const performanceMode = ref(localStorage.getItem('graphics_quality') === 'performance')

watch(highQuality, (val) => {
  localStorage.setItem('graphics_quality', val ? 'high' : 'performance')
  // Trigger PostProcessing component to enable/disable passes
})
</script>
```

**Conditional Rendering in PostProcessing.vue:**
```vue
<template>
  <EffectComposer v-if="!performanceMode">
    <UnrealBloomPass />
    <ShaderPass />  <!-- CRT shader -->
    <FilmPass />
  </EffectComposer>
  <!-- Direct render when performance mode -->
  <div v-else class="performance-warning">Effects disabled for performance</div>
</template>
```

### 4. Data Payload Size

**Expected Grid Size Range:** 10×10 to 30×30 (max 40×40)

**Payload Analysis:**
| Grid Size | Cell Count | HeightCell Size | Uncompressed |
|-----------|-------------|-----------------|--------------|
| 10×10 | 100 | ~40 bytes | ~4 KB |
| 20×20 | 400 | ~40 bytes | ~16 KB |
| 30×30 | 900 | ~40 bytes | ~36 KB |
| 40×40 (max) | 1,600 | ~40 bytes | ~64 KB |

**Conclusion:** Payload sizes remain well within acceptable limits even without compression. Compression is still recommended for WebSocket bandwidth optimization.

**File:** `battleui/config/cors.php` or middleware

```php
// Add response compression
'gzip' => \Illuminate\Http\Middleware\CompressResponseMiddleware::class,
```

**Verification:** Test payload sizes with curl:
```bash
curl -H "Accept-Encoding: gzip" http://localhost:8000/api/v1/game/ID | wc -c
# Should be ~10-15% of uncompressed size
```

---

## Current Architecture Analysis

### Data Flow

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│ upsilonmapmaker│───▶│  upsilonmapdata│───▶│   upsilonapi   │
│ (Generator)    │    │   (Data Model) │    │ (Game Engine)  │
└─────────────────┘    └─────────────────┘    └────────┬────────┘
                                                       │
                                                       ▼
                                              ┌──────────────────────┐
                                              ┌──▶│  Laravel API       │
                                              │    │ (Gateway/Webhook) │
┌───────────────────┐    ┌───────────────┐   │    └────────┬───────────┘
│  upsiloncli      │    │   battleui     │───┘             │
│  (CLI Frontend) │◀───│  (Web Frontend)│◀────────────────┘
│ (2D ASCII only) │    │  (Vue 3 + SVG)│    HTTP/WS
└───────────────────┘    └───────────────┘
```

### Current Grid Data Model

**upsilonmapdata/grid/grid.go** (Internal 3D):
```go
type Grid struct {
    Width  int
    Length int
    Height int
    Cells map[position.Position]*cell.Cell
}

type Cell struct {
    Type     CellType      // Ground, Water, Obstacle, Dirt
    Position position.Position // (X, Y, Z) - TRUE 3D!
    EntityIDs []uuid.UUID
    EffectIDs []uuid.UUID
}
```

**upsilonapi/api/output.go** (Current API Output - FLATTENED):
```go
type Cell struct {
    EntityID string `json:"entity_id"`
    Obstacle bool   `json:"obstacle"`
    // MISSING: Z-height, CellType (Ground/Water/Dirt)
}

type Grid struct {
    Width  int      `json:"width"`
    Height int      `json:"height"`
    Cells  [][]Cell `json:"cells"` // 2D array ONLY
}
```

**Current Limitation:** The API flattens the 3D grid structure to 2D, losing:
- Terrain elevation (Z-height per cell)
- Cell type distinction (Ground vs Water vs Dirt)
- Underground information

**CLI Compatibility:** The CLI receives the same flattened 2D structure and renders ASCII. This is BACKWARD COMPATIBLE with changes since it only uses `obstacle` boolean and `entity_id` for display.

---

## Phase 1: Core Architecture & Uplink

### 1.1 Install TresJS Dependencies

**File:** `battleui/package.json`

```diff
  "dependencies": {
    "@inertiajs/vue3": "^2.0.0",
    "@tresjs/core": "^4.0.0",
    "@tresjs/cientos": "^4.0.0",
    "three": "^0.170.0",
    "uuid": "^13.0.0"
+   "postprocessing": "^6.35.0"
  }
```

### 1.2 Extend API Data Payload

**File:** `upsilonapi/api/output.go`

```go
// NEW: TerrainType for detailed cell classification
type TerrainType string

const (
    TerrainGround TerrainType = "ground"
    TerrainWater TerrainType = "water"
    TerrainDirt  TerrainType = "dirt"
)

// NEW: HeightCell extends Cell with Z information
type HeightCell struct {
    EntityID     string      `json:"entity_id"`
    Obstacle     bool        `json:"obstacle"`
    Height       int         `json:"height"`       // NEW: Z-coordinate
    TerrainType  TerrainType `json:"terrain_type"` // NEW: Ground/Water/Dirt
}

// UPDATED: Grid3D extends Grid with height data
type Grid3D struct {
    Width  int          `json:"width"`
    Height int          `json:"height"`
    MaxHeight int       `json:"max_height"`  // NEW: For bounding box
    Cells  [][]HeightCell `json:"cells"`
}
```

**File:** `upsilonapi/api/output.go` - NewBoardState function

```go
func NewBoardState3D(matchID uuid.UUID, g *grid.Grid, ...) BoardState {
    // ... existing code ...

    bs.Grid = Grid3D{
        Width:     g.Width,
        Height:    g.Length,
        MaxHeight: g.Height, // NEW
        Cells:     make([][]HeightCell, g.Width),
    }

    for x := 0; x < g.Width; x++ {
        bs.Grid.Cells[x] = make([]HeightCell, g.Length)
        for y := 0; y < g.Length; y++ {
            z := g.TopMostCellAt(x, y)
            cl, ok := g.CellAt(position.New(x, y, z))
            if ok {
                // NEW: Map cell types
                terrain := TerrainGround
                if cl.Type == cell.Water {
                    terrain = TerrainWater
                } else if cl.Type == cell.Dirt {
                    terrain = TerrainDirt
                }

                bs.Grid.Cells[x][y] = HeightCell{
                    EntityID:    charID,
                    Obstacle:    cl.Type == cell.Obstacle,
                    Height:      z,  // NEW: Pass Z coordinate
                    TerrainType: terrain, // NEW: Cell type
                }
            }
        }
    }
    // ... rest of function ...
}
```

### 1.3 Update Communication Protocol

**File:** `communication.md`

```diff
### Grid
-**Grid: A 2D array of cells; for our purpose as in this implementation, height will be fixed at 1 for every cell giving us a flat map.**
+**Grid: A 3D grid representing terrain topography. Supports elevation and terrain type (ground/water/dirt).**

```json
{
  "width": 10,
  "height": 10,
+ "max_height": 15,
  "cells": [
    [
      {
        "entity_id": "uuid...",
        "obstacle": false,
+       "height": 3,
+       "terrain_type": "ground"
      }
    ]
  ]
}
```

**Backward Compatibility Note:** CLI clients will ignore new fields (`height`, `terrain_type`, `max_height`) and continue to work with `obstacle` and `entity_id` only.

---

## Phase 2: The Wasteland Grid (Rendering)

### 2.1 TresJS Grid Component

**New File:** `battleui/resources/js/Components/Arena/ThreeGrid.vue`

```vue
<script setup>
import { TresCanvas, useRenderLoop } from '@tresjs/core'
import { OrbitControls } from '@tresjs/cientos'
import { ref, computed, onMounted, watch } from 'vue'
import { InstancedMesh } from 'three'

const props = defineProps({
  grid: { type: Object, required: true },
  entities: { type: Array, default: () => [] },
  currentEntityId: { type: String, default: '' },
  teamColors: { type: Object, default: () => ({}) },
  highlightedCells: { type: Array, default: () => [] },
  // Performance mode to disable expensive effects
  performanceMode: { type: Boolean, default: false },
})

const TILE_SIZE = 1.0
const TILE_HEIGHT = 0.2

// InstancedMesh refs for direct Three.js access
const groundMesh = ref(null)
const waterMesh = ref(null)
const obstacleMesh = ref(null)

const { onLoop } = useRenderLoop()

// OPTIMIZED: Static matrices computed ONCE when grid dimensions change
// Cell content changes (entity_id, highlights) are handled via direct matrix updates
const instanceMatrices = computed(() => {
  const matrices = []
  const { width, height, cells } = props.grid
  let instanceCount = 0

  for (let x = 0; x < width; x++) {
    for (let y = 0; y < height; y++) {
      const cell = cells[x]?.[y]
      if (!cell) continue

      const matrix = new THREE.Matrix4()
      // Three.js: X=grid.x, Y=cell.height, Z=grid.y
      matrix.setPosition(x * TILE_SIZE, cell.height * TILE_HEIGHT, y * TILE_SIZE)
      matrices.push(matrix)
      instanceCount++
    }
  }
  return matrices
})

// Update specific instance matrix without full recalculation
function updateInstanceMatrix(meshRef, index, x, y, z) {
  if (!meshRef.value) return
  const matrix = new THREE.Matrix4()
  matrix.setPosition(x * TILE_SIZE, z * TILE_HEIGHT, y * TILE_SIZE)
  meshRef.value.setMatrixAt(index, matrix)
  meshRef.value.instanceMatrix.needsUpdate = true
}

// Only rebuild instances on grid dimension changes (shape change)
watch(() => [props.grid.width, props.grid.height, props.grid.max_height], () => {
  // Triggered only when map grows/shrinks, not cell content changes
}, { immediate: true })

onMounted(() => {
  // Set up instanced mesh buffers
})
</script>

<template>
  <TresCanvas shadows>
    <!-- OrbitControls with restricted pitch -->
    <OrbitControls
      :max-polar-angle="Math.PI / 2.2"
      :min-distance="5"
      :max-distance="30"
      :enable-damping="true"
    />

    <!-- Ambient: low and gritty -->
    <TresAmbientLight :intensity="0.3" />

    <!-- Directional: sunlight/moonlight -->
    <TresDirectionalLight
      :position="[10, 15, 10]"
      :intensity="0.8"
      cast-shadow
    />

    <!-- PointLights: neon glows for active cells -->
    <TresPointLight
      v-for="(cell, i) in highlightedCells"
      :key="'highlight-' + i"
      :position="[cell.x, cell.height * TILE_HEIGHT + 1, cell.y]"
      :color="cell.type === 'attack' ? '#ff00ff' : '#00f2ff'"
      :intensity="0.8"
      :distance="5"
    />

    <!-- Instanced Ground Tiles -->
    <TresInstancedMesh
      ref="groundMesh"
      :args="[TILE_SIZE, TILE_HEIGHT, TILE_SIZE]"
      :count="instanceMatrices.instanceCount"
      color="#2a2a2a"
    >
      <TresMeshStandardMaterial :roughness="0.8" :metalness="0.3" />
    </TresInstancedMesh>

    <!-- Water surface -->
    <TresInstancedMesh
      ref="waterMesh"
      :args="[TILE_SIZE, TILE_HEIGHT, TILE_SIZE]"
      color="#0a1520"
      transparent
      :opacity="0.8"
    >
      <TresMeshStandardMaterial :roughness="0.1" :metalness="0.1" />
    </TresInstancedMesh>

    <!-- Obstacles (stacked cubes based on height) -->
    <TresInstancedMesh
      ref="obstacleMesh"
      :args="[TILE_SIZE, TILE_HEIGHT * 3, TILE_SIZE]"
      color="#3d2b1f"
    >
      <TresMeshStandardMaterial :roughness="0.9" :metalness="0.6" />
    </TresInstancedMesh>
  </TresCanvas>
</template>
```

### 2.2 Character Pawns in 3D

**Modify File:** `battleui/resources/js/Components/Arena/CharacterPawn.vue`

**IMPORTANT: Coordinate System Correction**

Three.js uses **Y-axis for vertical height** and **Z-axis for depth**.
The API entity position is 2D (x, y) where y represents grid depth.
The grid now includes `height` field for terrain elevation.

```vue
<script setup>
import { TresMesh } from '@tresjs/core'

const props = defineProps({
  entity: { type: Object, required: true },
  teamColor: { type: String, default: '#39ff13' },
  isActive: { type: Boolean, default: false },
  // NEW: Grid cell height for terrain-aware positioning
  terrainHeight: { type: Number, default: 0 },
})

const TILE_SIZE = 1.0
const TILE_HEIGHT = 0.2

// FIXED: Three.js coordinate mapping
// X-axis: entity.x (grid column)
// Y-axis: terrainHeight (elevation from grid.height)
// Z-axis: entity.y (grid row depth)
const position = computed(() => [
  props.entity.position.x * TILE_SIZE,
  props.terrainHeight * TILE_HEIGHT + TILE_HEIGHT / 2,
  props.entity.position.y * TILE_SIZE
])

// LookAt target for rotation
const rotation = computed(() => {
  // Smooth rotation toward target if attacking
  return [0, 0, 0]
})
</script>

<template>
  <TresMesh
    :position="position"
    :rotation="rotation"
    cast-shadow
  >
    <!-- Cone shape for robotic units -->
    <TresConeGeometry :args="[0.3, 0.8, 4]" />
    <TresMeshStandardMaterial
      :color="teamColor"
      :emissive="teamColor"
      :emissive-intensity="isActive ? 0.5 : 0"
      :roughness="0.3"
      :metalness="0.8"
    />
  </TresMesh>
</template>
```

---

## Phase 3: Camera & Tactical Control

### 3.1 Orbit Controls Integration

**Configuration in ThreeGrid.vue:**
- **Max Polar Angle:** `Math.PI / 2.2` (prevents going under map)
- **Min Distance:** 5 units
- **Max Distance:** 30 units
- **Enable Damping:** Smooth camera movement
- **Auto Rotate:** OFF (tactical precision required)

### 3.2 Pawn Orientation

**Function:** Add to battleui game logic

```javascript
//战斗ui/resources/js/services/tactical.js
export function calculatePawnRotation(from, to) {
  if (!to) return [0, 0, 0]
  
  const angle = Math.atan2(to.x - from.x, to.y - from.y)
  return [0, angle, 0] // Rotate around Y axis
}

// For attack animations
export function snapToTarget(attacker, target) {
  const angle = Math.atan2(
    target.position.y - attacker.position.y,
    target.position.x - attacker.position.x
  )
  return [0, angle, 0] // Abrupt, robotic motion
}
```

---

## Phase 4: Applying "Neon" (Post-Processing)

### 4.1 Post-Processing Setup

**New File:** `battleui/resources/js/Components/Arena/PostProcessing.vue`

```vue
<script setup>
import { useRenderLoop } from '@tresjs/core'
import { EffectComposer } from 'postprocessing'
import { UnrealBloomPass } from 'postprocessing'
import { ShaderPass } from 'postprocessing'
import { FilmPass } from 'postprocessing'

const composer = ref(null)
const bloom = ref(null)
const crt = ref(null)

// Custom CRT/Noise shader
const crtShader = {
  uniforms: {
    tDiffuse: { value: null },
    time: { value: 0 },
    scanlineIntensity: { value: 0.1 },
    noiseIntensity: { value: 0.05 },
    curvature: { value: 0.05 }
  },
  vertexShader: `
    varying vec2 vUv;
    void main() {
      vUv = uv;
      gl_Position = projectionMatrix * modelViewMatrix * vec4(position, 1.0);
    }
  `,
  fragmentShader: `
    uniform sampler2D tDiffuse;
    uniform float time;
    uniform float scanlineIntensity;
    uniform float noiseIntensity;
    uniform float curvature;

    varying vec2 vUv;

    float random(vec2 st) {
      return fract(sin(dot(st.xy, vec2(12.9898, 78.233))) * 43758.5453);
    }

    void main() {
      vec2 uv = vUv;

      // Screen curvature
      uv = (uv - 0.5) * (1.0 - curvature) + 0.5;

      // Scanlines
      float scanline = sin(uv.y * 800.0 + time * 5.0) * scanlineIntensity;

      // Noise/film grain
      float noise = random(uv + time) * noiseIntensity;

      vec4 color = texture2D(tDiffuse, uv);
      gl_FragColor = color - scanline - noise;
    }
  `
}

onLoop(({ elapsed }) => {
  if (crt.value) {
    crt.value.material.uniforms.time.value = elapsed
  }
})
</script>

<template>
  <EffectComposer ref="composer">
    <UnrealBloomPass
      ref="bloom"
      :strength="1.5"
      :radius="0.4"
      :threshold="0.1"
    />
    <ShaderPass
      ref="crt"
      :vertex-shader="crtShader.vertexShader"
      :fragment-shader="crtShader.fragmentShader"
      :uniforms="crtShader.uniforms"
    />
  </EffectComposer>
</template>
```

### 4.2 Neon Configuration

**Bloom Settings:**
- **Strength:** 1.5 (intense glow)
- **Radius:** 0.4 (spread)
- **Threshold:** 0.1 (only bright elements glow)

**What Glows:**
- Active unit outline (`emissive` material property)
- Move range tiles (cyan)
- Attack range tiles (magenta)
- Action feedback indicators

### 4.3 Performance Toggle (Conditional Rendering)

**File:** `battleui/resources/js/Components/Arena/SettingsPanel.vue`

```vue
<script setup>
import { ref, watch } from 'vue'
import { useSettingsStore } from '@/stores/settings'

const settings = useSettingsStore()
const quality = ref(settings.graphicsQuality === 'high')

watch(quality, (val) => {
  settings.setGraphicsQuality(val ? 'high' : 'performance')
})
</script>

<template>
  <div class="settings-toggle">
    <button :class="{ active: quality }" @click="quality = true">
      High Quality
    </button>
    <button :class="{ active: !quality }" @click="quality = false">
      Performance
    </button>
  </div>
</template>
```

**PostProcessing.vue with conditional rendering:**

```vue
<template>
  <EffectComposer v-if="!performanceMode">
    <UnrealBloomPass ref="bloom" />
    <ShaderPass ref="crt" />
    <FilmPass />
  </EffectComposer>
</template>

<script setup>
import { computed } from 'vue'
import { useSettingsStore } from '@/stores/settings'

const settings = useSettingsStore()
const performanceMode = computed(() => settings.graphicsQuality === 'performance')
</script>
```

---

## Phase 5: The HUD Overlay

### 5.1 Layering Strategy

**Architecture:** DOM overlay on top of WebGL canvas

```html
<div class="tactical-container">
  <!-- Layer 0: 3D Canvas (z-index: 0) -->
  <ThreeGrid />

  <!-- Layer 1: HUD Overlay (z-index: 100) -->
  <div class="hud-overlay">
    <CombatHeader />
    <TeamRosterPanel side="left" />
    <ActionPanel />
    <TeamRosterPanel side="right" />
    <InitiativeTimeline />
    <TacticalActionReport />
    <ConfirmModal />
  </div>
</div>
```

**CSS:**
```css
.tactical-container {
  position: relative;
  width: 100%;
  height: 100vh;
}

.hud-overlay {
  position: absolute;
  inset: 0;
  pointer-events: none; /* Allow clicks to pass through to canvas */
}

/* Re-enable pointer events for interactive HUD elements */
.hud-overlay > * {
  pointer-events: auto;
}
```

### 5.2 Keep Existing Vue Components

**No Changes Needed To:**
- `CombatHeader.vue` - Already DOM-based
- `TeamRosterPanel.vue` - Already DOM-based
- `ActionPanel.vue` - Already DOM-based
- `InitiativeTimeline.vue` - Already DOM-based
- `TacticalActionReport.vue` - Already DOM-based
- `ConfirmModal.vue` - Already DOM-based

**Changes Needed To:**
- `IsoBoardGrid.vue` → Replace/Deprecate (new `ThreeGrid.vue`)
- `HoloObstacle.vue` → Deprecate (obstacles now 3D instanced)
- `CharacterPawn.vue` → Update to TresJS rendering

---

## CLI Compatibility (No Caves)

### Constraint
We **WILL NOT** allow caves/underground navigation in this iteration.

### Rationale
1. CLI's ASCII map rendering (`upsiloncli/internal/display/printer.go:Board()`)
   - Only uses: `cell.obstacle`, `cell.entity_id`
   - Does NOT use Z-height for navigation

2. CLI will remain functionally correct as long as:
   - Ground cells are walkable (obstacle=false)
   - Obstacle cells are blocked (obstacle=true)
   - Entity positions are correct (x, y)

3. The API extension adds:
   - `height` field (ignored by CLI)
   - `terrain_type` field (ignored by CLI)

### Verification
```go
// upsiloncli/internal/display/printer.go
func (p *Printer) Board(bs *dto.BoardState, currentUserID string, players []dto.Player) {
    // ... existing code ...
    for y := 0; y < bs.Grid.Height; y++ {
        for x := 0; x < bs.Grid.Width; x++ {
            cell := bs.Grid.Cells[y][x]
            // CLI only reads these:
            if cell.EntityID != "" {
                // render entity symbol
            } else if cell.Obstacle {
                // render '#'
            } else {
                // render '.'
            }
            // New fields (height, terrain_type) are ignored
        }
    }
}
```

---

## Implementation Checklist

### Phase 1: Core Architecture
- [ ] Install TresJS and Three.js dependencies
- [ ] Create `HeightCell` struct in upsilonapi
- [ ] Create `Grid3D` struct in upsilonapi
- [ ] Update `NewBoardState` to populate height/terrain
- [ ] Update `communication.md` with new payload format
- [ ] Update upsilonapi ATD atoms for new fields
- [ ] **[NEW]** Enable GZIP compression in Laravel API gateway
- [ ] **[NEW]** Verify payload size reduction with compression enabled

### Phase 2: Wasteland Grid
- [ ] Create `ThreeGrid.vue` component
- [ ] Implement InstancedMesh for ground tiles
- [ ] Implement InstancedMesh for water surfaces
- [ ] Implement InstancedMesh for obstacles
- [ ] Add low ambient lighting
- [ ] Add directional sun/moonlight
- [ ] Add PointLights for highlighted cells
- [ ] Test with upsilonmapmaker Hill/River outputs
- [ ] **[NEW]** Implement static/dynamic separation for instance matrices
- [ ] **[NEW]** Add `updateInstanceMatrix()` for targeted updates
- [ ] **[NEW]** Watch grid dimensions (not content) for rebuild triggers

### Phase 3: Camera & Controls
- [ ] Integrate OrbitControls with pitch restriction
- [ ] Set min/max zoom distances
- [ ] Implement pawn LookAt rotation
- [ ] Add smooth camera damping
- [ ] Test camera angles for tactical visibility

### Phase 4: Post-Processing
- [ ] Install postprocessing package
- [ ] Configure UnrealBloomPass
- [ ] Implement custom CRT/Noise shader
- [ ] Add film grain effect
- [ ] Tune neon glow intensity
- [ ] Profile performance impact
- [ ] **[NEW]** Create `SettingsPanel.vue` with quality toggle
- [ ] **[NEW]** Implement conditional rendering for Performance mode
- [ ] **[NEW]** Test FPS with/without post-processing

### Phase 5: HUD Integration
- [ ] Create tactical container with z-index layering
- [ ] Integrate ThreeGrid.vue into BattleArena.vue
- [ ] Remove IsoBoardGrid.vue (deprecate)
- [ ] Update ActionPanel for 3D coordinates (if needed)
- [ ] Test HUD overlay positioning and clicks
- [ ] Ensure pointer-events pass-through works correctly

### Testing
- [ ] CLI continues to work with new API (ignore new fields)
- [ ] Web UI renders 3D grid correctly
- [ ] Hill/River maps display with proper elevation
- [ ] Water tiles are semi-transparent
- [ ] Obstacles have height visualization
- [ ] Neon effects glow properly
- [ ] Camera controls feel responsive
- [ ] Performance: >60 FPS on standard maps

---

## Risk Mitigation

| Risk | Mitigation |
|-------|------------|
| **Performance** | InstancedMesh reduces draw calls from ~1600 to 1 per material type (40×40 grid). Separated static/dynamic updates to avoid full rebuilds. |
| **Post-Processing** | Performance toggle allows users on low-end GPUs to disable bloom/CRT effects entirely. |
| **API Breaking** | CLI ignores new fields (`height`, `terrain_type`); maintains backward compatibility |
| **Payload Size** | Max grid 40×40 = 64KB uncompressed. GZIP compression still recommended for WebSocket bandwidth. |
| **Browser Support** | TresJS/Three.js supports modern browsers; fallback to 2D can be added if needed. |
| **ATD Drift** | Update all relevant atoms with new field specs. |
| **Coordinate Systems** | Fixed: Three.js uses Y for height, Z for depth. API (x,y) maps to (X, Z) with grid.height → Y. |

---

## References

- **External Party Input:** 2026-04-24
- **upsilonmapmaker:** `/upsilonmapmaker/gridgenerator/gridgenerator.go`
- **upsilonmapdata:** `/upsilonmapdata/grid/grid.go`
- **upsilonapi:** `/upsilonapi/api/output.go`
- **battleui:** `/battleui/resources/js/Components/Arena/IsoBoardGrid.vue`
- **upsiloncli:** `/upsiloncli/internal/display/printer.go`
- **communication:** `/communication.md`

---

## Appendix: Terrain Type Color Mapping

| TerrainType | Three.js Color | Description |
|-------------|-----------------|-------------|
| `ground` | `#2a2a2a` | Rusted metal/earth |
| `water` | `#0a1520` | Dark toxic liquid |
| `dirt` | `#1a1a1a` | Underground fill |

---

*End of Document*
