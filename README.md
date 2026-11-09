<p align="center">
  <img src="assets/logo.png" alt="SkillNODE logo" width="130" />
</p>

# <p align="center">SkillNODE</p>
<p align="center">
  A local-first cognitive and technical execution platform powered by an in-browser Go WebAssembly core with an asynchronous real-time telemetry conduit.
</p>

<p align="center">
  <img src="https://img.shields.io/badge/go-1.22+-00ADD8?logo=go&logoColor=white" alt="Go" />
  <img src="https://img.shields.io/badge/webassembly-WASM-654FF0?logo=webassembly&logoColor=white" alt="WebAssembly" />
  <img src="https://img.shields.io/badge/realtime-Socket.IO-010101?logo=socket.io&logoColor=white" alt="Socket.IO" />
  <img src="https://img.shields.io/badge/database-PostgreSQL%2016-4169E1?logo=postgresql&logoColor=white" alt="PostgreSQL" />
  <img src="https://img.shields.io/badge/orm-Prisma%206-2D3748?logo=prisma&logoColor=white" alt="Prisma" />
  <img src="https://img.shields.io/badge/license-VL--PRL-blue" alt="License" />
</p>

<p align="center">
  <a href="https://skillnode-azure.vercel.app"><strong>Live Application</strong></a>
  &nbsp;|&nbsp;
  <a href="https://skillnode-api.onrender.com/healthz"><strong>Telemetry Conduit Health</strong></a>
</p>

---

## Architectural Philosophy

SkillNODE departs from bloated client-side JavaScript frameworks. The entire application interface, state lifecycle, and 6 training tracks execute inside an in-browser **Go WebAssembly (WASM)** virtual runtime compiled from Go templates. 

The client ships with **zero npm dependencies** and requires **zero node_modules**. State persistence is managed through an in-memory **HMAC-SHA256 Cryptographic Vault** with hardware telemetry entropy to prevent local manipulation.

The backend is intentionally reduced to a lightweight, asynchronous **telemetry conduit and real-time message broker (dumb pipe)** — responsible strictly for peer room multiplexing via Socket.IO, database synchronization via Prisma, and external API relays (OpenAI hints, Resend transactional codes).

---

## Workspace Segregation

| Workspace | Runtime / Subsystem | Architectural Role |
|---|---|---|
| `client/` | Go 1.22 + Go Templates (`.tmpl`) + WebAssembly (`app.wasm`) | **Primary Application Core** — Static site generation, SPA routing, 6 training engines (Typing, Quant, Code, Grammar, Estimation, Reading), HMAC-SHA256 integrity vault, zero node_modules. |
| `server/` | Node.js 22 + Socket.IO 4 + Prisma 6 | **Asynchronous Telemetry Conduit** — Decoupled message broker for multiplayer room synchronization, PostgreSQL persistence, and transactional email relays. |

---

## Core Capabilities

- **Typing Velocity Engine** — Real-time character differential highlighter, 120ms cadence sampling ticker, WPM calculation, error tracking, and window blur fair-play detection.
- **Quant Sprint (Math Clash)** — 45-second high-pressure quantitative reasoning and arithmetic sprint with streak multiplier logic.
- **Algorithmic Code Arena** — Multi-language problem sets (JavaScript, Python, C++, TypeScript, SQL) with in-browser JavaScript evaluation sandbox and automated test assertion checklists.
- **Battle Net Real-Time Arena** — Room preset switching, private room code matching, peer presence broadcasting, and synchronized messaging via Socket.IO.
- **Cryptographic Storage Vault** — In-memory HMAC-SHA256 signed storage envelopes with device telemetry entropy to prevent local tampering and client-side privilege escalation.
- **Operator Profile Studio** — Local-first player identity with cryptographic role verification, discipline radar, and shareable performance records.

---

## Repository Structure

```
SkillNODE/
├── client/
│   ├── src/                    # Go WebAssembly core (compiled to wasm/app.wasm)
│   │   ├── api.go              # Asynchronous backend telemetry & health monitor
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
│   │   ├── ai/                 # OpenAI coaching conduit
│   │   ├── auth/               # Cryptographic JWT token utilities
│   │   ├── db/                 # Prisma client instance
│   │   ├── email/              # Transactional email dispatcher (Resend)
│   │   ├── middleware/         # CORS, Helmet, and auth guards
│   │   ├── realtime/           # Socket.IO room lifecycle and peer events
│   │   ├── routes/             # Telemetry & auth routing tables
│   │   ├── storage/            # Static upload handling
│   │   ├── env.js              # Strict environment variable validation
│   │   └── server.js           # Realtime broker entrypoint
│   ├── package.json            # Conduit dependencies and lifecycle scripts
│   └── package-lock.json       # Deterministic conduit lockfile
├── assets/                     # Shared branding assets (logo.png)
├── docs/
│   └── readme/                 # Technical interface documentation captures
├── .github/
│   └── workflows/
│       ├── bot-manifest.yml    # Cryptographic commit audit ledger action
│       └── ci-cd.yml           # Automated build & verification pipeline
├── bot-manifest.json           # Authenticated commit history ledger
├── Dockerfile                  # Production container definition for Render runtime
├── .dockerignore               # Docker build context exclusions
├── docker-compose.yml          # PostgreSQL 16 local service container
├── Makefile                    # Standardized build and deployment automation
├── render.yaml                 # Render web service deployment blueprint (Docker runtime)
├── vercel.json                 # Vercel static build routing configuration
├── package.json                # Root monorepo workspace definition
├── LICENSE                     # Voxion Labs Proprietary Research License
└── README.md                   # System documentation
```

---

## Interface Reference

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

### Telemetry Conduit Development

```bash
# Start local PostgreSQL database container
make db-up

# Push Prisma schema to the database
make db-push

# Start the realtime conduit server (http://localhost:5000)
make dev-server
```

---

## Deployment Pipeline

| Target | Service | Configuration | Deployment Hook |
|---|---|---|---|
| Frontend Client | Vercel Static Hosting | `vercel.json` | Automated via GitHub integration, serving pre-compiled `client/dist` |
| Telemetry Conduit | Render Web Service | `render.yaml` | Automated via GitHub integration, running Node.js real-time service |

---

## License

This project is proprietary software under the **Voxion Labs Proprietary Research License (VL-PRL)**.  
See [LICENSE](LICENSE) for the full terms.

---
<div align="center">
  (c) 2026 Voxion Labs & Rudranarayan Jena
</div>
