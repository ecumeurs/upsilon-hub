# Issue: Character Data Transfer for Battle Engine

**ID:** `20260423_character_data_transfer`
**Ref:** `ISS-076`
**Date:** 2026-04-23
**Severity:** High
**Status:** Open
**Component:** `battleui`, `upsilonapi`, `upsilonbattle`
**Affects:** Battle arena communication, entity initialization

---

## Summary

Define the communication schema for transferring character data from Laravel to the Go battle engine, ensuring that all necessary information (stats, skills, equipped items) is provided in a single structured format. This enables the Go engine to properly initialize entities with both V2 skill system and V2 item system.

---

## Technical Description

### Background

Character data is sent from Laravel to the Go engine via `ArenaStartRequest`, but the current `Entity` struct only contains basic stats (HP, Attack, Defense, Move) and skills. Items and equipment information are not included, making it impossible to use purchased items in battle.

### The Problem Scenario

1. **Skills loaded**: Only character skills are sent to Go engine
2. **No item data**: Purchased armor, weapons, movement items are not transferred
3. **Incomplete entity initialization**: Items in inventory but not applied to battle entity
4. **Missing slots**: Skill slot information is not sent to Go engine

### Character Transfer Architecture

**Enhanced Entity Request:**
```go
type CharacterData struct {
    // Existing fields
    ID       string   // Character ID
    PlayerID string   // Owner
    Team     int      // Team assignment
    Name     string   // Character name
    
    // V2 Stats
    HP       int
    MaxHP    int
    Attack   int
    Defense  int
    Move     int
    MaxMove  int
    
    // V2 Skill System (NEW)
    EquippedSkills []uuid.UUID    // Skills selected for battle
    SkillSlots     int              // Total available slots (1 + level/10)
    
    // V2 Item System (NEW)
    EquippedItems []EquippedItem  // Items equipped on character
}

type EquippedItem struct {
    ItemID    uuid.UUID     // Shop item reference
    Name      string        // Item name
    Type      ItemType      // Armor, Weapon, Movement
    Properties map[string]property.Property  // Item stat bonuses
}
```

### Enhanced Arena Start Request

```go
type ArenaStartRequest struct {
    MatchID     string
    CallbackURL string
    Players     []CharacterData  // Enhanced with skills and items
}
```

**Entity Initialization with Items:**
```go
func NewEntityFromCharacter(charData CharacterData) Entity {
    entity := NewEntity()
    
    // Set basic stats
    entity.ID = ParseUUID(charData.ID)
    entity.Name = charData.Name
    entity.RepertPropertyValue(property.HP, charData.HP)
    entity.RepertPropertyValue(property.MaxHP, charData.MaxHP)
    entity.RepertPropertyValue(property.Attack, charData.Attack)
    entity.RepertPropertyValue(property.Defense, charData.Defense)
    entity.RepertPropertyValue(property.Move, charData.Move)
    entity.RepertPropertyCMaxValue(property.MaxMove, charData.MaxMove)
    
    // Register equipped skills (V2 skill system)
    for _, skillID := range charData.EquippedSkills {
        skill := GetSkillByID(skillID)  // Load skill template
        entity.RegisterSkill(skill)
    }
    
    // Apply equipped items (V2 item system)
    for _, item := range charData.EquippedItems {
        ApplyItemProperties(entity, item)
    }
    
    return entity
}

func ApplyItemProperties(entity *Entity, item EquippedItem) {
    switch item.Type {
    case Armor:
        rating := entity.GetProperty(property.ArmorRating)
        bonus := item.Properties[property.ArmorRating]
        rating.SetValue(rating.I() + bonus.I())
        
    case Weapon:
        damage := entity.GetProperty(property.Attack)
        bonus := item.Properties[property.WeaponRating]
        damage.SetValue(damage.I() + bonus.I())
        // Apply other weapon properties (range, crit, etc.)
        
    case Movement:
        move := entity.GetProperty(property.Move)
        bonus := item.Properties[property.Move]
        move.SetValue(move.I() + bonus.I())
    }
}
```

### Communication Schema Updates

**Request Structure:**
```go
// Match start request includes enhanced character data
POST /arena/start
{
    "match_id": "...",
    "callback_url": "...",
    "players": [
        {
            "id": "char-1",
            "player_id": "player-1",
            "team": 1,
            "name": "Warrior",
            "hp": 35, "max_hp": 35,
            "attack": 12, "defense": 6, "move": 3, "max_move": 3,
            
            // V2 Skill System
            "equipped_skills": ["skill-1-id", "skill-2-id"],
            "skill_slots": 2,
            
            // V2 Item System
            "equipped_items": [
                {"item_id": "armor-1", "name": "Basic Armor", "type": "armor"},
                {"item_id": "weapon-1", "name": "Basic Sword", "type": "weapon"}
            ]
        }
    ]
}
```

### Where This Pattern Exists Today

- `upsilonapi/api/input.go:20` - `type ArenaActionRequest struct`, `type Entity struct`
- `upsilonbattle/battlearena/entity/entity.go:20` - `Skills map[uuid.UUID]skill.Skill` field exists
- `docs/rule_character_skill_slots.atom.md` - Skill slot system (NEW)
- `docs/mech_skill_equipment_system.atom.md` - Skill equipment system (NEW)
- `ISS-073` (Roguelike Skill System) - Combines skills + slots (NEW)
- `ISS-074` (Simple Shop) - Shop items catalog (NEW)
- `ISS-075` (Player Inventory) - Inventory tables (NEW)

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | High |
| Impact if triggered | High |
| Detectability | High |
| Current mitigant | Skills exist but items not included in entity |

---

## Recommended Fix

**Short term:** Update `ArenaStartRequest` in `upsilonapi/api/input.go` to include `EquippedSkills`, `SkillSlots`, and `EquippedItems`. Update Go engine's entity initialization to apply item bonuses.

**Medium term:** Create `EquippedItem` struct in Go codebase. Implement `ApplyItemProperties()` function. Add item property bonuses to damage calculations.

**Long term:** Implement full item system with item effects (on-hit triggers, passive bonuses). Add item durability and repair systems.

---

## User Experience Flow

**Player Preparation:**
1. Select 2 skills from inventory to equip (2/3 slots used)
2. Equip "Basic Armor" from purchased items
3. Equip "Basic Sword" from purchased items
4. Start match

**Battle Init:**
1. Server receives character data with skills and items
2. Go engine creates entity with:
   - 2 registered skills
   - +5 armor rating from armor item
   - +5 weapon rating from weapon item
3. Entity is ready for combat with full V2 systems

**Combat:**
1. Skills appear as action buttons (only equipped skills)
2. Attacks include weapon bonuses
3. All purchased items contribute to character power
