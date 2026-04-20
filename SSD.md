# Software Design Document (SSD) - Upsilon Battle

## 1. User Personas & Roles

### 1.1 Persona: The Player
A tactical RPG enthusiast looking for fast-paced, competitive matches with minimal onboarding friction.

### 1.2 Persona: The Administrator
A system maintainer responsible for managing users and ensuring database performance. Does not participate in combat but has high-level control over account states and match history.

### 1.3 Access Control List (ACL) Roles
Currently, the system operates with a single specialized role:
| Role | Responsibility | Data Access |
|---|---|---|
| **User (Player)** | Account management, joining queues, taking combat turns. | Own profile, own character rosters, public leaderboards, authorized game arenas. |
| **Administrator** | User soft deletion, match history review/cleaning. | Systematic overview (account names, win rates), match results. **No access to private address/birth date.** |

---

## 2. Global Architecture

### 2.1 System Architecture Diagram
The system follows a multi-tier "Proxy and Bridge" architecture to isolate complex Go-based combat logic from the web-facing Laravel gateway.

```mermaid
flowchart TD
    UI["BattleUI (Vue.js)"]
    LGW["Laravel Gateway (PHP)"]
    DB[("PostgreSQL DB")]
    UAPI["UpsilonAPI (Go Bridge)"]
    UBG["UpsilonBattle Engine (Go)"]

    UI -- "REST / Auth" --> LGW
    UI -- "WebSockets (State)" --> LGW
    LGW -- "Proxy Actions" --> UAPI
    LGW -- "Cache State" --> DB
    UAPI -- "Orchestrate Arenas" --> UBG
    UBG -- "Webhook Updates" --> LGW
```

### 2.2 Communication Protocols
| Link | Protocol | Description |
|---|---|---|
| **Vue.js <-> Laravel** | HTTP & WebSockets | Authentication (Sanctum) and real-time state streaming (Reverb). |
| **Laravel <-> UpsilonAPI** | HTTP & Webhooks | Proxied game actions and async state updates into the Laravel callback. |
| **UpsilonAPI <-> Engine** | Internal Go Channels | High-performance message passing between Ruler and Controllers. |

---

## 3. Software Entities (Detailed)

### 3.1 Laravel Gateway
- **Responsibility:** Serves as the primary ingress point. Manages authentication, session state, and metadata (characters, wins).
- **Inner Working:** 
    - Proxies combat actions to the Go bridge.
    - Listens for webhooks from Go to update the `game_matches` JSON state cache and broadcast to WebSockets.
- **Constraints:** Must not perform complex combat math; only record results and authorize requests.
- **References:** [[api_laravel_gateway]], [[req_security]].

### 3.2 UpsilonAPI (The Bridge)
- **Responsibility:** Provides an HTTP interface for the stateful Go engine. Orchestrates multiple concurrent arenas.
- **Inner Working:** Maintains a registry of active `Ruler` instances and maps them to `arena_ids`.
- **Constraints:** Must be stateless regarding player identity (delegates to Laravel).
- **References:** [[module_upsilonapi]], [[api_go_battle_engine]]

### 3.3 UpsilonBattle Engine
- **Responsibility:** The core TRPG logic processor.
- **Inner Working:** Implements the **Ruler/Controller** pattern.
    - **Ruler:** Acts as the Game Master, enforcing rules (initiative, timer, collision).
    - **Controller:** Acts as the player/AI interface for issuing move/attack commands.
- **Constraints:** Enforces the 30-second shot clock and +400 delay penalty for timeouts.
- **References:** [[module_game]], [[mech_controller_communication_sequence]], [[mech_action_economy]].

---

## 4. Implementation Architecture Details

### 4.1 Technology Stack
| Component | Technology | Purpose | Key Libraries |
|---|---|---|---|
| **Frontend** | Vue.js 3 + Laravel 10 | User interface & API gateway | Vue Router, Pinia, Laravel Sanctum |
| **Backend API** | Go 1.21+ | High-performance combat engine | Gin, UUID v7, logrus |
| **Database** | PostgreSQL 15+ | Data persistence & caching | - |
| **Real-time** | Laravel Reverb | WebSocket state broadcasting | - |
| **Testing** | Go testing + PHPUnit | Unit & integration tests | - |

### 4.2 Data Flow Architecture
```mermaid
sequenceDiagram
    participant U as User (Vue.js)
    participant L as Laravel Gateway
    participant D as Database
    participant R as Reverb (WebSocket)
    participant A as UpsilonAPI (Go)
    participant E as Battle Engine

    U->>L: POST /auth/login
    L->>D: Validate credentials
    D-->>L: User data + JWT token
    L-->>U: Auth response + token

    U->>L: POST /matchmaking/join
    L->>D: Create queue entry
    L->>A: POST /arena/start
    A->>E: Initialize battle
    E-->>A: Initial state
    A-->>L: Webhook: game.started
    L->>D: Cache match state
    L->>R: Broadcast to private channel
    R-->>U: WebSocket: Match ready

    U->>L: POST /game/{id}/action
    L->>A: POST /arena/{id}/action  
    A->>E: Execute combat logic
    E-->>A: Action result
    A-->>L: Webhook: board.updated
    L->>D: Update cached state
    L->>R: Broadcast state change
    R-->>U: WebSocket: New board state
```

### 4.3 State Management Strategy
- **Database Source of Truth**: Laravel PostgreSQL database holds authoritative match state
- **Engine Processing**: Go Engine processes actions and generates new states
- **Webhook Synchronization**: Engine pushes state updates to Laravel via webhooks
- **Version Control**: Monotonic version numbers prevent race conditions and duplicate processing
- **WebSocket Distribution**: Real-time state changes broadcast to connected clients

### 4.4 Security Architecture
- **Authentication**: Laravel Sanctum Bearer tokens with 15-minute expiration
- **Authorization**: Role-based access control (Player/Admin)
- **Transport Security**: HTTPS mandatory (self-signed certificates permitted for development)
- **Input Validation**: Multi-layer validation (Laravel requests + Go engine validation)
- **Identity Protection**: UUID v7 for internal IDs, Tactical IDs for client references

---

## 5. Requirement Traceability Matrix

| Requirement ID | Business Requirement | Software Component | Implementation Detail | ATD Reference | Status |
|---|---|---|---|---|---|
| **BR-01** | Frictionless Onboarding | Laravel Gateway | Name/pass/address/birth registration. | [[us_new_player_onboard]] | ✅ Complete |
| **BR-02** | GDPR Compliance | Laravel Gateway | Soft-delete and Anonymization hooks. | [[rule_gdpr_compliance]] | 🔄 Partial |
| **BR-03** | Data Portability | Laravel Gateway | `/api/profile/export` endpoint. | [[api_profile_export]] | ✅ Complete |
| **BR-04** | Tactical Combat Engine | UpsilonBattle Engine | Initiative-based turn logic. | [[module_game]] | ✅ Complete |
| **BR-05** | Turn Timeout Penalty | UpsilonBattle Engine | +400 delay cost logic in Ruler. | [[mech_action_economy]] | ✅ Complete |
| **BR-06** | Secure Transport | All Components | Mandatory HTTPS (Self-signed ok). | [[req_security]] | ✅ Complete |
| **BR-07** | Fair Progression | Laravel Gateway | Attribute point allocation gated by wins. | [[rule_progression]] | ✅ Complete |
| **BR-08** | Real-time Updates | Laravel Gateway | Reverb WebSockets broadcasting. | [[api_laravel_gateway]] | ✅ Complete |
| **BR-09** | Identity Safety | UpsilonBattle Engine | Friendly fire detection and blocking. | [[rule_friendly_fire]] | ✅ Complete |
| **BR-10** | System Administration | Laravel Gateway | User listing/deletion; History purge. | [[uc_admin_user_management]] | 🔄 Partial |
| **BR-11** | Admin Privacy Gate | Laravel Gateway | Masking sensitive user fields for admins. | [[rule_admin_access_restriction]] | ✅ Complete |
| **BR-12** | Secure Admin Seeding | Laravel Gateway | env-based admin account creation. | [[infra_seed_admin]] | ✅ Complete |
| **BR-13** | Tactical Action Reporting | UpsilonBattle Engine | Rich state mutation describing combat events. | [[requirement_customer_action_reporting]] | 🔄 Partial |
| **BR-14** | API-First Integration | All Components | Self-documenting registry and 1:1 API parity. | [[requirement_customer_api_first]] | 🔄 Partial |

---

## 6. Component Interaction Details

### 6.1 Laravel Gateway Responsibilities
- **Authentication**: JWT token issuance and validation
- **Session Management**: Token refresh and expiration handling
- **Request Proxying**: Forward combat actions to Go engine
- **State Caching**: Maintain authoritative match state in database
- **WebSocket Broadcasting**: Distribute state changes to connected clients
- **Business Logic**: Enforce progression rules, character limits, reroll restrictions

### 6.2 UpsilonAPI Bridge Responsibilities
- **Arena Orchestration**: Manage multiple concurrent battle instances
- **Request Routing**: Direct actions to appropriate engine instances
- **Webhook Delivery**: Send asynchronous state updates to Laravel
- **Health Monitoring**: Provide liveness and readiness probes
- **State Validation**: Ensure request integrity before engine processing

### 6.3 Battle Engine Responsibilities
- **Turn Management**: Initiative calculation and turn sequencing
- **Action Validation**: Move, attack, and ability validation
- **Combat Resolution**: Damage calculation, status effects, win detection
- **Timer Enforcement**: 30-second shot clock with auto-pass penalty
- **State Generation**: Create new board states after each action

---

## 7. Performance & Scalability Considerations

### 7.1 Current Limitations
- **Single Instance**: Not designed for horizontal scaling
- **Memory Usage**: Each arena maintains full state in memory
- **Database Load**: Frequent state updates can cause write contention
- **WebSocket Connections**: No connection pooling or load balancing

### 7.2 Optimization Opportunities
- **State Caching**: Redis layer for high-frequency state reads
- **Batch Processing**: Aggregate multiple actions before webhook delivery
- **Connection Management**: WebSocket connection pooling and reconnection logic
- **Database Indexing**: Optimize queries for leaderboard and match history

### 7.3 Monitoring & Observability
- **Health Checks**: `/health` endpoints for all services
- **Request Tracing**: UUID v7 request IDs across all services
- **Error Logging**: Structured logging with request context
- **Performance Metrics**: Turn processing time, action latency, state update frequency
