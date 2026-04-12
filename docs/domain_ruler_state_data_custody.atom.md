---
id: domain_ruler_state_data_custody
human_name: Data Custody Split
type: DOMAIN
layer: CUSTOMER
version: 1.0
status: STABLE
priority: 5
tags: []
parents:
  - [[domain_ruler_state]]
dependents: []
---
# Data Custody Split

## INTENT
Independently holds the Grid Data and the Entities map.

## THE RULE / LOGIC
The Ruler independently holds the Grid Data and the Entities map. Clients must independently request and synchronize with this state payload.

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[domain_ruler_state_data_custody]]`
