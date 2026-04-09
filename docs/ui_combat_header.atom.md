---
id: ui_combat_header
status: STABLE
human_name: Combat Header UI
version: 1.0
priority: 5
parents:
  - [[ui_battle_arena]]
dependents: []
type: UI
layer: ARCHITECTURE
tags: [ui, combat, header, hp-bar, timer]
---

# New Atom

## INTENT
The fighting-game-style combat header displaying team HP bars, remaining character counts, match duration timer, and shot clock with neon color transitions.

## THE RULE / LOGIC
- **Left Side:** Current player team's remaining character count + total HP bar (fills left-to-right).
- **Right Side:** Adversary team's remaining character count + total HP bar (fills right-to-left).
- **HP Bar Colors:** Transitions from neon green (>60%) → orange (30-60%) → red (<30%) based on remaining HP percentage.
- **Center Top:** Match duration timer (ticks every second, format MM:SS).
- **Center Bottom:** Shot clock (30s per turn countdown). Text turns orange at ≤10s, red at ≤5s. Always neon-styled.
- **Visual Style:** Evokes fighting game (Street Fighter/Mortal Kombat) HP bar aesthetics with neon glow.

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[ui_combat_header]]`
- **Component:** `CombatHeader.vue`

## EXPECTATION
- HP bars visually shrink as team HP decreases.
- Shot clock numbers change color at 10s and 5s thresholds.
- Match timer ticks up every second.
