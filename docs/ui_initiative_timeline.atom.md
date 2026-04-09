---
id: ui_initiative_timeline
status: DRAFT
parents:
  - [[ui_battle_arena]]
  - [[mech_action_economy]]
human_name: Initiative Timeline UI
layer: ARCHITECTURE
priority: 5
tags: [ui, combat, initiative, timeline, turn-order]
dependents: []
type: UI
version: 1.0
---

# New Atom

## INTENT
A horizontal timeline at the bottom of the arena showing the turn order queue, each character's next expected activation tick, and projected action cost shadows.

## THE RULE / LOGIC
- **Layout:** Horizontal bar spanning the full width of the center panel.
- **Tokens:** Character icons placed along the timeline based on their `delay` value (lower delay = further left = sooner to act).
- **Active Character:** Delay=0 character glows brightly at the leftmost position.
- **Action Shadows:** When hovering an action, projected future position shown as a ghost token (e.g., "Move 3 tiles → next turn at tick +60").
- **Color Coding:** Tokens match team colors (blue/green for allies, red/purple for enemies).
- **Scale:** Timeline normalized to the maximum delay value in the turn queue.

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[ui_initiative_timeline]]`
- **Component:** `InitiativeTimeline.vue`
- **Props:** `turns[]` (sorted by delay), `teamColors` (map of player_id to color)

## EXPECTATION
- Turn order is visually clear from left (next) to right (furthest).
- Active character token is prominently highlighted.
- Tokens use correct team colors.
