---
id: mec_pre_post_execution_costs
human_name: Pre and Post Execution Costs Mechanic
type: MECHANIC
layer: IMPLEMENTATION
version: 2.0
status: DRAFT
priority: 5
tags: [skills, costs, time-based]
parents: []
dependents: []
---

# Pre and Post Execution Costs Mechanic

## INTENT
To implement the dual-cost system for skills with pre-execution costs (SP, MP, Channeling) paid upfront and post-execution costs (Delay) paid after effect completes.

## THE RULE / LOGIC
**Cost Separation:**
- **Pre-Execution Costs:** Paid immediately when skill execution begins
- **Post-Execution Costs:** Applied after skill effect completes

**Pre-Execution Costs:**
- **SP Leech:** SP points consumed upfront
- **MP Leech:** MP points consumed upfront
- **Channeling:** Casting time delay (risk premium cost)
- **Resource Risk:** Pre-execution costs are consumed even if skill is interrupted

**Post-Execution Costs:**
- **Delay:** Recovery delay added after effect completes
- **Example:** Basic attack = +100 delay after dealing damage
- **Timing:** Applied to caster's CurrentDelay after effect execution

**Cost Payment Flow:**
```go
// Skill execution sequence
1. Pre-execution checks (SP, MP, Channeling available?)
2. Deduct pre-execution costs (SP, MP)
3. Create temporary channeling entity if channeling
4. Execute skill effect (when channeling completes)
5. Apply post-execution costs (Delay)
6. Update entity state and return to player
```

**Channeling Special Case:**
- **Pre-Execution:** Pay SP/MP, create channeling entity
- **Delay:** Added after channeling completes and effect executes
- **Total Timeline:** Channeling delay → Effect execution → Recovery delay

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[mec_pre_post_execution_costs]]`
- **Related Files:** `upsilonbattle/battlearena/ruler/rules/skill.go`
