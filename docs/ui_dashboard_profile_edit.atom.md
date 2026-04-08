---
id: ui_dashboard_profile_edit
status: DRAFT
type: UI
layer: ARCHITECTURE
version: 1.0
priority: 5
tags: [profile, user, gdpr]
dependents: []
human_name: Profile Management UI
parents:
  - [[ui_dashboard]]
---

# New Atom

## INTENT
To allow users to manage their personal profile information and exercise their GDPR rights.

## THE RULE / LOGIC
- Display fields for address and birth date.
- Provide a button to trigger data export (GDPR compliance).
- Provide a button to delete account (with multi-step confirmation).
- Access: Only via authenticated dashboard session.

## TECHNICAL INTERFACE
- **API Endpoints:** `PATCH /v1/profile`, `GET /v1/profile/export`, `DELETE /v1/profile`
- **Code Tag:** `@spec-link [[ui_dashboard_profile_edit]]`
- **Logic Rule:** [rule_gdpr_compliance.atom.md](file:///workspace/docs/rule_gdpr_compliance.atom.md)

## EXPECTATION
- Users must be able to update their address and birth date in a dedicated profile section.
- Users must be able to export all their personal data in compliance with GDPR.
- Users must be able to request account deletion.
- A confirmation step is mandatory for account deletion and sensitive data changes.
