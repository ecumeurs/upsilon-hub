---
id: us_win_progression
human_name: Post-Win Progression Story
type: USER_STORY
layer: BUSINESS
version: 1.0
status: STABLE
priority: 5
tags: []
parents: []
dependents: []
---
# Post-Win Progression Story

## INTENT
As a player, after winning a match I am shown a progression screen where I allocate Character Points, with results reflected immediately on my Dashboard.

## THE RULE / LOGIC
After a match victory the player allocates progression to their roster. Acceptance criteria:
- A progression screen appears only after a match victory.
- The player allocates earned Character Points (CP) to upgradable stats via the CP point-buy system defined in [[rule_progression]] (per-attribute CP costs and the global spend cap; no fixed "1 point per win" allocation, and no movement-every-5-wins gate).
- Applying an upgrade immediately reflects on the player's Dashboard character stats.

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[us_win_progression]]`
