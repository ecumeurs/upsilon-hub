---
id: ui_action_panel
status: STABLE
type: UI
version: 1.0
parents:
  - [[ui_battle_arena]]
  - [[mech_action_economy]]
dependents: []
human_name: Action Panel Component
layer: ARCHITECTURE
priority: 4
tags: [ui, combat, actions, turn]
---

# New Atom

## INTENT
A self-contained, boxed combat action panel that enables or disables player actions depending on whether the authenticated user owns the current turn.

## THE RULE / LOGIC
- **Box Structure:** Outer border container with a status header row and a button row.
- **Header row contains:**
  - Status pill: `YOUR TURN` (cyan, pulsing dot) or `WAITING` (muted, gray dot) based on `isPlayerTurn`.
  - Character context: active entity name, HP, movement remaining.
  - Owner label: nickname of the player whose turn it is, or `⬡ Sending…` while processing.
- **Button row:** MOVE (+20/tile), ATTACK (+100), PASS (+300), separator, FORFEIT.
- **Disabled state:** When `!isPlayerTurn` or `isProcessing`, the button row gets `opacity: 0.28`, `filter: saturate(0.3) brightness(0.7)`, and `pointer-events: none`. A translucent lock overlay appears below the header with the text `Actions locked — awaiting your turn`.
- **Active state:** When `isPlayerTurn`, the outer box glows with a subtle cyan border/shadow.
- **Turn authority is determined at the player level:** `isPlayerTurn = (current_player_id === authenticated_user.id)`. Any character belonging to that player activates the panel.
- **Selected action highlight:** MOVE and ATTACK buttons show a colored selected state when clicked, cleared on board update or cancel.");

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[ui_action_panel]]`
- **File:** `resources/js/Components/Arena/ActionPanel.vue`
- **Props:** `isPlayerTurn`, `isProcessing`, `canMove`, `canAttack`, `moveCostPerTile`, `attackCost`, `passCost`, `selectedAction`, `activeCharacter`, `activePlayerName`
- **Emits:** `action` with payload string `'move' | 'attack' | 'pass' | 'forfeit'`

## EXPECTATION
- Panel renders correctly in both YOUR TURN and WAITING states.
- In WAITING state: button row is visually grayed out (desaturated + dimmed), overlay text appears, status pill shows muted WAITING.
- In YOUR TURN state: box border glows cyan, status pill has pulsing cyan dot, buttons are fully interactive.
- Character context (name, HP, MOV) is visible in the header for the currently active entity.
- Active player nickname is shown on the right of the header when it is not the user's turn.
- FORFEIT button requires only `isPlayerTurn`; MOVE additionally requires `canMove`; ATTACK requires `canAttack`.
- All actions are completely unreachable (pointer-events: none) when locked.
