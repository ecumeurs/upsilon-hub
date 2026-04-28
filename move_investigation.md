
[{2026-04-28T12:49:23.514Z}] [Bot-01] [36m[1mTACTICAL FEED — MATCH DATA[0m
[{2026-04-28T12:49:23.514Z}] [Bot-01] [2m────────────────────────────────────────[0m
     0 1 2 3 4 5 6
 0 │[2m:[0m [31m[1mZ[0m [31m[1mX[0m [32m[1mC[0m [2m:[0m [2m.[0m [2m.[0m │
 1 │[2m.[0m [2m#[0m [2m:[0m [32m[1mA[0m [2m:[0m [2m#[0m [2m:[0m │
 2 │[2m:[0m [2m#[0m [2m:[0m [2m.[0m [2m:[0m [2m.[0m [2m:[0m │
 3 │[2m:[0m [2m:[0m [2m:[0m [2m:[0m [2m:[0m [2m:[0m [2m:[0m │
 4 │[2m:[0m [2m:[0m [2m:[0m [2m:[0m [2m:[0m [32m[1mB[0m [2m:[0m │
 5 │[2m:[0m [2m#[0m [2m:[0m [2m:[0m [2m:[0m [31m[1m[42mY[0m [2m:[0m │
 6 │[2m.[0m [2m:[0m [2m:[0m [2m:[0m [2m:[0m [2m:[0m [2m:[0m │

  [1mID  UNIT NAME       OWNER        HP/MAX     MVT     DELAY [0m
  [2m──────────────────────────────────────────────────────────────────────[0m
  [32m[1mA[0m Herald          credit_bot_3686_0 30/30      3/3     252  
  [32m[1mB[0m Saboteur        credit_bot_3686_0 29/30      3/3     36   
  [32m[1mC[0m Slayer          credit_bot_3686_0 30/30      3/3     286  
  [31m[1mX[0m Entropy_125a    Echo_e04e    3/3        2/2     43   
[36m> [0m[31m[1mY[0m Null_Zero_666   Echo_e04e    5/5        2/2     0    
  [31m[1mZ[0m Ghost_Vermin_Alpha Echo_e04e    5/5        1/1     75   


  [2mSuggested next steps:[0m [32mredraw[0m
[{2026-04-28T12:49:23.514Z}] [Bot-01] Received board update, but it's not my turn yet. Continuing wait...

[{2026-04-28T12:49:23.610Z}] [Bot-01] [35m[1m[WS][0m turn.started event received.
  [2m{
    "request_id": "019dd423-1a98-70ab-95c4-2a65368a70e1",
    "message": "Board Updated",
    "success": true,
    "data": {
      "match_id": "019dd423-1909-7104-88f2-9e5f23883eff",
      "players": [
        {
          "nickname": "credit_bot_3686_0",
          "entities": [
            {
              "id": "019dd423-17ee-73d8-b288-766fa93194c6",
              "team": 1,
              "name": "Herald",
              "hp": 30,
              "max_hp": 30,
              "attack": 10,
              "defense": 5,
              "move": 3,
              "max_move": 3,
              "position": {
                "x": 3,
                "y": 1
              },
              "equipped_items": null,
              "buffs": [],
              "equipped_skills": [],
              "is_self": true,
              "dead": false
            },
            {
              "id": "019dd423-17f1-720f-a91f-c0609d46b867",
              "team": 1,
              "name": "Saboteur",
              "hp": 29,
              "max_hp": 30,
              "attack": 10,
              "defense": 5,
              "move": 3,
              "max_move": 3,
              "position": {
                "x": 5,
                "y": 4
              },
              "equipped_items": null,
              "buffs": [],
              "equipped_skills": [],
              "is_self": true,
              "dead": false
            },
            {
              "id": "019dd423-17f3-70f5-8778-83f69f26edeb",
              "team": 1,
              "name": "Slayer",
              "hp": 30,
              "max_hp": 30,
              "attack": 10,
              "defense": 5,
              "move": 3,
              "max_move": 3,
              "position": {
                "x": 3,
                "y": 0
              },
              "equipped_items": null,
              "buffs": [],
              "equipped_skills": [],
              "is_self": true,
              "dead": false
            }
          ],
          "team": 1,
          "ia": false,
          "is_self": true
        },
        {
          "nickname": "Echo_e04e",
          "entities": [
            {
              "id": "0b6a5439-dd54-4f06-951c-a4b75b91ebcb",
              "team": 2,
              "name": "Entropy_125a",
              "hp": 3,
              "max_hp": 3,
              "attack": 2,
              "defense": 3,
              "move": 2,
              "max_move": 2,
              "position": {
                "x": 2,
                "y": 0
              },
              "equipped_items": null,
              "buffs": [],
              "equipped_skills": [],
              "is_self": false,
              "dead": false
            },
            {
              "id": "42a4e4b0-02df-45a3-af78-3396ab6613b2",
              "team": 2,
              "name": "Null_Zero_666",
              "hp": 5,
              "max_hp": 5,
              "attack": 2,
              "defense": 1,
              "move": 2,
              "max_move": 2,
              "position": {
                "x": 5,
                "y": 5
              },
              "equipped_items": null,
              "buffs": [],
              "equipped_skills": [],
              "is_self": false,
              "dead": false
            },
            {
              "id": "ab2e8860-bb3c-49d4-8475-d48a4e7fe513",
              "team": 2,
              "name": "Ghost_Vermin_Alpha",
              "hp": 5,
              "max_hp": 5,
              "attack": 2,
              "defense": 2,
              "move": 1,
              "max_move": 1,
              "position": {
                "x": 1,
                "y": 0
              },
              "equipped_items": null,
              "buffs": [],
              "equipped_skills": [],
              "is_self": false,
              "dead": false
            }
          ],
          "team": 2,
          "ia": true,
          "is_self": false
        }
      ],
      "grid": {
        "width": 7,
        "height": 7,
        "max_height": 2,
        "cells": [
          [
            {
              "entity_id": null,
              "obstacle": false,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 0
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 0
            }
          ],
          [
            {
              "entity_id": "ab2e8860-bb3c-49d4-8475-d48a4e7fe513",
              "obstacle": false,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": true,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": true,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": true,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 1
            }
          ],
          [
            {
              "entity_id": "0b6a5439-dd54-4f06-951c-a4b75b91ebcb",
              "obstacle": false,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 1
            }
          ],
          [
            {
              "entity_id": "019dd423-17f3-70f5-8778-83f69f26edeb",
              "obstacle": false,
              "height": 1
            },
            {
              "entity_id": "019dd423-17ee-73d8-b288-766fa93194c6",
              "obstacle": false,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 0
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 1
            }
          ],
          [
            {
              "entity_id": null,
              "obstacle": false,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 1
            }
          ],
          [
            {
              "entity_id": null,
              "obstacle": false,
              "height": 0
            },
            {
              "entity_id": null,
              "obstacle": true,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 0
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 1
            },
            {
              "entity_id": "019dd423-17f1-720f-a91f-c0609d46b867",
              "obstacle": false,
              "height": 1
            },
            {
              "entity_id": "42a4e4b0-02df-45a3-af78-3396ab6613b2",
              "obstacle": false,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 1
            }
          ],
          [
            {
              "entity_id": null,
              "obstacle": false,
              "height": 0
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 1
            }
          ]
        ]
      },
      "turn": [
        {
          "delay": 7,
          "entity_id": "0b6a5439-dd54-4f06-951c-a4b75b91ebcb",
          "is_self": false,
          "team": 0
        },
        {
          "delay": 39,
          "entity_id": "ab2e8860-bb3c-49d4-8475-d48a4e7fe513",
          "is_self": false,
          "team": 0
        },
        {
          "delay": 216,
          "entity_id": "019dd423-17ee-73d8-b288-766fa93194c6",
          "is_self": true,
          "team": 1
        },
        {
          "delay": 250,
          "entity_id": "019dd423-17f3-70f5-8778-83f69f26edeb",
          "is_self": true,
          "team": 1
        },
        {
          "delay": 364,
          "entity_id": "42a4e4b0-02df-45a3-af78-3396ab6613b2",
          "is_self": false,
          "team": 0
        }
      ],
      "current_entity_id": "019dd423-17f1-720f-a91f-c0609d46b867",
      "timeout": "2026-04-28T12:49:53.56605791Z",
      "start_time": "2026-04-28T12:49:23.56605781Z",
      "winner_team_id": null,
      "version": 8589934592,
      "current_player_is_self": true,
      "game_finished": false
    },
    "meta": {}
  }[0m

[{2026-04-28T12:49:23.652Z}] [Bot-01] [35m[1m[WS][0m board.updated event received.
  [2m{
    "request_id": "019dd423-1ac3-72b2-9b59-876f796e3989",
    "message": "Board Updated",
    "success": true,
    "data": {
      "match_id": "019dd423-1909-7104-88f2-9e5f23883eff",
      "players": [
        {
          "nickname": "credit_bot_3686_0",
          "entities": [
            {
              "id": "019dd423-17ee-73d8-b288-766fa93194c6",
              "team": 1,
              "name": "Herald",
              "hp": 30,
              "max_hp": 30,
              "attack": 10,
              "defense": 5,
              "move": 3,
              "max_move": 3,
              "position": {
                "x": 3,
                "y": 1
              },
              "equipped_items": null,
              "buffs": [],
              "equipped_skills": [],
              "is_self": true,
              "dead": false
            },
            {
              "id": "019dd423-17f1-720f-a91f-c0609d46b867",
              "team": 1,
              "name": "Saboteur",
              "hp": 29,
              "max_hp": 30,
              "attack": 10,
              "defense": 5,
              "move": 3,
              "max_move": 3,
              "position": {
                "x": 5,
                "y": 4
              },
              "equipped_items": null,
              "buffs": [],
              "equipped_skills": [],
              "is_self": true,
              "dead": false
            },
            {
              "id": "019dd423-17f3-70f5-8778-83f69f26edeb",
              "team": 1,
              "name": "Slayer",
              "hp": 30,
              "max_hp": 30,
              "attack": 10,
              "defense": 5,
              "move": 3,
              "max_move": 3,
              "position": {
                "x": 3,
                "y": 0
              },
              "equipped_items": null,
              "buffs": [],
              "equipped_skills": [],
              "is_self": true,
              "dead": false
            }
          ],
          "team": 1,
          "ia": false,
          "is_self": true
        },
        {
          "nickname": "Echo_e04e",
          "entities": [
            {
              "id": "0b6a5439-dd54-4f06-951c-a4b75b91ebcb",
              "team": 2,
              "name": "Entropy_125a",
              "hp": 3,
              "max_hp": 3,
              "attack": 2,
              "defense": 3,
              "move": 2,
              "max_move": 2,
              "position": {
                "x": 2,
                "y": 0
              },
              "equipped_items": null,
              "buffs": [],
              "equipped_skills": [],
              "is_self": false,
              "dead": false
            },
            {
              "id": "42a4e4b0-02df-45a3-af78-3396ab6613b2",
              "team": 2,
              "name": "Null_Zero_666",
              "hp": 5,
              "max_hp": 5,
              "attack": 2,
              "defense": 1,
              "move": 2,
              "max_move": 2,
              "position": {
                "x": 5,
                "y": 5
              },
              "equipped_items": null,
              "buffs": [],
              "equipped_skills": [],
              "is_self": false,
              "dead": false
            },
            {
              "id": "ab2e8860-bb3c-49d4-8475-d48a4e7fe513",
              "team": 2,
              "name": "Ghost_Vermin_Alpha",
              "hp": 5,
              "max_hp": 5,
              "attack": 2,
              "defense": 2,
              "move": 1,
              "max_move": 1,
              "position": {
                "x": 1,
                "y": 0
              },
              "equipped_items": null,
              "buffs": [],
              "equipped_skills": [],
              "is_self": false,
              "dead": false
            }
          ],
          "team": 2,
          "ia": true,
          "is_self": false
        }
      ],
      "grid": {
        "width": 7,
        "height": 7,
        "max_height": 2,
        "cells": [
          [
            {
              "entity_id": null,
              "obstacle": false,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 0
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 0
            }
          ],
          [
            {
              "entity_id": "ab2e8860-bb3c-49d4-8475-d48a4e7fe513",
              "obstacle": false,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": true,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": true,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": true,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 1
            }
          ],
          [
            {
              "entity_id": "0b6a5439-dd54-4f06-951c-a4b75b91ebcb",
              "obstacle": false,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 1
            }
          ],
          [
            {
              "entity_id": "019dd423-17f3-70f5-8778-83f69f26edeb",
              "obstacle": false,
              "height": 1
            },
            {
              "entity_id": "019dd423-17ee-73d8-b288-766fa93194c6",
              "obstacle": false,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 0
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 1
            }
          ],
          [
            {
              "entity_id": null,
              "obstacle": false,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 1
            }
          ],
          [
            {
              "entity_id": null,
              "obstacle": false,
              "height": 0
            },
            {
              "entity_id": null,
              "obstacle": true,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 0
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 1
            },
            {
              "entity_id": "019dd423-17f1-720f-a91f-c0609d46b867",
              "obstacle": false,
              "height": 1
            },
            {
              "entity_id": "42a4e4b0-02df-45a3-af78-3396ab6613b2",
              "obstacle": false,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 1
            }
          ],
          [
            {
              "entity_id": null,
              "obstacle": false,
              "height": 0
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 1
            },
            {
              "entity_id": null,
              "obstacle": false,
              "height": 1
            }
          ]
        ]
      },
      "turn": [
        {
          "delay": 7,
          "entity_id": "0b6a5439-dd54-4f06-951c-a4b75b91ebcb",
          "is_self": false,
          "team": 0
        },
        {
          "delay": 39,
          "entity_id": "ab2e8860-bb3c-49d4-8475-d48a4e7fe513",
          "is_self": false,
          "team": 0
        },
        {
          "delay": 216,
          "entity_id": "019dd423-17ee-73d8-b288-766fa93194c6",
          "is_self": true,
          "team": 1
        },
        {
          "delay": 250,
          "entity_id": "019dd423-17f3-70f5-8778-83f69f26edeb",
          "is_self": true,
          "team": 1
        },
        {
          "delay": 364,
          "entity_id": "42a4e4b0-02df-45a3-af78-3396ab6613b2",
          "is_self": false,
          "team": 0
        }
      ],
      "current_entity_id": "019dd423-17f1-720f-a91f-c0609d46b867",
      "timeout": "2026-04-28T12:49:53.611002456Z",
      "start_time": "2026-04-28T12:49:23.611002286Z",
      "winner_team_id": null,
      "action": {
        "type": "pass",
        "actor_id": "42a4e4b0-02df-45a3-af78-3396ab6613b2"
      },
      "version": 8589934592,
      "current_player_is_self": true,
      "game_finished": false
    },
    "meta": {}
  }[0m
[{2026-04-28T12:49:23.653Z}] [Bot-01] [33m[1m[SYSTEM][0m Tactical feed updated.

[{2026-04-28T12:49:23.653Z}] [Bot-01] [36m[1mTACTICAL FEED — MATCH DATA[0m
[{2026-04-28T12:49:23.653Z}] [Bot-01] [2m────────────────────────────────────────[0m
  [35m[1m[FEEDBACK][0m Unit [31m[1mY[0m passed their turn
  [2m────────────────────────────────────────[0m
     0 1 2 3 4 5 6
 0 │[2m:[0m [31m[1mZ[0m [31m[1mX[0m [32m[1mC[0m [2m:[0m [2m.[0m [2m.[0m │
 1 │[2m.[0m [2m#[0m [2m:[0m [32m[1mA[0m [2m:[0m [2m#[0m [2m:[0m │
 2 │[2m:[0m [2m#[0m [2m:[0m [2m.[0m [2m:[0m [2m.[0m [2m:[0m │
 3 │[2m:[0m [2m:[0m [2m:[0m [2m:[0m [2m:[0m [2m:[0m [2m:[0m │
 4 │[2m:[0m [2m:[0m [2m:[0m [2m:[0m [2m:[0m [32m[1m[42mB[0m [2m:[0m │
 5 │[2m:[0m [2m#[0m [2m:[0m [2m:[0m [2m:[0m [31m[1mY[0m [2m:[0m │
 6 │[2m.[0m [2m:[0m [2m:[0m [2m:[0m [2m:[0m [2m:[0m [2m:[0m │

  [1mID  UNIT NAME       OWNER        HP/MAX     MVT     DELAY [0m
  [2m──────────────────────────────────────────────────────────────────────[0m
  [32m[1mA[0m Herald          credit_bot_3686_0 30/30      3/3     216  
[36m> [0m[32m[1mB[0m Saboteur        credit_bot_3686_0 29/30      3/3     0    
  [32m[1mC[0m Slayer          credit_bot_3686_0 30/30      3/3     250  
  [31m[1mX[0m Entropy_125a    Echo_e04e    3/3        2/2     7    
  [31m[1mY[0m Null_Zero_666   Echo_e04e    5/5        2/2     364  
  [31m[1mZ[0m Ghost_Vermin_Alpha Echo_e04e    5/5        1/1     39   


  [2mSuggested next steps:[0m [32mredraw[0m
[{2026-04-28T12:49:23.654Z}] [Bot-01] Initiative shift:  -> 019dd423-17f1-720f-a91f-c0609d46b867. Clearing turn memory.
[{2026-04-28T12:49:23.654Z}] [Bot-01] --- My Turn! Acting with Saboteur (credit_bot_3686_0, team 1) (v8589934592) ---
[{2026-04-28T12:49:23.654Z}] [Bot-01] Auto Action: Move (5,3) -> (5,2) -> (4,2)

[{2026-04-28T12:49:23.654Z}] [Bot-01] [36m[1m[CURL][0m [2mcurl -X POST -H 'Accept: application/json' -H 'Authorization: Bearer 32|wRqrkvV7LVxfzNyGD6tNwiwASkeS5z65O4L17SzQcd7cb941' -H 'Content-Type: application/json' -d '{"request_id":"019dd423-1ac6-7ab9-bcf4-4ff0c6ef666c","message":"CLI Request: POST /api/v1/game/019dd423-1909-7104-88f2-9e5f23883eff/action","success":true,"data":{"entity_id":"019dd423-17f1-720f-a91f-c0609d46b867","player_id":"","target_coords":[{"x":5,"y":3},{"x":5,"y":2},{"x":4,"y":2}],"type":"move"},"meta":{}}' 'http://127.0.0.1:8000/api/v1/game/019dd423-1909-7104-88f2-9e5f23883eff/action'[0m
[{2026-04-28T12:49:23.687Z}] [Bot-01] [31m[1m[REPLY 400][0m 
  [2m{
    "request_id": "019dd423-1ac6-7ab9-bcf4-4ff0c6ef666c",
    "message": "Invalid path",
    "success": false,
    "data": [],
    "meta": {
      "error_key": "entity.path.notvalid"
    }
  }[0m
[{2026-04-28T12:49:23.687Z}] [Bot-01] [CALL_ERROR] Route game_action failed: Invalid path
[{2026-04-28T12:49:23.687Z}] [Bot-01] JS Exception: [object Object]
	at github.com/ecumeurs/upsiloncli/internal/script.(*Agent).jsAutoBattleTurn-fm (native)
	at <eval>:85:31(232)
