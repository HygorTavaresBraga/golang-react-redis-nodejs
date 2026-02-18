package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"support-backend/internal/service"
)

type TicketHandler struct {
	svc *service.TicketService
}

func NewTicketHandler(svc *service.TicketService) *TicketHandler {
	return &TicketHandler{svc: svc}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func decodeJSON(r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(dst)
}

func (h *TicketHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input service.CreateTicketInput
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "request body must be valid JSON")
		return
	}

	t, err := h.svc.Create(r.Context(), input)
	if err != nil {
		if errors.Is(err, service.ErrInvalidInput) {
			writeError(w, http.StatusBadRequest, "invalid_input", "title/description are invalid")
			return
		}
		writeInternalError(w)
		return
	}

	writeJSON(w, http.StatusCreated, t)
}

func (h *TicketHandler) List(w http.ResponseWriter, r *http.Request) {
	tickets, err := h.svc.List(r.Context())
	if err != nil {
		writeInternalError(w)
		return
	}
	writeJSON(w, http.StatusOK, tickets)
}

func (h *TicketHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, ok := extractID(r.URL.Path, "/api/v1/tickets/")
	if !ok {
		writeError(w, http.StatusBadRequest, "bad_request", "invalid ticket id")
		return
	}

	t, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			writeError(w, http.StatusNotFound, "not_found", "ticket not found")
			return
		}
		writeInternalError(w)
		return
	}

	writeJSON(w, http.StatusOK, t)
}

func (h *TicketHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := extractID(r.URL.Path, "/api/v1/tickets/")
	if !ok {
		writeError(w, http.StatusBadRequest, "bad_request", "invalid ticket id")
		return
	}

	var input service.UpdateTicketInput
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", "request body must be valid JSON")
		return
	}

	t, err := h.svc.Update(r.Context(), id, input)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrNotFound):
			writeError(w, http.StatusNotFound, "not_found", "ticket not found")
		case errors.Is(err, service.ErrInvalidInput):
			writeError(w, http.StatusBadRequest, "invalid_input", "title/description are invalid")
		case errors.Is(err, service.ErrInvalidStatus):
			writeError(w, http.StatusBadRequest, "invalid_status", "status must be OPEN, PENDING, or CLOSED")
		default:
			writeInternalError(w)
		}
		return
	}

	writeJSON(w, http.StatusOK, t)
}

func (h *TicketHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := extractID(r.URL.Path, "/api/v1/tickets/")
	if !ok {
		writeError(w, http.StatusBadRequest, "bad_request", "invalid ticket id")
		return
	}

	err := h.svc.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			writeError(w, http.StatusNotFound, "not_found", "ticket not found")
			return
		}
		writeInternalError(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
