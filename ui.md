# Upsilon UI Design System

## Theme Overview
**Style**: Sci-fi Post-Apocalyptic / Cyberpunk Industrial
**Core Concept**: "Neon in the Dust" - High-tech elements clashing with gritty, rusty, and worn-out industrial textures.

---

## Color Palette

### Neon Accents (The Tech)
- **Neon Cyan**: `#00f2ff` (Primary Glow, UI Borders)
- **Neon Magenta**: `#ff00ff` (Secondary Glow, Alerts)
- **Neon Lime**: `#39ff13` (Success, Energy, Progress)

### Gritty Base (The World)
- **Deep Void**: `#0a0a0b` (Main Background)
- **Gunmetal**: `#1a1a1e` (Component Backgrounds)
- **Oxidized Iron**: `#3d2b1f` (Rust shadows, subtle textures)
- **Worn Steel**: `#4a4a4f` (Secondary text, inactive elements)

### Utility Colors
- **Warning Orange**: `#f59e0b` (Hazard signs, alerts)
- **Critical Red**: `#ef4444` (Damage, danger)

---

## Typography
- **Headings**: `Orbitron` (Sci-fi, Industrial)
- **Body**: `Inter` or `Geist` (Modern, readable)
- **Monospace**: `JetBrains Mono` or `Roboto Mono` (Data, logs, technical readouts)

---

## Tailwind Configuration Extensions

```javascript
// tailwind.config.js
module.exports = {
  theme: {
    extend: {
      colors: {
        upsilon: {
          cyan: '#00f2ff',
          magenta: '#ff00ff',
          lime: '#39ff13',
          void: '#0a0a0b',
          gunmetal: '#1a1a1e',
          rust: '#3d2b1f',
          steel: '#4a4a4f',
        }
      },
      fontFamily: {
        scifi: ['Orbitron', 'sans-serif'],
        mono: ['JetBrains Mono', 'monospace'],
      },
      boxShadow: {
        'glow-cyan': '0 0 10px rgba(0, 242, 255, 0.5), 0 0 20px rgba(0, 242, 255, 0.2)',
        'glow-magenta': '0 0 10px rgba(255, 0, 255, 0.5), 0 0 20px rgba(255, 0, 255, 0.2)',
      },
      backgroundImage: {
        'rust-texture': "url('/assets/textures/rust-overlay.png')", // To be generated/referenced
      }
    }
  }
}
```

---

## UI Components Design

### Buttons
- **Primary (Neon)**: Sharp edges, neon border-glow, uppercase text.
- **Secondary (Gritty)**: Dark background, subtle rust-colored hover state, "warning tape" pattern on edges.

### Cards / Panels
- Semi-transparent dark backgrounds.
- "Corner-brackets" design instead of full borders.
- Scanline or noise overlays for a retro-tech look.
