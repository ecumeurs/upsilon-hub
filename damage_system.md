# Upsilon Combat & Damage System (Proposed)

This document outlines the "Tactical Mindset" for the Upsilon combat engine. This is a design specification for future implementation.

## 1. Core Computation Sequence
Damage resolution follows this strict tactical order:

1.  **Hit Test**: `Random(0, 100) < (Attacker.Accuracy - Target.Dodge)`.
2.  **Raw Physical Damage**: `RawPhys = (Attacker.Attack * Skill.Damage / 100) * CritMultiplier`.
    *   *Note: Status powers (Poison/Stun) do NOT add to this HP-damaging pool.*
3.  **Shield Step (The Bubble)**:
    *   `PhysAfterShield = max(RawPhys - Target.Shield, 0)`
    *   `Target.Shield = max(Target.Shield - RawPhys, 0)`
4.  **Poise Step (The Stance)**:
    *   **Backstab Check**: If hit originates from behind (see Section 2), skip halving.
    *   **Halving**: If `Target.Poise > 0`, `PhysForMitigation = floor(PhysAfterShield / 2)`.
    *   **Depletion**: `Target.Poise = max(Target.Poise - PhysAfterShield, 0)`.
5.  **Mitigation Step (The Guard)**:
    *   `FinalHPLevelDmg = max(PhysForMitigation - Target.Armor - Target.Defense, 0)`
6.  **Resolution**:
    *   **HP**: `Target.HP = Target.HP - FinalHPLevelDmg`.

## 2. Positioning & Backstabs
Strategic positioning allows bypassing Poise.
*   **Backstab Condition**: An attacker delivers a hit from the 90° quadrant directly behind the target.
*   **Effect**: Poise halving is bypassed. The hit deals full damage (after Shield/Mitigation). Poise is NOT depleted.

## 3. Status Effects (Poison & Stun)
Status effects are secondary applications tied to the success of the physical hit.

### Poison
*   **Application**: Only applied if `FinalHPLevelDmg > 0` (The weapon must break the skin).
*   **Value**: `max(PoisonPower - Defense, 0)`.
*   **Lifecycle**:
    *   **Damage**: Deals HP damage at **End of Turn** equal to the Poison value.
    *   **Floor**: Cannot reduce HP below 1.
    *   **Decay**: Halves each turn after damage.

### Stun
*   **Application**: Applied if the physical hit lands (Accuracy check) AND the Shield is broken or bypassed.
*   **Condition**: If `PhysAfterShield > 0` (Shielding mitigates everything, including the "stagger" of Stun).
*   **Value**: `max(StunPower - Armor, 0)`.
*   **Lifecycle**:
    *   **Condition**: If `StunCount > (MaxHP / 2)` at **Beginning of Turn**, the turn is skipped.
    *   **Decay**: Halves each turn if no skip occurs.

---
*Status: Design Specification (Draft)*
