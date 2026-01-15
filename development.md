# Developer Documentation

This document provides a guide for developers working on the BD_Bot project.

## 1. Local Development Setup

This section should detail the steps required to set up a local development environment. This includes:

*   **Prerequisites:** List all the tools and technologies that need to be installed on a developer's machine (e.g., Go, Node.js, PostgreSQL, etc.).
*   **Installation:** Provide a step-by-step guide on how to clone the repository and install dependencies.
*   **Configuration:** Explain how to set up the `config.yaml` file for local development, including how to connect to a local PostgreSQL instance and a mock LLM server.
    *   *Note for macOS (container platform):* If using the `container` platform on macOS, `localhost` may not resolve to the container. Use `container ls` to find the container's IP address (e.g., `192.168.64.2`) and update the `database_url` in `config.yaml` accordingly.
*   **Running the Application:** Provide commands to run the backend server and frontend application in development mode.

## 2. Architectural Deep-Dive

This section should provide a more detailed explanation of the project's architecture. This includes:

*   **Backend Architecture:** A detailed breakdown of the Go backend, including the structure of the web server, the role of the CLI and TUI, and how background jobs are managed.
*   **Frontend Architecture:** A detailed breakdown of the ReactJS frontend, including the component structure, state management, and how it communicates with the backend.
*   **Database Schema:** A description of the PostgreSQL database schema, including the purpose of each table and the relationships between them.
*   **LLM Integration:** A detailed explanation of how the application integrates with the internal LLM, including the structure of the API requests and responses.

## 3. Coding Conventions

This section should outline the coding conventions that all developers should follow. This includes:

*   **Go Backend:**
    *   Code style (e.g., `gofmt`).
    *   Naming conventions.
    *   Error handling.
    *   Testing best practices.
*   **React Frontend:**
    *   Code style (e.g., Prettier, ESLint).
    *   Component structure.
    *   State management patterns.
    *   Testing best practices.
*   **Git Workflow:**
    *   Branching strategy (e.g., GitFlow).
    *   Commit message format.
    *   Pull request process.