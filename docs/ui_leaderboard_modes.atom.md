---
id: ui_leaderboard_modes
status: DRAFT
human_name: Leaderboard Categories (Modes)
type: UI
layer: ARCHITECTURE
priority: 5
version: 1.0
dependents: []
parents:
  - [[ui_leaderboard]]
---

# New Atom

## INTENT
Define the categorical split of leaderboards by battle mode.

## THE RULE / LOGIC
- The leaderboard system must support 4 distinct categories:
  1. **1v1 PvP**: Single player vs Single player.
  2. **2v2 PvP**: Team vs Team.
  3. **1v1 PvE**: Single player vs AI.
  4. **2v2 PvE**: Team vs AI.
- Users must be able to switch between these categories via a toggle/tab interface.
- Statistics for one mode do not bleed into others.

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[ui_leaderboard_modes]]`

## EXPECTATION
