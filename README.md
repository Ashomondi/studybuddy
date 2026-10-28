# StudyBuddy AI

StudyBuddy AI is a learning companion built with a Go REST API, SQLite, and a React/Vite frontend. Development is organized as small, runnable increments.

## Sprint 1: application foundation

This increment provides a responsive welcome screen, a Go API health endpoint, SQLite connection/configuration, and a schema-migrations table. Goal creation and AI features will be added in later increments.

## Run locally

Requirements: Go 1.22+, Node.js 20+, npm, and a C compiler for the SQLite driver.

Start the API from the repository root:

```sh
go run ./cmd/api
```

It listens on `:8080` and creates `studybuddy.db` in the current directory. Set `API_ADDR` or `DATABASE_PATH` to override these defaults. Check `http://localhost:8080/api/health` for API/database status.

Start the frontend in another terminal:

```sh
cd web
npm install
npm run dev
```

Vite proxies `/api` requests to the Go server. Open the URL printed by Vite (usually `http://localhost:5173`).

## Configuration

Configuration is read from environment variables. AI credentials will be added to the Go backend configuration when an AI provider is selected; they must never be placed in frontend code.