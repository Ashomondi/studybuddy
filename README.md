# StudyBuddy AI

StudyBuddy AI is a learning companion built with a Go REST API, SQLite, and a React/Vite frontend. Development is organized as small, runnable increments.

## Current features

StudyBuddy currently supports learning goals with study preferences, manually authored learning plans, scheduled activities, weekly schedule views, and study session tracking. Plans are user-authored; AI generation is not implemented.

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

Configuration is read from environment variables. See `.env.example`.

- `API_ADDR` (default `:8080`) — API listen address
- `DATABASE_PATH` (default `studybuddy.db`) — SQLite file location
- `CORS_ORIGINS` (default `http://localhost:5173,http://localhost:4173,http://127.0.0.1:5173`) — comma-separated origins allowed to call the API
- `VITE_API_URL` (frontend, default `http://localhost:8080`) — API base URL; set at build/dev time, e.g. `VITE_API_URL=http://192.168.1.10:8080 npm run dev`

AI credentials will be added to the Go backend configuration when an AI provider is selected; they must never be placed in frontend code.

## API

Base URL: `http://localhost:8080/api`. JSON only. Errors: `{"error":{"code":...,"message":...}}`.

- `GET /api/health`
- `GET /api/goals` (`?status=active`)
- `POST /api/goals` `{title, subject?, description?, status?, target_date?}`
- `GET /api/goals/{id}`
- `PUT /api/goals/{id}`
- `DELETE /api/goals/{id}` (204)
- `GET /api/plans` (`?goal_id=`)
- `GET /api/plans/{id}`
- `POST /api/plans` `{goal_id, title, description?}`
- `PUT /api/plans/{id}` `{goal_id, title, description?}`
- `DELETE /api/plans/{id}` (204; activities are removed with the plan)
- `GET /api/plans/{id}/topics` (returns `activities`)
- `POST /api/activities` `{plan_id, title, description?, position?, estimated_minutes?, scheduled_date?, status?}`
- `PUT /api/activities/{id}` (same fields)
- `DELETE /api/activities/{id}` (204)
- `GET /api/schedule?from=YYYY-MM-DD&to=YYYY-MM-DD` (maximum 32-day range)
- `POST /api/sessions` `{goal_id, plan_topic_id?, started_at?}`
- `GET /api/sessions` (`?goal_id=`, `?status=`)
- `GET /api/sessions/{id}`

Schema is auto-migrated on startup (see `internal/migrations`). SQLite file is created if missing.

Goal create/update also accepts `current_level`, `weekly_study_minutes`, and `preferred_study_days` (weekday numbers, Sunday `0` through Saturday `6`). Activity statuses are `pending`, `in_progress`, and `completed`; completed activities cannot be reopened.
