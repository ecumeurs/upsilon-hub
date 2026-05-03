# Dashboard Review: Current State & Improvement Plan

**Target File:** `battleui/resources/js/Pages/Dashboard.vue`
**Status:** Half-functional, fragmented state management.

---

## 1. Current State of Affairs

The Dashboard serves as the central "Tactical Command" hub, orchestrating several critical sub-systems via modals and child components. However, the current implementation suffers from **fragmented data ownership**:

*   **Scattered State:** Each major component (`CharacterRoster`, `CharacterDetailModal`, `InventoryModal`, `ShopModal`) manages its own internal data fetching and state.
*   **Poor Coordination:** Actions taken in one modal (e.g., purchasing an item in `ShopModal`) do not consistently trigger refreshes in other visible components (e.g., `CharacterRoster`'s internal inventory).
*   **Stale Data:** Re-fetching logic is inconsistent. Some components refresh on mount, others on prop changes, but there is no global "Sync" event to keep the dashboard unified.

---

## 2. Entities & Data Flow

### Involved Entities
-   **User:** Owns credits and global stats.
-   **Characters:** Have stats (HP, ATK, etc.), equipment links, and skill loadouts.
-   **Inventory Items:** Linked to ShopItems, can be equipped to character slots (`armor`, `weapon`, `utility`).
-   **Character Skills:** Individual skill instances acquired via Roulette or progression.

### Data Pulling Mechanisms
1.  **Dashboard.vue:** Fetches `user` from auth service on mount. Polls global match stats every 60s via API.
2.  **CharacterRoster.vue:** Fetches all characters and full player inventory on mount.
3.  **CharacterDetailModal.vue:** Fetches specific character details, skills, equipment, and full player inventory whenever it opens or the target character changes.
4.  **InventoryModal.vue / ShopModal.vue:** Fetch their respective manifests (Inventory list / Shop catalog) on open.

---

## 3. Action Manifest & Use Cases

| Feature | Component | Current Logic | Known Failures / Gaps |
| :--- | :--- | :--- | :--- |
| **Character Listing** | `CharacterRoster` | Internal `characters` ref loaded on mount. | Doesn't reflect equipment changes made in `CharacterDetailModal` without manual refresh. |
| **Item Purchase** | `ShopModal` | Calls `shopService.purchase`, emits `credits-updated`. | **Does not refresh inventory**. `CharacterRoster` and `InventoryModal` will show stale data. |
| **Equip Hardware** | `CharacterDetailModal` | Calls `inventoryService.equip`, then re-loads modal data. | UI updates only inside the modal. The background roster remains stale. |
| **Skill Roulette** | `SkillRouletteModal` | Triggered from `CharacterDetailModal`. | Often fails to appear due to `roulette_available` flag check or UI positioning in the side pane. |
| **Skill Slot Management**| `CharacterDetailModal` | Uses `skill_slots` count and `skills` list. | Slots often don't appear if `skill_slots` is not correctly returned or defaulted. |

---

## 4. Why it feels "Half-Functional" (Root Causes)

1.  **Race Conditions & Stale Refs:** Because `CharacterRoster` and `CharacterDetailModal` both fetch the same data (inventory/characters) independently, they frequently fall out of sync.
2.  **Unidirectional Credit Updates:** `credits-updated` only flows from modals to the `Dashboard`'s `user` object. It doesn't propagate back down to other components that might need the updated `user` (or they ignore it because they have their own fetch logic).
3.  **Modal Isolation:** Modals are treated as separate "islands" rather than views into the same shared state.
4.  **UI Clipping:** The Skill Roulette button is placed at the bottom of a scrollable sidebar in the detail modal, making it easily missed or clipped.

---

## 5. Recommended Refactor Path

### Phase 1: Shared State (Composables)
-   Move character and inventory fetching into a shared composable (e.g., `useDashboardData`).
-   Use a single "Source of Truth" for the inventory and roster at the `Dashboard` level.

### Phase 2: Unified Refresh Event
-   Implement a `global-refresh` or `sync-required` event that triggers a re-fetch of all dashboard data (credits, inventory, and roster) after any state-changing action (purchase, equip, roll).

### Phase 3: UI/UX Rework
-   **Skill Roulette:** Move the trigger to a more prominent location or ensure it's always visible (not buried in a scroll list).
-   **Skill Slots:** Ensure `skill_slots` defaults correctly to `1` and handle empty states more gracefully.
-   **Inventory Sync:** Explicitly refresh the inventory manifest after any `ShopModal` purchase.

---

## 6. Verification Checklist for Improvements
- [ ] Purchasing an item immediately shows it in the `Inventory Archive`.
- [ ] Equipping an item in the detail modal immediately updates the character card in the roster.
- [ ] Skill Roulette button visibility is deterministic based on the `roulette_available` flag.
- [ ] Skill slots correctly render based on character capacity.
