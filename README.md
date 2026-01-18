# BD_Bot: Business Development Intelligence Portal

BD_Bot is an automated intelligence platform designed to revolutionize how organizations discover and pursue government business opportunities. It combines high-performance web scraping with local Large Language Models (LLM) to deliver personalized, actionable leads directly to your inbox.

## Key Features

*   ** automated Discovery:** A background bot scrapes government solicitation portals (e.g., Georgia Procurement Registry) daily.
*   **Intelligent Matching:** An internal AI analyzes your organizational "Narrative" against thousands of opportunities to score relevance.
*   **Collaborative Dashboard:** A centralized portal for teams to review, claim, and track pursuits.
*   **Dual-Mode Authentication:** Supports both standalone (local password) and Enterprise SSO integration.

---

## 1. Installation & Setup

### Prerequisites
*   **Docker** (for PostgreSQL database)
*   **Go 1.25+** (for backend)
*   **Node.js 20+** (for frontend build)

### Step 1: Clone and Build
```bash
git clone https://github.com/techmuch/bd_bot.git
cd bd_bot

# Install Go dependencies
make deps

# Build the unified binary (Frontend + Backend)
make build
```

### Step 2: Infrastructure
Start the database container:
```bash
make db-up
```

### Step 3: Configuration
Run the interactive setup wizard to generate `config.yaml`:
```bash
./bd_bot config init
```
*   *Tip:* If running on macOS/Container, use `container ls` to find the DB IP address.

### Step 4: Database Schema
Apply the latest migrations:
```bash
./bd_bot migrate up
```

---

## 2. Administration

### User Management
BD_Bot includes a CLI for managing users, useful for initial setup or admin tasks.

**Create a new user:**
```bash
./bd_bot user create -e admin@example.com -n "System Admin"
```

**Set a password (Standalone Auth):**
```bash
./bd_bot user passwd -e admin@example.com -p StrongPassword123!
```

**List all users:**
```bash
./bd_bot user list
```

### Scraper Management
Trigger the scraper manually or check its status:
```bash
./bd_bot scraper run-now
```

---

## 3. Usage Guide

### Log In
1.  Start the server: `./bd_bot serve`
2.  Navigate to `http://localhost:8080`.
3.  Click **Login** and enter your credentials (or use SSO if configured).

### Set Your Narrative (Crucial!)
The AI needs to know who you are to find matches.
1.  Click on your **User Profile** (top right).
2.  Scroll down to **Business Capability Narrative**.
3.  Paste your capabilities statement, past performance, or core competencies.
4.  Click **Save Narrative**.

### View Matches
1.  Navigate to the **Inbox**.
2.  Review opportunities scored by the AI.
3.  Use the **Analytics Dashboard** to filter by urgency or agency.

---

## 4. Deployment

BD_Bot is designed as a **Single Binary Deployment**. The React frontend is embedded into the Go binary.

1.  **Build:** Run `make build` to generate the `bd_bot` binary.
2.  **Deploy:** Copy `bd_bot` and `config.yaml` to your server.
3.  **Run:** Execute `./bd_bot serve`.
    *   *Production Tip:* Use a systemd service or Docker container to keep it running.
    *   *Env Vars:* Configuration can be overridden with environment variables (see `config.go`).

## Documentation Links
- [Requirements & Architecture](requirements.md)
- [Developer Guide (Internals)](development.md)