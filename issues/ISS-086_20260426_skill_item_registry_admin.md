# Issue: Skill and Item Registries with Admin CRUD

**ID:** `20260426_skill_item_registry_admin`
**Ref:** `ISS-086`
**Date:** 2026-04-26
**Severity:** High
**Status:** Open
**Component:** `battleui`, `upsilonapi`
**Affects:** Admin dashboard, item creation, skill management, shop catalog

---

## Summary

Transition from a fixed, hardcoded list of items and skills to a dynamic, database-backed registry. Implement a central repository for skills (templates) and items (patterns), including an admin section for full CRUD management. These registries will allow designers to create and modify game content without code changes. Skills will serve as building blocks for item effects (e.g., weapon skills or consumable effects).

---

## Technical Description

### Background
Currently, items are defined in a fixed list (seeded via `ShopItemsSeeder.php`). Skills are referenced but lack a central, manageable registry. The shop catalog is essentially a static set of entries in the `shop_items` table.

### The Problem Scenario
1. **Hardcoded Content**: Adding a new skill or item requires a code change and re-seeding.
2. **Lack of Reuse**: No central "template" for common items or skills.
3. **No Admin UI**: No way for non-developers to balance or add content.
4. **Disconnected Systems**: Skills and items are managed separately, making it hard to create items with complex skill-based effects.

### Proposed Architecture

#### 1. Skill Registry (`skill_templates`)
- Stores "blueprints" for skills.
- Fields: `id`, `name`, `type`, `logic_type` (e.g., projectile, area, buff), `properties` (JSON), `metadata` (JSON).
- Properties include: damage, range, cooldown, stamina cost, etc.

#### 2. Item Registry (`item_patterns`)
- Stores "blueprints" for items.
- Fields: `id`, `name`, `type`, `slot`, `properties` (JSON), `effect_skill_id` (FK to `skill_templates`).
- Patterns can be "instantiated" into the shop with specific costs and availability.

#### 3. Admin Dashboard
- **Skill Manager**: List, Create, Edit, Delete skill templates.
- **Item Manager**: List, Create, Edit, Delete item patterns.
- **Form Builder**: Dynamic forms that adapt to the properties required by the skill/item type.

#### 4. Integration
- Items can "carry" a skill template. When an item is equipped, the character gains access to the associated skill.

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | Medium |
| Impact if triggered | High |
| Detectability | High |
| Current mitigant | Seeders provide a fallback, but they are hard to maintain. |

---

## Recommended Fix

**Short term:**
- Design the database schema for `skill_templates` and `item_patterns`.
- Create the migration and models in `battleui`.

**Medium term:**
- Build the Admin CRUD views using Laravel Nova or custom Vue components (Inertia.js).
- Implement dynamic form builders for property editing.

**Long term:**
- Refactor the Shop to use the registries instead of hardcoded seeder data.
- Integrate the Go engine to fetch skill/item data from these registries (or via the API).

---

## References

- [ShopItemsSeeder.php](file:///workspace/battleui/database/seeders/ShopItemsSeeder.php)
- [ISS-073_20260423_roguelike_skill_system.md](file:///workspace/issues/ISS-073_20260423_roguelike_skill_system.md)
- [ISS-074_20260423_comprehensive_item_system.md](file:///workspace/issues/ISS-074_20260423_comprehensive_item_system.md)
- [entity_shop_item.atom.md](file:///workspace/upsilonbattle/docs/entity_shop_item.atom.md)
