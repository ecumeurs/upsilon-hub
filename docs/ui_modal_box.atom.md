---
id: ui_modal_box
status: DRAFT
layer: ARCHITECTURE
tags: [ui, modal, design-system]
parents:
  - [[requirement_customer_user_account]]
dependents: []
human_name: UI Modal Box Specification
type: UI
priority: 5
version: 1.0
---

# New Atom

## INTENT
To provide a consistent visual and behavioral specification for interactive modal dialogs.

## THE RULE / LOGIC
- **Visuals:** Dark semi-transparent backdrop (`upsilon-void/80`). 
- **Layout:** Centered panel with `upsilon-gunmetal/40` background and `upsilon-magenta/30` borders.
- **Header:** Scifi-style title in `Orbitron` font, uppercase.
- **Decorations:** Magenta glow effects (`shadow-glow-magenta`) and top-left/bottom-right corner accents.
- **Interactions:** Close button in the top right; supports ESC to close.

## TECHNICAL INTERFACE
- **Component:** `ModalBox.vue`
- **Code Tag:** `@spec-link [[ui_modal_box]]`
- **Related Spec:** [[ui_theme]]

## EXPECTATION
