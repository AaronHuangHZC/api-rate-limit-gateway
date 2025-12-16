# API Rate Limiting Gateway

A production-ready HTTP reverse-proxy gateway with per-API-key rate limiting, built with Go and Next.js.

## Architecture

```
Clients → Load Balancer → Gateway Instances (N) → Redis (rate limit state)
                                                → Postgres (config + analytics)
                                                → Upstream Services
                                                → Next.js Dashboard
```

## Features

- ✅ Per-API-key + per-route rate limiting (Sliding Window & Token Bucket)
- ✅ Distributed consistency via Redis + Lua scripts (atomic operations)
- ✅ Configurable failure policy (fail-open vs fail-closed)
- ✅ Structured logging + Prometheus metrics
- ✅ Admin dashboard for key management and analytics

## Tech Stack

- **Gateway**: Go 1.21+, Chi router
- **Rate Limiting**: Redis 7+ with Lua scripts
- **Config Storage**: PostgreSQL 16
- **Dashboard**: Next.js 14, TypeScript, Tailwind CSS
- **Observability**: Prometheus metrics, structured JSON logs (zerolog)

## Quick Start

### Prerequisites

- Go 1.21+
- Docker & Docker Compose
- Node.js 18+ (for dashboard)

### 1. Start Infrastructure

```bash
make docker-up
```

This starts Redis and PostgreSQL with health checks.

### 2. Run Migrations

```bash
# (Will be available after Phase 1.4)
make migrate-up
```

### 3. Run Gateway

```bash
make run
```

The gateway will listen on `http://localhost:8080` by default.

### 4. Run Dashboard

```bash
cd dashboard
npm install
npm run dev
```

Dashboard will be available at `http://localhost:3000`.

## Development

See `Makefile` for available commands:

```bash
make help
```

## Project Status

🚧 **Under Active Development**

Currently implementing Phase 1 (Foundation).

