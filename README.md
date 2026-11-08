<p align="center">
  <img src="assets/logo.png" alt="SkillNODE logo" width="130" />
</p>

# <p align="center">SkillNODE</p>
<p align="center">
  A high-throughput skill training and evaluation platform built on a Go WebAssembly client architecture with a decoupled Node.js real-time backend.
</p>

<p align="center">
  <img src="https://img.shields.io/badge/go-1.22+-00ADD8?logo=go&logoColor=white" alt="Go" />
  <img src="https://img.shields.io/badge/webassembly-WASM-654FF0?logo=webassembly&logoColor=white" alt="WebAssembly" />
  <img src="https://img.shields.io/badge/node-22.x-339933?logo=node.js&logoColor=white" alt="Node" />
  <img src="https://img.shields.io/badge/express-4-000000?logo=express&logoColor=white" alt="Express" />
  <img src="https://img.shields.io/badge/prisma-6-2D3748?logo=prisma&logoColor=white" alt="Prisma" />
  <img src="https://img.shields.io/badge/PostgreSQL-16-4169E1?logo=postgresql&logoColor=white" alt="PostgreSQL" />
  <img src="https://img.shields.io/badge/Socket.IO-4-010101?logo=socket.io&logoColor=white" alt="Socket.IO" />
  <img src="https://img.shields.io/badge/license-VL--PRL-blue" alt="License" />
</p>

<p align="center">
  <a href="https://skillnode-azure.vercel.app"><strong>Live Application</strong></a>
  &nbsp;|&nbsp;
  <a href="https://skillnode-api.onrender.com/healthz"><strong>API Health Status</strong></a>
</p>

---

## System Overview

SkillNODE is a dual-runtime monorepo engineered for low-latency cognitive and technical execution training. The client runtime runs on a compiled Go WebAssembly engine rendered through Go HTML templates, eliminating client-side JavaScript bundle overhead and node_modules dependencies. Data integrity and role assignments are cryptographically enforced in the browser through an HMAC-SHA256 vault.

The monorepo is partitioned into two independent workspaces:

| Workspace | Runtime / Architecture | Core Responsibility |
|---|---|---|
| `client/` | Go 1.22 + Go Templates (`.tmpl`) + WebAssembly (`app.wasm`) | Zero-dependency client runtime, static site generation, SPA navigation, 6 training engines, local HMAC-SHA256 vault |
| `server/` | Node.js 22 + Express 4 + Prisma 6 + Socket.IO 4 | REST API, PostgreSQL persistence, AI coaching bridge, real-time WebSocket room coordination, telemetry ingestion |

---

## Core Capabilities

- **Typing Velocity Engine** — Real-time character differential rendering, 120ms cadence sampling, WPM calculation, error tracking, and window blur fair-play detection.
- **Quant Sprint (Math Clash)** — 45-second high-pressure arithmetic and quantitative reasoning challenges with streak multipliers.
- **Algorithmic Code Arena** — Multi-language problem sets (JavaScript, Python, C++, TypeScript, SQL) with in-browser JavaScript evaluation sandbox and automated test assertion checklists.
- **Battle Net Multiplayer** — Room preset switching, private room code matching, peer presence broadcasting, and synchronized messaging via Socket.IO.
- **Cryptographic Storage Vault** — In-memory HMAC-SHA256 signed storage envelopes with device telemetry entropy to prevent local tampering and client-side privilege escalation.
- **Operator Profile Studio** — Local-first player identity with cryptographic role verification, discipline radar, and shareable performance records.

---

## Technical Architecture

```
SkillNODE/
├── client/
│   ├── src/                    # Go WebAssembly core (compiled to wasm/app.wasm)
│   │   ├── api.go              # Asynchronous backend syncing & health monitor
│   │   ├── coding.go           # Algorithmic challenge evaluation engine
│   │   ├── drills.go           # Math Clash quant sprint (45s ticker)
│   │   ├── main.go             # WASM entrypoint and subsystem bootstrap
│   │   ├── multiplayer.go      # Battle Net room coordination & presence
│   │   ├── store.go            # HMAC-SHA256 encrypted local vault & zero-trust role guards
│   │   ├── typing.go           # High-precision typing cadence & diff engine
│   │   └── ui.go               # SPA routing, theme coordinator, dynamic DOM bindings
│   ├── templates/              # Go html/template source files
│   │   ├── head.tmpl           # Metadata, fonts, and favicon configuration
│   │   ├── styles.tmpl         # Purged monolithic CSS stylesheet (40.8 KB)
│   │   ├── layout.tmpl         # Base viewport frame and shell wrappers
│   │   ├── app_header.tmpl     # Telemetry status bar and quick action triggers
│   │   ├── app_sidebar.tmpl    # Navigation, theme toggle, and operator status
│   │   ├── dashboard.tmpl      # Performance metrics and session history feed
│   │   ├── activities.tmpl     # Training track registry
│   │   ├── arena_typing.tmpl   # Typing velocity execution surface
│   │   ├── arena_math.tmpl     # Quantitative sprint surface
│   │   ├── arena_code.tmpl     # Multi-language algorithmic workspace
│   │   ├── arena_multiplayer.tmpl # Real-time Battle Net lobby and live chat
│   │   ├── profile.tmpl        # Operator credentials studio
│   │   ├── auth_modal.tmpl     # Session initialization modal
│   │   ├── app_toast.tmpl      # Glassmorphic notification container
│   │   ├── loading.tmpl        # Startup spinner and WebAssembly loader
│   │   ├── scripts.tmpl        # Bootstrap orchestration script
│   │   └── wasm_exec.tmpl      # Go runtime WebAssembly bridge
│   ├── static/
│   │   └── bridge.js           # Frozen DOM & Socket.IO client interface (5.6 KB)
│   ├── build.go                # Sub-second Go SSG compiler & development server
│   ├── go.mod                  # Go module definition (go 1.22)
│   ├── package.json            # Client build script mappings (0 dependencies)
│   ├── vercel.json             # Vercel deployment rewrites
│   └── logo.png                # Client logo asset
├── server/
│   ├── prisma/
│   │   └── schema.prisma       # Database schema and relational models
│   ├── src/
│   │   ├── ai/                 # OpenAI coaching engine interface
│   │   ├── auth/               # JWT token generation and validation
│   │   ├── controllers/        # Request handling and domain logic
│   │   ├── db/                 # Prisma client instance
│   │   ├── email/              # Transactional email dispatcher (Resend)
│   │   ├── middleware/         # CORS, Helmet, and auth guards
│   │   ├── models/             # Domain schemas and data structures
│   │   ├── realtime/           # Socket.IO room lifecycle and peer events
│   │   ├── routes/             # REST routing tables
│   │   ├── services/           # Background email triggers and AI services
│   │   ├── storage/            # Local static upload handling
│   │   ├── env.js              # Strict environment variable validation
│   │   └── server.js           # HTTP + WebSocket server initialization
│   ├── package.json            # Server dependencies and lifecycle scripts
│   └── package-lock.json       # Deterministic server dependency lockfile
├── assets/                     # Shared branding assets (logo.png)
├── docs/
│   └── readme/                 # Technical interface documentation captures
├── .github/
│   └── workflows/
│       └── bot-manifest.yml    # Cryptographic commit audit ledger action
├── bot-manifest.json           # Authenticated commit history ledger
├── docker-compose.yml          # PostgreSQL 16 local service container
├── Makefile                    # Standardized build and deployment automation
├── render.yaml                 # Render web service deployment blueprint
├── vercel.json                 # Vercel static build routing configuration
├── package.json                # Root monorepo workspace definition
├── LICENSE                     # Voxion Labs Proprietary Research License
└── README.md                   # System documentation
```

---

## Screenshots

<table>
  <tr>
    <td width="50%" valign="top">
      <img src="docs/readme/demo-landing.png" alt="SkillNODE Access Gateway & Landing" />
      <br />
      <strong>1. Access Gateway & Landing</strong>
      <br />
      Landing interface with 0-click guest session initialization and theme coordinator.
    </td>
    <td width="50%" valign="top">
      <img src="docs/readme/demo-typing-lab.png" alt="SkillNODE Typing Lab & Telemetry" />
      <br />
      <strong>2. Execution Lab & Fair-Play</strong>
      <br />
      Focused practice arena tracking cadence, real-time WPM, accuracy, and window blur anomalies.
    </td>
  </tr>
  <tr>
    <td width="50%" valign="top">
      <img src="docs/readme/demo-multiplayer.png" alt="SkillNODE Multiplayer Arena" />
      <br />
      <strong>3. Multiplayer Battle Rooms</strong>
      <br />
      Real-time Socket.IO room coordination, instant invite codes, and synchronized peer presence.
    </td>
    <td width="50%" valign="top">
      <img src="docs/readme/demo-profile.png" alt="SkillNODE Profile Studio" />
      <br />
      <strong>4. Operator Profile Studio</strong>
      <br />
      Local-first player identity with cryptographic integrity seal, discipline radar, and session logs.
    </td>
  </tr>
</table>

---

## Security Model

1. **HMAC-SHA256 Storage Envelopes**: Local state entries (`skillnode_profile`, `skillnode_results`) are sealed with a SHA-256 HMAC key held exclusively within Go WebAssembly linear memory. Hardware telemetry entropy is hashed alongside payloads; manual modification via browser DevTools triggers validation failure, logs an audit violation, and quarantines storage.
2. **Zero-Trust Role Segregation**: Roles (`admin`, `staff`, `moderator`) cannot be granted via client storage. Local profiles default strictly to `player` or `guest`. Elevated roles require server-signed cryptographic JWT validation.
3. **Anti-Cheat Validation**: Activity submissions are checked against physical boundaries (WPM threshold <= 260, minimum race durations >= 1500ms, mathematical score caps, clamp bounds). Unrealistic submissions are dropped prior to storage or synchronization.
4. **JavaScript Sandbox Lockdown**: Browser bridges (`window.SkillNodeBridge`, `window.SkillNodeSocket`) are locked using `Object.freeze()` and defined as non-writable and non-configurable on the `window` object to mitigate console prototype hijacking.

---

## Getting Started

### Prerequisites

| Requirement | Minimum Version |
|---|---|
| Go | `1.22+` |
| Node.js | `22.x` |
| Docker | Latest (optional, for local PostgreSQL) |

### Installation & Build

```bash
# Clone the repository
git clone https://github.com/Voxion-Labs/SkillNODE.git && cd SkillNODE

# Build the client production bundle (Go SSG + WebAssembly)
make build-client

# Start local client development server (http://localhost:8080)
make dev-client
```

### Server Development

```bash
# Start local PostgreSQL database container
make db-up

# Push Prisma schema to the database
make db-push

# Start the Node.js backend server (http://localhost:5000)
make dev-server
```

---

## Deployment Pipeline

| Target | Service | Configuration | Deployment Hook |
|---|---|---|---|
| Frontend | Vercel Static Hosting | `vercel.json` | Automated via GitHub integration, building from `client/dist` |
| Backend API | Render Web Service | `render.yaml` | Automated via GitHub integration, running Node.js runtime |

---

## License

This project is proprietary software under the **Voxion Labs Proprietary Research License (VL-PRL)**.  
See [LICENSE](LICENSE) for the full terms.

---
<div align="center">
  (c) 2026 Voxion Labs & Rudranarayan Jena
</div>
