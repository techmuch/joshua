# Developer Documentation

This document provides a comprehensive guide for new developers working on the BD_Bot project.

## 1. Local Development Setup

### Prerequisites
*   **Go:** 1.25.5 or later.
*   **Node.js & npm:** For frontend development.
*   **Container Platform:** Docker or macOS `container` platform.
*   **Make:** For build automation.

### Quick Start
1.  **Clone & Install:**
    ```bash
    git clone github.com:techmuch/bd_bot.git
    cd bd_bot
    make deps
    ```
2.  **Database:**
    ```bash
    make db-up
    ./bd_bot config init
    # Note: On macOS/container, verify DB IP with `container ls` and edit config.yaml
    ./bd_bot migrate up
    ```
3.  **Create Admin User:**
    ```bash
    ./bd_bot user create -e admin@example.com -n "Admin User"
    ./bd_bot user passwd -e admin@example.com -p secret123
    ```
4.  **Run Full Stack:**
    ```bash
    make run
    # Access UI at http://localhost:8080
    ```

### Live Frontend Development (HMR)
For rapid UI iteration:
1.  Start Backend: `make build-go && ./bd_bot serve`
2.  Start Frontend: `cd web && npm run dev`
3.  Access: `http://localhost:5173` (Proxies API calls to backend).

## 2. Project Architecture

### Backend (`/internal`)
*   **`api/`**: REST handlers (`auth.go`, `user.go`, `solicitations.go`).
*   **`cli/`**: Cobra commands (`root.go`, `user.go`, `match.go`).
*   **`repository/`**: PostgreSQL data access logic.
*   **`ai/`**: LLM integration logic.
*   **`scraper/`**: GPR scraping engine.

### Frontend (`/web/src`)
*   **`components/`**:
    *   `UserProfile.tsx`: Unified profile, password, and narrative editor.
    *   `PersonalInbox.tsx`: AI-matched opportunities with analytics.
    *   `SolicitationList.tsx`: Global library view.
    *   `DashboardCharts.tsx`: Reusable Recharts component.
    *   `LoginButton.tsx`: Auth modal and navigation tab.
*   **`context/`**: `AuthContext.tsx` handles session state.
*   **`hooks/`**: `useAnalytics.ts` encapsulates cross-filtering logic.

## 3. Key Workflows

### Authentication
The system supports Dual-Mode Auth:
1.  **Standalone (Local):**
    *   Uses `bcrypt` for password hashing.
    *   Managed via CLI (`user create`, `user passwd`).
    *   UI: Modal with Email/Password inputs.
2.  **SSO (Future):**
    *   Schema supports `auth_provider` and `external_id`.
    *   Mock SSO: If no password is set, login can auto-provision (dev mode).

### Data Pipeline
1.  **Ingestion:** `make scrape` runs the scraper -> DB.
2.  **Matching:** `./bd_bot match [user_id]` runs LLM analysis -> `matches` table.
3.  **Consumption:** User views Inbox -> `PersonalInbox.tsx` renders matches + analytics.

## 4. API Reference

| Method | Endpoint | Description | Auth |
|:---|:---|:---|:---|
| `POST` | `/api/auth/login` | Login (Local or Mock SSO) | No |
| `POST` | `/api/auth/password` | Change Password | Yes |
| `GET` | `/api/solicitations` | List all opportunities | No |
| `GET` | `/api/matches` | List user matches | Yes |
| `PUT` | `/api/user/profile` | Update Name/Email/Avatar | Yes |
| `POST` | `/api/user/avatar` | Upload Profile Picture | Yes |

## 5. Troubleshooting

**Port 8080 already in use:**
*   Check if `bd_bot` is running: `ps aux | grep bd_bot`
*   Kill it: `killall bd_bot`

**Database Connection Failed:**
*   Verify config: `cat config.yaml`
*   Check container IP: `container ls` or `docker ps`
*   Ensure `make db-up` was run.

**Frontend "White Screen":**
*   Check browser console.
*   Ensure backend is running (proxies rely on it).

**CLI "Duplicate Command":**
*   Check `internal/cli/user.go` `init()` for duplicate `AddCommand` calls.

## 6. Coding Standards
*   **Go:** `gofmt`, `goimports`. Use `slog` for logging.
*   **React:** Functional components, Hooks (`useAuth`, `useAnalytics`).
*   **CSS:** Responsive, full-width layouts. Avoid fixed widths in main containers.
