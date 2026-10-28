# StudyBuddy AI

StudyBuddy AI is a learning companion built with a Go REST API, SQLite, and a React/Vite frontend. Development is organized as small, runnable increments.

## Sprint 1: application foundation

This increment provides a responsive welcome screen, a Go API health endpoint, SQLite connection/configuration, and a schema-migrations table. Goal creation and AI features will be added in later increments.

## Run locally

Requirements: Go 1.22+, Node.js 20+, npm, and a C compiler for the SQLite driver.

The backend and frontend run independently — each can be started alone.

Start the API from the repository root:

```sh
go run ./cmd/api
```

It listens on `:8080` and creates `studybuddy.db` in the current directory. Check `http://localhost:8080/api/health` for API/database status.

Start the frontend in another terminal:

```sh
cd web
npm install
npm run dev
```

Open the URL printed by Vite (usually `http://localhost:5173`). The frontend calls the API directly at `http://localhost:8080` by default.

## Configuration

Configuration is read from environment variables:

- `API_ADDR` (default `:8080`) — API listen address
- `DATABASE_PATH` (default `studybuddy.db`) — SQLite file location
- `CORS_ORIGINS` (default `http://localhost:5173,http://localhost:4173,http://127.0.0.1:5173`) — comma-separated origins allowed to call the API
- `VITE_API_URL` (frontend, default `http://localhost:8080`) — API base URL; set at build/dev time, e.g. `VITE_API_URL=http://192.168.1.10:8080 npm run dev`

AI credentials will be added to the Go backend configuration when an AI provider is selected; they must never be placed in frontend code.
