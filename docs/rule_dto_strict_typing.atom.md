---
id: rule_dto_strict_typing
status: DRAFT
priority: 1
human_name: Strict DTO Typing
tags: [api, architecture, type-safety]
version: 1.0
parents:
  - [[contract_upsilon_contract]]
dependents:
  - [[mechanic_mec_skill_payload_resolution]]
type: RULE
layer: ARCHITECTURE
---

# New Atom

## INTENT
To ensure API contracts are predictable and strongly typed by forbidding the use of 'any' or 'interface{}' in DTOs.

## THE RULE / LOGIC
- DTOs (Data Transfer Objects) must never use `any` or `interface{}` fields for input or output.
- Every field must have a concrete type that describes its structure.
- In cases where external systems provide inconsistent JSON (e.g., empty arrays instead of empty objects), a custom type with a strict Unmarshaler must be used to normalize the data into a concrete Go structure.

## TECHNICAL INTERFACE
- **Code Tag:** `@spec-link [[rule_dto_strict_typing]]`

## EXPECTATION
Every DTO must be validated against a formal schema (JSON Schema or Go Struct tags) and fail-fast if an unexpected type is provided.
