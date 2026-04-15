---
id: ui_tactical_action_report
human_name: "Tactical Action Report"
type: UI
layer: ARCHITECTURE
version: 1.0
status: STABLE
priority: 2
tags: [ui, combat, feedback]
parents:
  - [[requirement_customer_action_reporting]]
dependents: []
---

# Tactical Action Report

## INTENT
Provide immediate, human-readable visual feedback of tactical actions within the Battle Arena to eliminate the need for manual state diffing.

## THE RULE / LOGIC
- **Trigger:** Must activate whenever the `BoardState.action` field is populated in a state update.
- **Persistence:** Should remain visible for 3 seconds before fading out.
- **Content Mapping:**
  - `attack`: Display "DAMAGED", the numerical value of damage dealt, and the HP transition (PrevHP -> NewHP).
  - `move`: Display "REPOSITIONED" and the number of tiles crossed (path length).
  - `pass`: Display "STANCE RESET" to indicate turn end.
- **Styling:** Must use high-contrast holographic aesthetics (Cyan #00f2ff for general, Red #ff2020 for damage).

## TECHNICAL INTERFACE (The Bridge)
- **Vue Component:** `TacticalActionReport.vue`
- **Code Tag:** `@spec-link [[ui_tactical_action_report]]`
- **Props:**
  - `action`: Object containing `{type, damage, prev_hp, new_hp, path}`.
  - `show`: Boolean controlling visibility transition.

## EXPECTATION (For Testing)
- **UI Transition:** Component must use a fade transition when toggling `show`.
- **Data Integrity:** Action details must exactly match the `ActionFeedback` DTO received from the Go engine.
- **Concurrent Actions:** If a new action is received while one is already displaying, the 3-second timer must reset and the content refresh.
