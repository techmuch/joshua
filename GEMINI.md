# BD_Bot (Business Development Intelligence Portal)

## Project Overview

The BD_Bot is a high-performance, cross-platform portal designed to automate the discovery and pursuit of government business development opportunities. It features a background bot for scraping solicitation sites, an internal Large Language Model (LLM) for matching opportunities to user expertise, and a collaborative dashboard for organizational "coalition building" and transparency.

**Main Technologies:**
*   **Backend:** Go (1.25.5+) single binary.
*   **Frontend:** ReactJS Single Page Application (SPA) embedded via `go:embed`.
*   **CLI & TUI:** Cobra framework, Charmbracelet/huh (interactive forms), Bubble Tea (rich terminal feedback).
*   **Database:** PostgreSQL.
*   **External AI:** Internal, OpenAI-compatible LLM server.

**Architecture:**
The application follows a client-server architecture with a Go backend serving a ReactJS frontend. A CLI/TUI provides administrative and monitoring capabilities. Data persistence is handled by PostgreSQL, and an internal LLM is used for intelligent opportunity matching.

## Building and Running

The project utilizes a unified Go binary that embeds the frontend.

**Getting Started:**

1.  **Clone the repository:**
    ```bash
    git clone <repository-url>
    cd bd_bot
    ```
2.  **Generate configuration:** This command will assist in setting up the `config.yaml` with interactive forms.
    ```bash
    bd_bot config init
    ```
    For silent mode, generating default config:
    ```bash
    bd_bot config init --silent
    ```
3.  **Run Database Migrations:**
    ```bash
    bd_bot migrate up
    ```
4.  **Start the Portal:** This will serve the web application and run background processes.
    ```bash
    bd_bot serve
    ```

**CLI & TUI Commands:**

*   **Scraper:** Manually trigger a scraper run or monitor active processes:
    ```bash
    bd_bot scraper run-now
    ```
*   **User Management:** Manage identities and roles with a list view and search:
    ```bash
    bd_bot user
    ```

## Development Conventions

A comprehensive guide for developers is available in the `development.md` file. This includes:
*   Local development setup
*   Architectural deep-dive
*   Coding conventions (Go, React, Git)

## Key Files

*   `requirements.md`: Detailed project requirements, application architecture, and feature specifications.
*   `development.md`: A comprehensive guide for new developers, capturing all necessary details to facilitate a smooth onboarding process.
*   `GEMINI.md`: This file, providing a high-level overview and quick start guide for the project.

## Quality Assurance & CI/CD

*   **Unit Testing:** Development must include a comprehensive suite of unit tests for both the Go backend and React frontend.
*   **CI/CD Pipeline:** A fully automated pipeline is required for linting, security scanning, automated testing, building the unified binary, and artifact storage.
*   **Deployment Gate:** The application must pass all unit tests and pipeline stages before being considered for deployment.
