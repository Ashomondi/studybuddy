package httpapi

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"studybuddy/internal/models"
	"studybuddy/internal/service"
	"studybuddy/internal/store"
)

type PlanHandlers struct{ Svc *service.PlanService }

func (h *PlanHandlers) List(w http.ResponseWriter, r *http.Request) {
	var goalID *int64
	if v := r.URL.Query().Get("goal_id"); v != "" {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil || id <= 0 {
			WriteError(w, http.StatusBadRequest, "validation_error", "invalid goal_id")
			return
		}
		goalID = &id
	}
	plans, err := h.Svc.List(goalID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			WriteError(w, http.StatusNotFound, "not_found", "goal not found")
			return
		}
		WriteError(w, http.StatusInternalServerError, "internal", "could not list plans")
		return
	}
	if plans == nil {
		plans = []models.StudyPlan{}
	}
	WriteJSON(w, http.StatusOK, map[string]any{"plans": plans})
}

func (h *PlanHandlers) Create(w http.ResponseWriter, r *http.Request) {
	var in models.PlanInput
	if err := DecodeJSON(r, &in); err != nil {
		WriteError(w, 400, "validation_error", "invalid JSON body")
		return
	}
	p, err := h.Svc.Create(in)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	WriteJSON(w, 201, p)
}
func (h *PlanHandlers) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	var in models.PlanInput
	if err := DecodeJSON(r, &in); err != nil {
		WriteError(w, 400, "validation_error", "invalid JSON body")
		return
	}
	p, err := h.Svc.Update(id, in)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			WriteError(w, 404, "not_found", "plan or goal not found")
		} else {
			writeServiceError(w, err)
		}
		return
	}
	WriteJSON(w, 200, p)
}
func (h *PlanHandlers) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	if err := h.Svc.Delete(id); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			WriteError(w, 404, "not_found", "plan not found")
		} else {
			WriteError(w, 500, "internal", "could not delete plan")
		}
		return
	}
	w.WriteHeader(204)
}
func (h *PlanHandlers) Topics(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	topics, err := h.Svc.Topics(id)
	if err != nil {
		WriteError(w, 404, "not_found", "plan not found")
		return
	}
	WriteJSON(w, 200, map[string]any{"activities": topics})
}
func (h *PlanHandlers) CreateTopic(w http.ResponseWriter, r *http.Request) {
	var in models.TopicInput
	if err := DecodeJSON(r, &in); err != nil {
		WriteError(w, 400, "validation_error", "invalid JSON body")
		return
	}
	t, err := h.Svc.CreateTopic(in)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			WriteError(w, 404, "not_found", "plan not found")
		} else {
			writeServiceError(w, err)
		}
		return
	}
	WriteJSON(w, 201, t)
}
func (h *PlanHandlers) UpdateTopic(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	var in models.TopicInput
	if err := DecodeJSON(r, &in); err != nil {
		WriteError(w, 400, "validation_error", "invalid JSON body")
		return
	}
	t, err := h.Svc.UpdateTopic(id, in)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			WriteError(w, 404, "not_found", "activity or plan not found")
		} else {
			writeServiceError(w, err)
		}
		return
	}
	WriteJSON(w, 200, t)
}
func (h *PlanHandlers) DeleteTopic(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	if err := h.Svc.DeleteTopic(id); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			WriteError(w, 404, "not_found", "activity not found")
		} else {
			WriteError(w, 500, "internal", "could not delete activity")
		}
		return
	}
	w.WriteHeader(204)
}
func (h *PlanHandlers) Schedule(w http.ResponseWriter, r *http.Request) {
	from, to := r.URL.Query().Get("from"), r.URL.Query().Get("to")
	if from == "" || to == "" {
		WriteError(w, 400, "validation_error", "from and to dates are required")
		return
	}
	a, e1 := time.Parse("2006-01-02", from)
	b, e2 := time.Parse("2006-01-02", to)
	if e1 != nil || e2 != nil || b.Before(a) || b.Sub(a) > 31*24*time.Hour {
		WriteError(w, 400, "validation_error", "schedule range must be YYYY-MM-DD and no longer than 32 days")
		return
	}
	items, err := h.Svc.Schedule(from, to)
	if err != nil {
		WriteError(w, 500, "internal", "could not load schedule")
		return
	}
	WriteJSON(w, 200, map[string]any{"activities": items})
}

func (h *PlanHandlers) Get(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	p, err := h.Svc.Get(id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			WriteError(w, http.StatusNotFound, "not_found", "plan not found")
			return
		}
		WriteError(w, http.StatusInternalServerError, "internal", "could not fetch plan")
		return
	}
	WriteJSON(w, http.StatusOK, p)
}
