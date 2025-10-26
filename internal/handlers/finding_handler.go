package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/presstronic/recontronic-server/internal/middleware"
	"github.com/presstronic/recontronic-server/internal/models"
	"github.com/presstronic/recontronic-server/internal/services"
)

// FindingHandler handles finding management endpoints
type FindingHandler struct {
	findingService *services.FindingService
	validator      Validator
}

// NewFindingHandler creates a new finding handler
func NewFindingHandler(findingService *services.FindingService, validator Validator) *FindingHandler {
	return &FindingHandler{
		findingService: findingService,
		validator:      validator,
	}
}

// CreateFinding handles POST /api/v1/findings
func (h *FindingHandler) CreateFinding(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req models.CreateFindingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validator.Validate(&req); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	finding, err := h.findingService.CreateFinding(r.Context(), user.ID, &req)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, finding)
}

// GetFinding handles GET /api/v1/findings/{id}
func (h *FindingHandler) GetFinding(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid finding ID")
		return
	}

	finding, err := h.findingService.GetFinding(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, finding)
}

// ListFindings handles GET /api/v1/findings
func (h *FindingHandler) ListFindings(w http.ResponseWriter, r *http.Request) {
	programIDStr := r.URL.Query().Get("program_id")
	limitStr := r.URL.Query().Get("limit")

	limit := 50 // default
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if programIDStr == "" {
		respondError(w, http.StatusBadRequest, "program_id is required")
		return
	}

	programID, err := strconv.Atoi(programIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid program_id")
		return
	}

	findings, err := h.findingService.ListFindingsByProgram(r.Context(), programID, limit)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"findings": findings,
		"total":    len(findings),
	})
}

// UpdateFinding handles PATCH /api/v1/findings/{id}
func (h *FindingHandler) UpdateFinding(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid finding ID")
		return
	}

	var req models.UpdateFindingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validator.Validate(&req); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	finding, err := h.findingService.UpdateFinding(r.Context(), id, &req)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, finding)
}

// GetBountyStats handles GET /api/v1/findings/stats
func (h *FindingHandler) GetBountyStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.findingService.GetBountyStats(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, stats)
}
