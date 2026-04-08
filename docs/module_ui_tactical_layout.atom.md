---
id: module_ui_tactical_layout
status: STABLE
tags: [ui, layout]
parents:
  - [[ui_dashboard]]
dependents: []
human_name: Tactical UI Layout
type: MODULE
layer: ARCHITECTURE
priority: 5
version: 1.0
---

# New Atom

## INTENT
To provide a consistent, tactical user interface framework across all authenticated pages.

## THE RULE / LOGIC
- Composes `TacticalHeader` and `TacticalFooter`.
- Applies the "Neon in the Dust" global styling (backgrounds, fonts, colors).
- Provides a slot for main content.

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[ui_tactical_layout]]`
- **File:** [TacticalLayout.vue](file:///workspace/battleui/resources/js/Layouts/TacticalLayout.vue)

## EXPECTATION
- The layout must wrap all authenticated tactical pages.
- The header must persist across navigation.
- The footer must be visible at all times for system status monitoring.
