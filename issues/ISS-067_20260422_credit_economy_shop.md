# Issue: Credit Economy & Shop System

**ID:** `20260422_credit_economy_shop`
**Ref:** `ISS-067`
**Date:** 2026-04-22
**Severity:** High
**Status:** Open
**Component:** `upsilonapi/api`, `battleui`
**Affects:** `upsilonbattle/battlearena`, database schema

---

## Summary

Implement comprehensive credit economy with multiple earning mechanisms (damage, healing, support, status effects) and shop system for purchasing skills and equipment. Credits are earned through combat performance and spent on character progression.

---

## Technical Description

### Background

No economy system exists. Players have no way to earn or spend currency. Skill selection and equipment acquisition lack progression mechanics.

### The Problem Scenario

1. **Combat Ends**: Player deals damage, heals allies, mitigates damage
2. **No Reward**: No credit earning system exists
3. **No Spending**: No shop or purchase mechanisms
4. **No Progression**: Players can't acquire new skills or equipment

### Credit Earning System

**Base Rule:** 1 HP damage = 1 coin (healing also earns credits)

**Support Credits:**
- Damage mitigation: 1 HP mitigated = 1 coin
- Shield caster earns credits when shield blocks damage
- Effect must track caster for proper credit assignment

**Status Effect Credits (Option A - Flat Rate):**
- Poison/Stun/Buff: SkillWeight/10 credits per application
- 100 SW poison skill = 10 credits per poison application
- Applies at moment of effect application, not per-turn

**Effect Caster Tracking:**
```go
type Effect struct {
    Properties []property.Property
    Name       string
    CasterID   uuid.UUID  // Track creator for credit assignment
    OriginTime time.Time  // When effect was applied
}
```
- Effects remember caster until effect ends
- Credits go to original caster even if they die later
- Critical for shield/healing credit assignment

### Shop System

**Pricing Formula:** Credit Cost = Total Positive SW × 2
- Grade I Basic Attack (100 SW) = 200 credits
- Grade V Meteor Swarm (800 SW) = 1600 credits

**Shop Inventory:**
- Skills: Available based on player level/grade
- Equipment: Armor, Utility, Weapon categories
- Credit spending: Deduct from character balance

### Where This Pattern Exists Today

- `upsilonapi/api/input.go` - Player structure
- `battleui/app/Models/User.php` - User model
- `docs/rule_progression.atom.md` - Current progression rules

---

## Risk Assessment

| Factor | Value |
|---|---|
| Likelihood | High |
| Impact if triggered | High |
| Detectability | High |
| Current mitigant | None |

---

## Recommended Fix

**Short term:** Implement base credit earning (1 HP = 1 coin). Create credits field in character schema. Build basic shop UI.

**Medium term:** Add support credit earning (damage mitigation). Implement status effect credit formula. Create skill/equipment purchasing system.

**Long term:** Build advanced shop features (filters, recommendations). Implement credit economy balancing tools. Add credit caps and inflation controls.

---

## References

- `third_party_reply.md` - Credit earning discussion
- `V2_ARCHITECTURAL_DECISIONS.md` - Credit economy decisions
- `docs/entity_users.atom.md` - User database entity
