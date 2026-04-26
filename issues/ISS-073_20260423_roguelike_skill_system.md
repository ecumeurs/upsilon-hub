# Issue: Roguelike Skill System - Inventory, Slots & Equipment

**ID:** `20260423_roguelike_skill_system`
**Ref:** `ISS-073`
**Date:** 2026-04-23
**Severity:** High
**Status:** Open
**Component:** `battleui`, `upsilonapi`, `upsilonbattle`
**Affects:** `battleui/app/Models/Character.php`, character progression, battle engine integration

---

## Summary

Implement comprehensive roguelike-style skill system with character skill inventory, slot-based equipment, and battle integration. Characters gain skill slots every 10 levels (starting with 1), acquire skills through selection/shop, and equip chosen skills before battle. This transforms the V2 skill system from simple acquisition to strategic inventory management.

---

## Wiring instructions (post-ISS-074, 2026-04-26)

ISS-074 (item system) shipped the structural groundwork for skills inside the same engine bootstrap path. **All scaffolding is in place — this issue is now data-side only**. To complete:

1. **Schema:** Create the `character_skills` join table (`character_id` FK, `skill_id` FK, `equipped_at`, `acquired_at`). Decide whether equip is a flag on this join, or a separate `equipped_skills` mirror table — the engine doesn't care, it just receives the equipped subset.

2. **Laravel resource:** In `app/Http/Resources/API/Upsilon/UpsilonEntityResource.php`, populate the `equipped_skills` array (currently emitted as `[]`) with the equipped skill UUIDs from `character_skills WHERE equipped=true`.

3. **Go bridge:** In `upsilonapi/bridge/bridge.go`, the entity-bootstrap loop already reads `entity.EquippedSkills []string` from the request payload (added by ISS-074, currently always empty). Add the resolution step:
   ```go
   for _, skillID := range entity.EquippedSkills {
       skill := skillregistry.Lookup(uuid.MustParse(skillID))  // or your registry of choice
       e.RegisterSkill(skill)
   }
   ```
   `RegisterSkill` already exists on `Entity` (see `entity.go`).

4. **Drop the placeholder comment:** Once ISS-073 lands, delete the `// reserved for ISS-073` comment in `upsilonapi/api/input.go` next to the `EquippedSkills` field.

5. **Item-carried skills (already wired):** Items with an `Effect` property (a skill ID) already register their skill at equip time via `[[upsilonbattle:mec_item_buff_application]]`. ISS-073 doesn't need to handle weapon-as-skill — it's already done.

**Key references from ISS-074:**
- Reserved field: `upsilonapi/api/input.go` — `Entity.EquippedSkills []string`
- Insertion point: `upsilonapi/bridge/bridge.go` (the entity-bootstrap loop, just after the item-buff loop)
- Atom relationships: `[[upsilonbattle:mec_item_buff_application]]`, `[[upsilonbattle:entity_character_skill_inventory]]`, `[[upsilonbattle:mech_skill_equipment_system]]`

---

## Technical Description

### Background

Current skill system (V2 planning) has skill selection at character creation and reforging mechanics, but lacks:
- **Skill inventory system** - No mechanism to store skills separate from active usage
- **Skill slot system** - No limit on how many skills a character can use in battle
- **Skill equipment** - No preparation phase where players select active skills

Entity struct already supports skills via `Skills map[uuid.UUID]skill.Skill`, but this is unbounded and doesn't distinguish between owned vs equipped skills.

### The Problem Scenario

1. **Character Creation:** Player selects 1 of 3 skills, but no distinction between owning vs using them
2. **Skill Accumulation:** As player levels, they acquire multiple skills but all are immediately available
3. **No Strategic Depth:** Players cannot make equipment decisions before battle
4. **No Skill Inventory:** Skills cannot be stored, swapped, or managed separately

### Roguelike Skill System Architecture

**Three-Layer System:**

1. **Skill Inventory** (Unlimited storage)
   - All acquired skills stored in character record
   - Viewable, sortable, filterable
   - Persistent across matches

2. **Skill Slots** (Progressive cap)
   - Base: 1 slot at level 1
   - Formula: Slots = 1 + (CharacterLevel / 10)
   - Level 1-9: 1 slot
   - Level 10-19: 2 slots
   - Level 20-29: 3 slots
   - Level 30-39: 4 slots
   - Level 40+: 5 slots (soft cap)

3. **Skill Equipment** (Battle preparation)
   - Select skills from inventory to fill available slots
   - Equipment decisions are strategic (offensive loadout vs defensive)
   - Only equipped skills appear in battle action panel
   - Persistent across matches

### Database Schema Updates

**Characters Table:**
```sql
-- Skill slot tracking
ALTER TABLE characters ADD COLUMN skill_slots INTEGER DEFAULT 1;

-- Current equipped skills (JSON array)
ALTER TABLE characters ADD COLUMN equipped_skills JSON;

-- All owned skills (JSON array or pivot table)
-- Option A: JSON column (simpler)
ALTER TABLE characters ADD COLUMN skill_inventory JSON;

-- Option B: Pivot table (more relational)
CREATE TABLE character_skills (
    id UUID PRIMARY KEY,
    character_id UUID REFERENCES characters(id),
    skill_id UUID NOT NULL,
    acquired_at TIMESTAMP DEFAULT NOW(),
    reforge_count INTEGER DEFAULT 0
);
```

### API Endpoints

**Skill Inventory Management:**
- `GET /api/v1/character/{id}/skill-inventory` - View all owned skills
- `GET /api/v1/skills/templates` - Browse available skill templates
- `POST /api/v1/character/{id}/skill-purchase` - Buy skill from shop

**Skill Equipment:**
- `POST /api/v1/character/{id}/equip-skill` - Equip skill from inventory
- `DELETE /api/v1/character/{id}/unequip-skill/{skillId}` - Unequip skill
- `POST /api/v1/character/{id}/set-loadout` - Set full equipment loadout

**Skill Review & Inspection:**
- `GET /api/v1/character/{id}/skill/{skillId}` - View detailed skill information
- `GET /api/v1/skills/compute-grade` - Compute grade from skill properties
- CLI command: `skill inspect <skillId>` - Display skill properties, effects, and statistics

### CLI Integration

**Skill Usage Commands:**
- `skill list` - View all owned skills with equipped status
- `skill equip <skillId>` - Equip a specific skill from inventory
- `skill unequip <skillId>` - Unequip a skill
- `skill inspect <skillId>` - Display detailed skill properties and effects

**Skill Review UI Features:**
- Skill detail modal showing all properties and effects
- Visual grade indicator (I-V)
- Property breakdown (Damage, Range, Cooldown, etc.)
- Effect descriptions with tooltip explanations

**Character Progression:**
- Should appear in the character profile already as a specialized attribute.

### Battle Engine Integration

**Entity Initialization:**
```go
// When creating entity for battle, load equipped skills only
func NewEntityFromCharacter(char Character) Entity {
    entity := NewEntity()

    // Load only equipped skills, not all owned skills
    for _, skillID := range char.EquippedSkills {
        skill := char.GetSkillFromInventory(skillID)
        entity.RegisterSkill(skill)
    }

    return entity
}
```

**UI Integration:**
- Character sheet shows both inventory and equipped skills
- Pre-battle screen allows equipment adjustments
- Battle panel only shows equipped skill buttons
- Visual indicator for available slots vs used slots
- Skill detail modal for reviewing properties and effects
- Skill statistics panel showing usage metrics

**CLI Integration:**
- `skill list`, `skill equip`, `skill unequip`, `skill inspect` commands
- CLI displays skill properties in human-readable format
- Skill effects clearly explained with tooltips
- Statistics tracking accessible via command line

### Where This Pattern Exists Today

- `battleui/app/Models/Character.php` - Character model (needs skill fields)
- `upsilonbattle/battlearena/entity/entity.go:47` - Entity.Skills map exists
- `upsilonbattle/battlearena/entity/entity.go:142` - RegisterSkill() method exists
- `docs/mech_skill_selection_progression.atom.md` - Skill selection (DRAFT)
- `docs/mech_skill_reforging_mechanic.atom.md` - Skill reforging (DRAFT)

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | High |
| Impact if triggered | High |
| Detectability | High |
| Current mitigant | Entity.Skills map exists, needs slot validation |

---

## Recommended Fix

**Short term (Week 1):**
- Create database schema updates for skill_inventory, equipped_skills, skill_slots
- Update Character model with new fields and relationships
- Create basic skill inventory view API

**Medium term (Week 2):**
- Implement skill equipment/unequipment endpoints
- Add slot validation logic
- Create equipment management UI

**Long term (Week 3-4):**
- Integrate equipment loadout into battle entity creation
- Build pre-battle preparation screen
- Add skill statistics tracking (usage, damage, credits per skill)

---

## Integration with V2 Issues

**Related Issues:**
- **ISS-065** (Skill Weight & Grading) - Skill inventory uses grading for pricing
- **ISS-066** (Time-Based Mechanics) - Equipped skills can have channeling
- **ISS-067** (Credit Economy) - Shop purchases add to inventory
- **ISS-071** (Starting Stats & Progression) - Level determines slot count
- **ISS-074** (Comprehensive Item System) - Items use similar equipment system pattern

**New ATD Atoms Created:**
- `rule_character_skill_slots` - Slot progression formula
- `entity_character_skill_inventory` - Inventory management
- `mech_skill_equipment_system` - Equipment operations

---

## User Experience Flow

**Character Creation:**
1. Player rolls for 3 random skills (Grade I-II)
2. Player selects 1 skill
3. Skill added to inventory AND auto-equipped (only slot 1 available)
4. Character created with 1 equipped skill, 0 free slots

**Level Progression (Every 10 Levels):**
1. Character gains new skill slot (1 → 2 → 3, etc.)
2. System offers 3 new skills based on level access
3. Player selects 1 skill
4. Skill added to inventory
5. Player must manually equip new skill if slots full

**Pre-Battle Preparation:**
1. Player views character sheet showing inventory + equipped skills
2. Available slots: "2/3" (2 used, 1 free)
3. Player unequips old skill, equips new skill
4. Loadout saved for battle

**During Battle:**
1. Only equipped skills appear as action buttons
2. Entity has exactly 3 registered skills
3. Strategic decisions made before battle matter

---

## References

- `docs/mech_skill_selection_progression.atom.md` - Skill selection timeline
- `docs/mech_skill_reforging_mechanic.atom.md` - Reforging operations
- `docs/rule_skill_grading_system.atom.md` - Grade-based pricing
- `docs/rule_character_skill_slots.atom.md` - Slot calculation (NEW)
- `docs/entity_character_skill_inventory.atom.md` - Inventory system (NEW)
- `docs/mech_skill_equipment_system.atom.md` - Equipment operations (NEW)
