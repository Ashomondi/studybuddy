package httpapi

import (
	"errors"
	"net/http"
	"strconv"

	"studybuddy/internal/models"
	"studybuddy/internal/service"
	"studybuddy/internal/store"
)

type SessionHandlers struct{ Svc *service.SessionService }

func (h *SessionHandlers) Create(w http.ResponseWriter, r *http.Request) {
	var in models.SessionInput
	if err := DecodeJSON(r, &in); err != nil {
		WriteError(w, http.StatusBadRequest, "validation_error", "invalid JSON body")
		return
	}
	s, err := h.Svc.Create(in)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			WriteError(w, http.StatusNotFound, "not_found", "goal not found")
			return
		}
		WriteError(w, http.StatusBadRequest, "validation_error", err.Error())
		return
	}
	WriteJSON(w, http.StatusCreated, s)
}

func (h *SessionHandlers) Get(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	s, err := h.Svc.Get(id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			WriteError(w, http.StatusNotFound, "not_found", "session not found")
			return
		}
		WriteError(w, http.StatusInternalServerError, "internal", "could not fetch session")
		return
	}
	WriteJSON(w, http.StatusOK, s)
}

func (h *SessionHandlers) List(w http.ResponseWriter, r *http.Request) {
	var goalID *int64
	if v := r.URL.Query().Get("goal_id"); v != "" {
		id, err := strconv.ParseInt(v, 10, 64)
		if err != nil || id <= 0 {
			WriteError(w, http.StatusBadRequest, "validation_error", "invalid goal_id")
			return
		}
		goalID = &id
	}
	status := r.URL.Query().Get("status")
	sessions, err := h.Svc.List(goalID, status)
	if err != nil {
		WriteError(w, http.StatusInternalServerError, "internal", "could not list sessions")
		return
	}
	if sessions == nil {
		sessions = []models.StudySession{}
	}
	WriteJSON(w, http.StatusOK, map[string]any{"sessions": sessions})
}
