# BD_Bot

Business Development Intelligence Portal - Automating government solicitation discovery and pursuit.

## Overview

BD_Bot is a high-performance portal designed to scrape government solicitation sites, use an internal LLM to match opportunities to user expertise, and provide a collaborative dashboard for teams.

## Documentation

- [Requirements](requirements.md) - Detailed project requirements and architecture.
- [Development Guide](development.md) - Onboarding and development instructions.
- [Gemini Context](GEMINI.md) - High-level overview for AI assistants.

## Quick Start

### 1. Infrastructure
Ensure you have Docker installed and run:
```bash
docker-compose up -d
```

### 2. Build
Build the unified binary:
```bash
# Build frontend
cd web && npm install && npm run build && cd ..
# Build backend
go build -o bd_bot ./cmd/bd_bot
```

### 3. Initialize & Run
```bash
# Setup Config
./bd_bot config init

# Run Migrations
./bd_bot migrate up

# Create Admin User (Standalone Auth)
./bd_bot user create -e admin@example.com -n "Admin User"
./bd_bot user passwd -e admin@example.com -p secret123

# Start Server
./bd_bot serve
```

## Architecture

- **Backend:** Go 1.25.5+
- **Frontend:** React (TypeScript) embedded in the Go binary.
- **Database:** PostgreSQL
- **Authentication:** Dual-mode (Standalone/SSO)
- **CLI/TUI:** Cobra & Charmbracelet (Huh?, Bubble Tea, Lip Gloss)
