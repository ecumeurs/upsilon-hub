---
id: rule_gdpr_compliance
human_name: GDPR Compliance Rules
type: RULE
version: 1.0
status: STABLE
priority: CORE
tags: [gdpr, privacy]
parents: [req_security]
dependents: []
---

# GDPR Compliance Rules

## INTENT
Ensures personal data protection through secure deletion (soft delete) and anonymization of sensitive information.

## THE RULE / LOGIC
- **Account Deletion (Right to be Forgotten):** 
  - Deletion of an account MUST be a **soft delete**.
  - The record is marked as deleted but remains in the database for audit/integrity until a purge cycle.
- **Anonymization:**
  - Upon soft deletion or upon request for anonymization, sensitive data fields in `entity_player` MUST be overwritten with non-identifiable placeholders.
  - Sensitive Fields: `full_address`, `birth_date`.
  - Placeholder: `ANONYMIZED`.
- **Right to Portability:**
  - Authenticated users MUST have the ability to download a machine-readable dump of all their personal data stored in the system (e.g., JSON format).
  - Scope: All fields in `entity_player`, win/loss records, and character rosters.

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[rule_gdpr_compliance]]`
- **Test Names:** `TestPlayerSoftDelete`, `TestPlayerDataAnonymization`

## EXPECTATION (For Testing)
- Request account delete -> `deleted_at` timestamp set -> User cannot login.
- Audit `entity_player` table -> `full_address` and `birth_date` show "ANONYMIZED" for deleted user.
