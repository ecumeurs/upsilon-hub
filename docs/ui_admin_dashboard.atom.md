---
id: ui_admin_dashboard
human_name: "Admin Dashboard Page UI"
type: MODULE
layer: ARCHITECTURE
version: 1.0
status: STABLE
priority: 5
tags: [admin, ui]
parents:
  - [[uc_admin_login]]
dependents: []
---

# Admin Dashboard Page UI

## INTENT
To serve as the primary landing hub for authorized Administrators to manage system maintenance.

## THE RULE / LOGIC
1. Accessible only after successful **Admin Login** [[uc_admin_login]].
2. Provides navigation to **User Management** (UC-5) and **History Management** (UC-6).
3. Enforces [[rule_admin_access_restriction]] at the presentation layer.

## TECHNICAL INTERFACE (The Bridge)
- **Frontend Page:** `Admin/Dashboard.vue`
- **Code Tag:** `@spec-link [[ui_admin_dashboard]]`
- **Related Issue:** `#admin-dashboard`

## EXPECTATION (For Testing)
- Only users with `Admin` role can access the dashboard.
- Redirection to Login occurs if no active Admin session exists.
