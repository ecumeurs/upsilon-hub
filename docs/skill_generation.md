---
id: req_skill_generation_overhaul
human_name: Skill Generation — Design Spec
type: REQUIREMENT
layer: BUSINESS
version: 1.0
status: DRAFT
priority: 5
tags: [skills, generation, categories, icons, ui, progression]
parents:
  - [[req_tech_debt_backlog]]
dependents:
  - [[domain_skill_system]]
  - [[mech_skill_selection_progression]]
  - [[ui_action_panel]]
---

# Skill Generation — Design Spec

## INTENT

Upgrade the skill roll pipeline from a structureless stat-generator into a category-aware, tag-driven system with diegetic names, neon polygon icons, and proper action slots in the Battle Arena. Every rolled skill must be readable — a player should understand what a skill does from its icon stack and name without opening a tooltip.

---

## 1. Cost → Grade Reference

Source of truth: `upsilontypes/entity/skill/skillweight/skillweight.go` + atom `[[shared:rule_skill_grading_system]]`.

| Grade | Positive SW Band | Credit Cost (PSW × 2) | Unlock Window (player wins) |
|-------|------------------|-----------------------|-----------------------------|
| I     | 0 – 150          | 0 – 300               | 0+ (creation & early)       |
| II    | 151 – 300        | 302 – 600              | 0+ (creation & early)       |
| III   | 301 – 500        | 602 – 1000             | 10+ wins                    |
| IV    | 501 – 750        | 1002 – 1500            | 20+ wins                    |
| V     | 750+             | 1500+                  | 30+ wins                    |

> **Cadence decision:** Player wins is the single progression metric. There is no character-level counter. Slot count and grade-unlock window both follow `[[rule_character_skill_slots]]` (1 slot base, +1 per 10 wins, cap 5). Atom `[[mech_skill_selection_progression]]` should be updated to state unlock windows in player-wins terms.
>
> **Auto-roll:** Stays **manual** for now. The roulette is triggered from `CharacterDetailModal.vue` only. Layout rework is pending separately.

### Skill Weight rules (positive SW)

| Effect / Property | SW gain |
|---|---|
| Damage | +1 SW per 1% (100 Damage = +100 SW) |
| Heal | +15 SW per 10 heal |
| ShieldPower | +10 SW per point |
| Range > 1 | +10 SW per extra cell |
| Zone cells > 1 | +50 SW per extra cell |
| TargetingMechanics = Anywhere | +40 SW flat |
| StunChance | +2 SW per 1% |
| CriticalChance | +2 SW per 1% |
| PoisonPower | +15 SW per dmg/turn |

Negative SW payments (Delay, Channeling, Leeches, Cooldown) reduce net SW but do **not** reduce grade — grade is computed on **positive SW only**. The generator closes to `netSW = 0` by adjusting Delay.

---

## 2. What Is Rollable at Grade 1 Today

The Go generator (`upsilontypes/entity/skill/skillgenerator/skillgenerator.go`) is a narrow slice of the type system. Each bucket has a 50% inclusion roll; the effect loop forces ≥ 1 effect.

| Bucket | Property | Rolled Range | PSW impact |
|---|---|---|---|
| Targeting | `Accuracy` | 50–150 | 0 (not in SW table) |
| Effect (one of) | `Damage` | 50–200 | ≤ 200 PSW |
| Effect (one of) | `Heal` | 50–150 | ≤ 225 PSW |
| Effect (one of) | `ShieldPower` | 10–50 | ≤ 500 PSW |
| Cost | `Cooldown` | 1–5 turns | −25/turn |
| Cost | `HPLeech` | 1–10 | −20/HP |
| Cost | `MPLeech` | 1–10 | −15/MP |
| Cost | `SPLeech` | 1–10 | −10/SP |

**Key gaps:**
- The generator has **no grade target** — it can easily overshoot grade I (Damage > 150, Shield ≥ 16, Heal ≥ 110 all push past PSW 150).
- `Behavior` is always `Direct` (default, never randomised). No Reaction / Passive / Counter / Trap can roll.
- No `Zone`, no `Range > 1`, no `StunChance`, `PoisonChance`, `CriticalChance`, no `TargetType` variation.
- Skill name = raw property key (`"Damage"`, `"Heal"`, `"Shield"`). Not diegetic.

---

## 3. Category Taxonomy

A skill carries an ordered list of **tags** derived from its final properties. Tags drive icon composition, name generation, and action-panel rendering. They are **not authored** — they are inferred post-generation.

### 3.1 Tag Vocabulary

| Tag | Classification Rule |
|---|---|
| `melee` | Range == 1 AND has Damage or StunPower |
| `ranged` | Range ≥ 2 AND has Damage |
| `aoe` | Zone cell-count > 1 |
| `heal` | has positive Heal effect |
| `shield` | has positive ShieldPower effect |
| `buff` | Duration > 0 AND TargetType ∈ {Self, FriendOnly} AND any positive non-damage effect |
| `debuff` | Duration > 0 AND TargetType = EnemyOnly AND negative status effect |
| `dot` | PoisonPower > 0 |
| `stun` | StunPower > 0 OR StunChance > 0 |
| `crit` | CriticalChance ≥ 25 |
| `trap` | Behavior == Trap |
| `counter` | Behavior == Counter |
| `reaction` | Behavior == Reaction |
| `passive` | Behavior == Passive |
| `mobility` | TargetType == Self AND modifies Movement with Duration > 0 |
| `channeled` | Channeling > 0 |
| `instant` | Delay ≤ 100 AND Channeling == 0 |

Tags are stored in `instance_data.tags []string`. No DB column change needed (instance_data is JSON).

### 3.2 Tag Ordering

Order determines icon stacking and naming priority:

1. **Behavior tag** — if non-Direct (`trap`, `counter`, `reaction`, `passive`).
2. **Effect family** — `heal`, `shield`, `dot`, `stun`, `buff`, `debuff`.
3. **Delivery** — `melee`, `ranged`, `aoe`.
4. **Modifiers** — `crit`, `channeled`, `instant`, `mobility`.

The first tag drives the major icon glyph; subsequent tags add overlays.

### 3.3 Multi-tag examples

| Example | Tags (ordered) | Reading |
|---|---|---|
| Friendly trap that heals on entry | `[trap, heal]` | Trap behavior > healing effect |
| Ranged bolt that buffs the caster | `[buff, ranged]` | Buff effect > ranged delivery |
| Poison cloud | `[aoe, dot, debuff]` | AoE delivery > DoT > debuff modifier |
| Melee backstab crit | `[crit, melee]` | Critical modifier > melee delivery |
| Self-heal passive | `[passive, heal]` | Passive behavior > healing effect |

### 3.4 Deferred: Siphon

A `siphon` category was considered — skills that drain a stat from the target and gain it on the caster (e.g. −1 CritChance on target, +1 CritChance on caster). This requires paired / mirrored effects that `effect.Effect` cannot currently express. Deferred until a compound-effect system is designed.

---

## 4. Generator Architecture (Go)

Replace `GenerateRandomSkill()` with a **grade-aware, category-dispatched** system.

### 4.1 File layout

```
upsilontypes/entity/skill/skillgenerator/
  skillgenerator.go      ← new dispatcher + Generate() entrypoint
  blueprint.go           ← SW-budget DSL (AddDamage, SetRange, AddZone…)
  classifier.go          ← Skill → []string tags
  namegen.go             ← (primaryTag, secondaryTags, grade) → name
  producer_melee.go
  producer_ranged.go
  producer_aoe.go
  producer_heal.go
  producer_shield.go
  producer_buff.go
  producer_debuff.go
  producer_trap.go
  producer_counter.go
  producer_reaction.go
  producer_passive.go
  producer_dot.go
  producer_stun.go
  producer_mobility.go
```

### 4.2 Entrypoint

```go
type GenerateRequest struct {
    TargetGrade string   // "I"…"V"; empty defaults to "I"
    AllowedTags []string // empty = any category
    ForbidTags  []string // exclude categories
}

// Generate returns skill, ordered tags, error.
// Replaces GenerateRandomSkill(); kept as alias: Generate(GenerateRequest{TargetGrade:"I"})
func Generate(req GenerateRequest) (skill.Skill, []string, error)
```

Dispatcher flow:
1. Compute PSW budget from `TargetGrade` (see table §4.3).
2. Pick primary producer uniformly from `AllowedTags` (or all producers if empty).
3. Primary producer builds the skill within `[budgetLo, budgetHi]` PSW.
4. With grade-dependent probability (§4.3) invoke a secondary producer to layer additional properties, staying in band.
5. `skillweight.Calculate` + Delay closer → `netSW = 0`.
6. `classifier.Classify(sk)` → ordered tag list.
7. `namegen.Name(tags[0], tags[1:], grade)` → skill name.
8. Return.

### 4.3 PSW budget per grade

| Grade | PSW band | Secondary layer chance |
|---|---|---|
| I    | 60 – 150    | 25% |
| II   | 151 – 300   | 50% |
| III  | 301 – 500   | 65% |
| IV   | 501 – 750   | 75% |
| V    | 751 – 1000  | 85% |

### 4.4 Producer responsibilities

| Producer | Behavior set | Primary effect | Targeting |
|---|---|---|---|
| `melee` | Direct | Damage / StunPower | Range=1, EnemyOnly |
| `ranged` | Direct | Damage | Range 2–4, EnemyOnly, optionally LOS |
| `aoe` | Direct (layered) | — | Zone pattern (Neighbours / Line) |
| `heal` | Direct | Heal | FriendOnly |
| `shield` | Direct | ShieldPower | Self or FriendOnly |
| `buff` | Direct | stat boost + Duration | Self or FriendOnly |
| `debuff` | Direct | negative stat + Duration | EnemyOnly |
| `dot` | Direct (layered) | PoisonPower + PoisonChance | — |
| `stun` | Direct | StunPower / StunChance | EnemyOnly |
| `trap` | Trap | any | TargetType=Tile, TriggerType set |
| `counter` | Counter | Damage or shield | Self or EnemyOnly |
| `reaction` | Reaction | Damage or heal | depends |
| `passive` | Passive | HP regen / stat constant | Self |
| `mobility` | Direct | Movement buff + Duration | Self |

### 4.5 API changes

- `upsilonapi/handler/skill_generate.go` — accept optional JSON body `{ "grade": "II", "allowed_tags": ["heal","ranged"] }`.
- `upsilonapi/api/output.go` — add `Tags []string` to `SkillGenerateResponse`.
- `battleui/app/Services/SkillGeneratorBridge.php` — pass optional `grade` and `allowed_tags` from the roll endpoint.
- `battleui/app/Http/Controllers/API/CharacterSkillController.php` — `roll` may accept `?grade=I` query param; enforces that the requested grade is within the character's level window.

---

## 5. Name Generation

### 5.1 Template

```
[Modifier prefix]  [Subject]  [Suffix]
```

Max 24 characters total. Segments may be empty.

### 5.2 Dictionaries

**Modifier prefix** — driven by secondary tag(s) (first match wins):

| Secondary tag | Prefix pool |
|---|---|
| `dot` | `Cinder`, `Sludge`, `Rot_`, `Verm_` |
| `crit` | `Razor`, `Fang_`, `Spike_` |
| `aoe` | `Flux`, `Cascade_`, `Wave_` |
| `channeled` | `Drift_`, `Bleed_`, `Slow_` |
| `buff` / `debuff` | `Echo_`, `Static`, `Hex_` |
| none | ∅, `Null_`, `Void_`, `Ghost_` (haxxor pool) |

**Subject** — driven by primary tag:

| Primary tag | Subject pool |
|---|---|
| `melee` | Strike, Bash, Cleaver, Smash |
| `ranged` | Bolt, Lance, Pulse, Tracer |
| `aoe` | Burst, Field, Storm, Bloom |
| `heal` | Mend, Patch, Pulse, Suture |
| `shield` | Bulwark, Aegis, Plate, Shell |
| `trap` | Mine, Snare, Tripwire, Hex |
| `counter` | Riposte, Ricochet, Rebuke |
| `reaction` | Reflex, Backlash |
| `passive` | Aura, Cycle, Drift |
| `stun` | Jolt, Stutter, Lockdown |
| `mobility` | Sprint, Phase, Vector |

**Suffix** — pick one:
- Grade flavored: `_I`, `_II`, … (30% chance)
- Haxxor: `_X`, `v2`, `_Z`, `_Bot`, `_666`, `_Alpha` (70% chance)

### 5.3 Sample outputs

| Tags | Sample name |
|---|---|
| `[trap, dot]` | `Cinder Mine_X` |
| `[ranged]` | `Void_Bolt v2` |
| `[heal, passive]` | `Null_Suture _I` |
| `[aoe, stun, debuff]` | `Flux Jolt_666` |
| `[melee, crit]` | `Razor Strike_Z` |
| `[shield]` | `Ghost_Aegis v2` |

### 5.4 Implementation

File: `upsilontypes/entity/skill/skillgenerator/namegen.go`

```go
func Name(primaryTag string, secondaryTags []string, grade string) string
```

Called at the end of `Generate()` before returning. Name is written to `skill.Name`. The atom `[[mech_skill_name_generation]]` (new stub) documents the dictionaries as the authoritative source.

---

## 6. Icon System

### 6.1 Neon polygon vocabulary

Each tag maps to an SVG path glyph. Glyphs are stroked only (no fill), 1.5px stroke, neon drop-shadow glow.

| Tag | Glyph description | Theme color |
|---|---|---|
| `melee` | Crossed daggers (×) | magenta `#ff00ff` |
| `ranged` | Chevron / arrow (▶) | cyan `#00f2ff` |
| `aoe` | Concentric hex rings | cyan `#00f2ff` |
| `heal` | Plus cross (+) | lime `#39ff13` |
| `shield` | Hexagon outline | cyan `#00f2ff` |
| `buff` | Upward triangle (△) | lime `#39ff13` |
| `debuff` | Downward triangle (▽) | magenta `#ff00ff` |
| `dot` | Dripping rhombus | lime `#39ff13` |
| `stun` | Zigzag bolt (⚡) | amber `#fbbf24` |
| `crit` | Star / asterisk (*) | amber `#fbbf24` |
| `trap` | Trapezoid mine outline | amber `#fbbf24` |
| `counter` | Mirrored arrows (⇄) | magenta `#ff00ff` |
| `reaction` | Broken / bounced arrow | magenta `#ff00ff` |
| `passive` | Infinity loop (∞) | steel `#4a4a4f` |
| `mobility` | Sprint chevrons (»›) | cyan `#00f2ff` |
| `channeled` | Hourglass | steel `#4a4a4f` |
| `instant` | Lightning chevron | cyan `#00f2ff` |

### 6.2 Composition rule

| Layer | Tag index | Size in action slot | Size in detail panel | Opacity |
|---|---|---|---|---|
| Major glyph | tags[0] | 32 × 32 | 64 × 64 | 100% |
| Minor glyph | tags[1] | 14 × 14, top-right | 24 × 24 | 70% |
| Tertiary glyph (grade ≥ III only) | tags[2] | 10 × 10, bottom-right | 16 × 16 | 50% |

Background: thin hex-border ring, colored by grade:
`I = steel` · `II = cyan` · `III = lime` · `IV = magenta` · `V = amber`

### 6.3 Implementation

- New component: `battleui/resources/js/Components/Skill/SkillIcon.vue`
  ```vue
  <SkillIcon :tags="['trap','dot']" :grade="'II'" />
  ```
- New registry: `battleui/resources/js/Components/Skill/skillIconRegistry.js`
  — JS object mapping tag → inline SVG `<path>` string. No asset pipeline.
- Update `SkillCard.vue` to swap the current behavior ASCII glyph for `<SkillIcon>`.
- Update `SkillDetail.vue` to render `<SkillIcon>` at 64px and display `tags[]` as a pill row.
- Update `SkillSlotPill.vue` to show `<SkillIcon>` at 20px beside the skill name.
- New atom stub: `[[ui_skill_icon]]` parented to `[[req_ui_look_and_feel]]` + `[[ui_theme]]`.

---

## 7. Battle Arena — Action Panel

### 7.1 Split layout

`ActionPanel.vue` gains two new zones between the fixed actions and Forfeit:

**Active skill row** — skills with `behavior ∈ {Direct, Reaction, Counter, Trap}`:
- `<SkillIcon>` 32 × 32 on the left.
- Stacked label: name (Orbitron 8px) + cost summary (`−3 MP · CD 2`) in JetBrains Mono 7px cyan.
- Hover: `<SkillDetail>` tooltip in a neon panel.
- Disabled when on cooldown or insufficient resource.
- Click → `emit('action', { type: 'skill', skillId })`.

**Passive rail** — skills with `behavior == Passive`:
- Smaller icon (20 × 20), name only, `cursor: default`.
- No cost row, no click handler, no disabled state.
- Faint pulse animation to signal "always active".
- Hover: `<SkillDetail>` tooltip for inspection only.
- Renders below the active row, separated by a thin divider.

> Reaction and Counter skills remain in the **active row** — the player still arms them; the engine decides when they fire.

### 7.2 Engine wiring (flagged, not specced here)

The action panel needs the equipped skills' current cooldown state from the engine snapshot. The skill action route already handles dispatching (`[[api_go_battle_action]]`). Equipping and passing skills to the engine at arena init is handled by `[[entity_character_skill_inventory]]` + `[[api_battle_proxy]]`.

---

## 8. Test Coverage

### 8.1 Go unit tests

File additions / extensions in `upsilontypes/entity/skill/skillgenerator/`:

| File | Tests |
|---|---|
| `skillgenerator_test.go` | Extend with `TestGenerate_GradeBand` (100 skills per grade, PSW in band), `TestGenerate_TagsNonEmpty` |
| `producer_melee_test.go` | Range==1, TargetType==EnemyOnly, has Damage/Stun effect |
| `producer_heal_test.go` | has Heal effect, TargetType==FriendOnly |
| *(one per producer)* | Each producer's contractual guarantees |
| `namegen_test.go` (new) | Table-driven: (primaryTag, secondaryTags, grade) → ≤24 chars, no empty segments, matches template regex |
| `classifier_test.go` (new) | Hand-built skills → expected tag list and ordering |

### 8.2 CLI E2E scenario

New file: `upsiloncli/tests/scenarios/e2e_skill_roll_naming.js`

Flow:
1. Create account + character via CLI helpers.
2. `POST /api/v1/profile/character/{id}/skills/roll` (existing endpoint).
3. Assert `instance_data.name` matches diegetic-name pattern (≤ 24 chars, no bare property keys).
4. Assert `instance_data.tags` is a non-empty array drawn from the 17-tag vocabulary.
5. Roll 20× and print `name | grade | tags` table for visual spot-check.

Run with: `scripts/trigger_one_ci_test.sh tests/scenarios/e2e_skill_roll_naming.js`

### 8.3 Playwright (allowed to fail)

Extend `battleui/tests/playwright/user_flows.spec.ts`:

- **Block: skill icon + name visible** — open dashboard → `CharacterDetailModal` → `SkillRouletteModal` → spin → stop → assert `<SkillIcon>` rendered and name passes regex. Mark `test.fixme` until UI rework stabilises.
- **Block: action panel passive split** — stub only with `test.fixme`; full spec when `ActionPanel.vue` rework lands.

---

## 9. Implementation Work Items

Listed as task-stubs for the implementation session (not tonight):

| Item | File(s) | Depends on |
|---|---|---|
| `blueprint.go` — SW-budget DSL | `skillgenerator/` | — |
| `classifier.go` — tag inference | `skillgenerator/` | blueprint |
| `namegen.go` — name generation | `skillgenerator/` | classifier |
| Per-producer files (×14) | `skillgenerator/` | blueprint |
| New `Generate()` dispatcher | `skillgenerator/skillgenerator.go` | all producers |
| API grade/tag params | `upsilonapi/handler/skill_generate.go` | Generate() |
| Tags in API response | `upsilonapi/api/output.go` | — |
| PHP bridge grade passthrough | `SkillGeneratorBridge.php` | API change |
| `SkillIcon.vue` + registry | `battleui/resources/js/Components/Skill/` | — |
| `SkillCard/Detail/SlotPill` updates | same dir | SkillIcon.vue |
| `ActionPanel.vue` skill rows | `Arena/ActionPanel.vue` | SkillIcon.vue |
| `BattleArena.vue` skill props | `Pages/BattleArena.vue` | ActionPanel |
| Atom stubs: `mech_skill_name_generation`, `ui_skill_icon` | `upsilonbattle/docs/`, `battleui/docs/` | this doc |
| Atom update: `mech_skill_selection_progression` | `upsilonbattle/docs/` | cadence decision |

---

## EXPECTATION

- Every rolled skill has a name that does not equal a raw property key (`Damage`, `Heal`, `Shield`).
- Every rolled skill has `instance_data.tags` populated and ordered per §3.2.
- Grade I skills have PSW ∈ [60, 150].
- `ActionPanel.vue` renders active-skill buttons and a passive rail; passives have no click handler.
- `<SkillIcon>` renders the major glyph at grade-border color and up to two overlay glyphs.
- `e2e_skill_roll_naming.js` passes in CI.
- Playwright skill-icon test is present and marked `test.fixme` if the UI is unstable.
