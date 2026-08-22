package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go-repair-center/internal/service"
	"go-repair-center/internal/transport/http/dto"
	"go-repair-center/internal/transport/http/middleware"
	"go-repair-center/internal/transport/http/response"
)

type SessionHandler struct {
	sessions *service.SessionService
}

func NewSessionHandler(sessions *service.SessionService) *SessionHandler {
	return &SessionHandler{sessions: sessions}
}

func (h *SessionHandler) Routes(mux interface {
	Get(string, http.HandlerFunc)
	Post(string, http.HandlerFunc)
}) {
	mux.Get("/sessions/me", h.MySessions)
	mux.Get("/sessions/users/{userID}", h.UserSessions)
	mux.Post("/sessions/{id}/revoke", h.RevokeByID)
	mux.Post("/sessions/revoke-all", h.RevokeAll)
	mux.Post("/sessions/cleanup", h.CleanupExpired)
}

func (h *SessionHandler) MySessions(w http.ResponseWriter, r *http.Request) {
	summary, err := h.sessions.Summary(r.Context(), middleware.UserID(r.Context()))
	if err != nil {
		writeError(w, r, err)
		return
	}
	response.Success(w, summary, middleware.RequestID(r.Context()))
}

func (h *SessionHandler) UserSessions(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		writeError(w, r, err)
		return
	}
	summary, err := h.sessions.Summary(r.Context(), userID)
	if err != nil {
		writeError(w, r, err)
		return
	}
	response.Success(w, summary, middleware.RequestID(r.Context()))
}

func (h *SessionHandler) RevokeByID(w http.ResponseWriter, r *http.Request) {
	sessionID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		writeError(w, r, err)
		return
	}
	if err := h.sessions.RevokeByID(r.Context(), sessionID); err != nil {
		writeError(w, r, err)
		return
	}
	response.Success(w, map[string]any{"revoked": true, "session_id": sessionID}, middleware.RequestID(r.Context()))
}

func (h *SessionHandler) RevokeAll(w http.ResponseWriter, r *http.Request) {
	var payload dto.SessionRevokeRequest
	_ = json.NewDecoder(r.Body).Decode(&payload)
	if payload.UserID <= 0 {
		payload.UserID = middleware.UserID(r.Context())
	}
	if err := h.sessions.RevokeByUser(r.Context(), payload.UserID); err != nil {
		writeError(w, r, err)
		return
	}
	response.Success(w, map[string]any{"revoked": true, "user_id": payload.UserID}, middleware.RequestID(r.Context()))
}

func (h *SessionHandler) CleanupExpired(w http.ResponseWriter, r *http.Request) {
	count, err := h.sessions.CleanupExpired(r.Context())
	if err != nil {
		writeError(w, r, err)
		return
	}
	response.Success(w, map[string]any{"cleaned": count}, middleware.RequestID(r.Context()))
}
