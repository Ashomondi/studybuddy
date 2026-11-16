package httpapi

import (
	"errors"
	"net/http"
	"strconv"

	"studybuddy/internal/models"
	"studybuddy/internal/service"
	"studybuddy/internal/store"
)

type GoalHandlers struct{ Svc *service.GoalService }

func (h *GoalHandlers) List(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	goals, err := h.Svc.List(status)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "internal", "could not list goals")
		return
	}
	if goals == nil {
		goals = []models.Goal{}
	}
	WriteJSON(w, http.StatusOK, map[string]any{"goals": goals})
}

func (h *GoalHandlers) Create(w http.ResponseWriter, r *http.Request) {
	var in models.GoalInput
	if err := DecodeJSON(r, &in); err != nil {
		WriteError(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}
	g, err := h.Svc.Create(in)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	w.Header().Set("Location", "/api/goals/"+strconv.FormatInt(g.ID, 10))
	WriteJSON(w, http.StatusCreated, g)
}

func (h *GoalHandlers) Get(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	g, err := h.Svc.Get(id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			WriteError(w, http.StatusNotFound, "not_found", "goal not found")
			return
		}
		WriteError(w, http.StatusInternalServerError, "internal", "could not fetch goal")
		return
	}
	WriteJSON(w, http.StatusOK, g)
}

func (h *GoalHandlers) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	var in models.GoalInput
	if err := DecodeJSON(r, &in); err != nil {
		WriteError(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}
	g, err := h.Svc.Update(id, in)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			WriteError(w, http.StatusNotFound, "not_found", "goal not found")
			return
		}
		writeServiceError(w, err)
		return
	}
	WriteJSON(w, http.StatusOK, g)
}

func (h *GoalHandlers) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	if err := h.Svc.Delete(id); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			WriteError(w, http.StatusNotFound, "not_found", "goal not found")
			return
		}
		WriteError(w, http.StatusInternalServerError, "internal", "could not delete goal")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func parseID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		WriteError(w, http.StatusBadRequest, "validation_error", "invalid id")
		return 0, false
	}
	return id, true
}

func writeServiceError(w http.ResponseWriter, err error) {
	if errors.Is(err, store.ErrNotFound) {
		WriteError(w, http.StatusNotFound, "not_found", "not found")
		return
	}
	// validation errors are plain errors from models
	WriteError(w, http.StatusBadRequest, "validation_error", err.Error())
}
