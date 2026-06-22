---
id: us_leaderboard_view
human_name: Leaderboard View Story
type: USER_STORY
layer: BUSINESS
version: 1.0
status: STABLE
priority: 5
tags: []
parents: []
dependents:
  - [[upsilonapi:rule_leaderboard_cycle]]
  - [[upsilonapi:rule_leaderboard_score_calculation]]
---
# Leaderboard View Story

## INTENT
As an authenticated player, I view a global leaderboard ranking all players by wins, with a link reachable from the Dashboard.

## THE RULE / LOGIC
An authenticated player views a global leaderboard ranking all players. Acceptance criteria:
- A leaderboard link is clearly present on the Dashboard.
- The leaderboard is accessible only to authenticated users (valid bearer token required).
- The leaderboard is sorted in descending order of total Wins, with ties broken by Win/Loss ratio.

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[us_leaderboard_view]]`
