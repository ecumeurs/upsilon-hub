---
id: mech_board_generation
human_name: Board Generation Mechanic
type: MECHANIC
version: 1.0
status: REVIEW
priority: CORE
tags: [board, combat, generation]
parents:
  - [[module_game]]
dependents:
  - [[ui_board]]
---

# Board Generation Mechanic

## INTENT
Defines the dimensional constraints and terrain composition for the tactical combat grid where battles take place.

## THE RULE / LOGIC
- Dimensions: The board is a standard grid (rectangle). Its width and height must each be randomly rolled between `5` and `15` tiles inclusive.
- Minimum Size Constraint: The total area (width × height) of the rolled board must be greater than or equal to `50` tiles. If a roll yields an area under 50, it must be rerolled until the condition is met.
- Terrain Obstacles: After dimension generation, randomly selected tiles will be designated as impassable "obstacles." The number of obstacle tiles will be up to a maximum of `10%` of the total board area (rounded down).

## TECHNICAL INTERFACE (The Bridge)
- **Code Tag:** `@spec-link [[mech_board_generation]]`
- **Test Names:** `TestBoardDimensions`, `TestBoardMinArea`, `TestBoardObstacleLimit`

## EXPECTATION (For Testing)
- Generate Board -> Width and Height are both between 5 and 15 -> Area is >= 50.
- Board size is 100 tiles -> Number of generated obstacles is between 0 and 10.
