---
id: ui_theme
status: STABLE
human_name: UI Theme Specification
dependents: []
type: UI
layer: ARCHITECTURE
version: 1.0
priority: 5
tags: [ui, styling, theme]
parents:
  - [[req_ui_look_and_feel]]
---

# New Atom

## INTENT
To provide a centralized specification for colors, typography, and styling tokens used across the application.

## THE RULE / LOGIC
### Color Palette
- **Neon Accents**: Cyan (`#00f2ff`), Magenta (`#ff00ff`), Lime (`#39ff13`).
- **Gritty Base**: Deep Void (`#0a0a0b`), Gunmetal (`#1a1a1e`), Oxidized Iron (`#3d2b1f`), Worn Steel (`#4a4a4f`).

### Readability Standards
- **Rule**: Text elements must never use the following base colors on dark backgrounds: `void`, `gunmetal`, `rust`, or `steel`.
- **Requirement**: Preferred text colors for legibility are `cyan`, `magenta`, `lime`, and `white`.

### Typography
- **Headings**: `Orbitron` (Variable: `--font-scifi`)
- **Body**: `Inter` (Variable: `--font-sans`)
- **Technical/Logs**: `JetBrains Mono` (Variable: `--font-mono`)
- **Standard Styling**: Titles must use `uppercase` and `tracking-[0.3em]` to evoke high-tech readouts.

### Component Styling (Neon in the Dust)
- **Panels / Containers**: 
    - Background: `bg-upsilon-gunmetal/30` or `/90` with `backdrop-blur-md/xl`.
    - Border: 1px solid `upsilon-cyan/30` or `upsilon-magenta/30`.
- **Corner Accents**: Panels should include 2px solid corner accents (e.g. `border-t-2 border-l-2`) to reinforce a "HUD" geometry.
- **Interactive States**:
    - **Hover**: Increase border-opacity and add a subtle glow (`shadow-glow-cyan/magenta`).
    - **Active**: Use pulsing indicators (`animate-pulse`) for strategic links.

### Utility Classes (Tailwind)
- Colors prefixed with `upsilon-*`.
- Background components using `hero-bg`, `panel-texture`.
- Text/Border effects: `shadow-neon`, `glow-cyan`.

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[ui_theme]]`
- **Tailwind Config:** [tailwind.config.js](file:///workspace/battleui/tailwind.config.js)
- **CSS:** [app.css](file:///workspace/battleui/resources/css/app.css)

## EXPECTATION
- The primary font must be Orbitron for all display/heading elements.
- The Tailwind config must reflect the specified color palette exactly.
- All colors must have appropriate contrast ratios for readability against dark backgrounds.
