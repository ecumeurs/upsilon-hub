---
id: us_queue_selection
human_name: Queue Selection Story
type: USER_STORY
layer: BUSINESS
version: 1.0
status: STABLE
priority: 5
tags: []
parents: []
dependents: []
---
# Queue Selection Story

## INTENT
As a logged-in player on the Dashboard, I can pick a game mode from clearly distinct queue buttons and see my record.

## THE RULE / LOGIC
A logged-in player selects a game type directly from the Dashboard. Acceptance criteria:
- The Dashboard displays exactly 4 distinct queue buttons: `1v1 PVE`, `1v1 PVP`, `2v2 PVE`, `2v2 PVP`.
- Clicking a PVE button starts the game immediately (no waiting room).
- Clicking a PVP button navigates to the Waiting Room.
- The Dashboard also displays the player's current Win/Loss record and Win/Loss ratio.

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[us_queue_selection]]`
