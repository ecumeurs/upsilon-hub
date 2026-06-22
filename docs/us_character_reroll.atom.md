---
id: us_character_reroll
human_name: Character Reroll Story
type: USER_STORY
layer: BUSINESS
version: 1.0
status: STABLE
priority: 5
tags: []
parents:
  - [[shared:req_tech_debt_backlog]]
dependents:
  - [[upsilonapi:api_profile_character]]
---
# Character Reroll Story

## INTENT
As a new player creating my characters, I can reroll my starting roster up to three times to get a set I'm happy with.

## THE RULE / LOGIC
A new player's character-creation flow allows limited rerolling of the starting roster. Acceptance criteria:
- A clear, noticeable "Reroll" button is present on the character creation screen.
- Clicking "Reroll" discards the current character set and regenerates three new characters.
- A visible counter shows the number of rerolls remaining (e.g. "Rerolls remaining: 2").
- After three successful rerolls the "Reroll" button is disabled to prevent further use.

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[us_character_reroll]]`
