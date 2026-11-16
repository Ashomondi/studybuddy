package httpapi

import (
	"net/http"

	"studybuddy/internal/service"
)

type RouterDeps struct {
	Goals    *service.GoalService
	Plans    *service.PlanService
	Sessions *service.SessionService
	DBPing   func() error
}

func NewRouter(origins string, deps RouterDeps) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		if deps.DBPing != nil {
			if err := deps.DBPing(); err != nil {
				WriteError(w, http.StatusServiceUnavailable, "unavailable", "database disconnected")
				return
			}
		}
		WriteJSON(w, http.StatusOK, map[string]string{"status": "ok", "database": "connected"})
	})

	gh := &GoalHandlers{Svc: deps.Goals}
	mux.HandleFunc("GET /api/goals", gh.List)
	mux.HandleFunc("POST /api/goals", gh.Create)
	mux.HandleFunc("GET /api/goals/{id}", gh.Get)
	mux.HandleFunc("PUT /api/goals/{id}", gh.Update)
	mux.HandleFunc("DELETE /api/goals/{id}", gh.Delete)

	ph := &PlanHandlers{Svc: deps.Plans}
	mux.HandleFunc("GET /api/plans", ph.List)
	mux.HandleFunc("GET /api/plans/{id}", ph.Get)
	mux.HandleFunc("POST /api/plans", ph.Create)
	mux.HandleFunc("PUT /api/plans/{id}", ph.Update)
	mux.HandleFunc("DELETE /api/plans/{id}", ph.Delete)
	mux.HandleFunc("GET /api/plans/{id}/topics", ph.Topics)
	mux.HandleFunc("GET /api/schedule", ph.Schedule)
	mux.HandleFunc("POST /api/activities", ph.CreateTopic)
	mux.HandleFunc("PUT /api/activities/{id}", ph.UpdateTopic)
	mux.HandleFunc("DELETE /api/activities/{id}", ph.DeleteTopic)

	sh := &SessionHandlers{Svc: deps.Sessions}
	mux.HandleFunc("POST /api/sessions", sh.Create)
	mux.HandleFunc("GET /api/sessions", sh.List)
	mux.HandleFunc("GET /api/sessions/{id}", sh.Get)

	handler := Chain(mux, Recovery, Logging, CORS(origins))
	return handler
}
