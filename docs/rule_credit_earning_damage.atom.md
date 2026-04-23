---
id: rule_credit_earning_damage
human_name: Credit Earning from Damage Rule
type: RULE
layer: ARCHITECTURE
version: 2.0
status: DRAFT
priority: 5
tags: [economy, credits, combat]
parents:
  - [[domain_credit_economy]]
dependents: []
---

# Credit Earning from Damage Rule

## INTENT
To establish the base credit earning mechanism where 1 HP of absolute damage dealt equals 1 credit earned, with healing also earning credits.

## THE RULE / LOGIC
**Base Credit Rule:** 1 HP absolute damage = 1 credit

**Damage Credits:**
- When dealing damage: Credits += damage amount
- Example: Deal 15 damage = Earn 15 credits
- Applies to all damage sources (attacks, skills, poison, etc.)

**Healing Credits:**
- When healing HP: Credits += healing amount
- Example: Heal 10 HP = Earn 10 credits
- Supports active playstyles and rewards healing contribution

**Absolute Damage Definition:**
- Use damage before mitigation/shields
- Poison damage counts as absolute damage
- Stun damage does not earn credits (status effect, not HP damage)

**Credit Assignment:**
- Credits assigned to damage dealer/healer
- Tracked per character per match
- Added to character's total credits after match completion

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[rule_credit_earning_damage]]`
- **Test Names:** `TestDamageCreditEarning`, `TestHealingCreditEarning`
