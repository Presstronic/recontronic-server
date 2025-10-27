package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/presstronic/recontronic-server/internal/middleware"
	"github.com/presstronic/recontronic-server/internal/models"
	"github.com/presstronic/recontronic-server/internal/repository"
)

// AnomalyHandler handles anomaly query endpoints
type AnomalyHandler struct {
	anomalyRepo *repository.AnomalyRepository
	validator   Validator
}

// NewAnomalyHandler creates a new anomaly handler
func NewAnomalyHandler(anomalyRepo *repository.AnomalyRepository, validator Validator) *AnomalyHandler {
	return &AnomalyHandler{
		anomalyRepo: anomalyRepo,
		validator:   validator,
	}
}

// ListAnomalies handles GET /api/v1/anomalies
func (h *AnomalyHandler) ListAnomalies(w http.ResponseWriter, r *http.Request) {
	programIDStr := r.URL.Query().Get("program_id")
	unreviewedOnly := r.URL.Query().Get("unreviewed_only") == "true"
	minPriorityStr := r.URL.Query().Get("min_priority")
	limitStr := r.URL.Query().Get("limit")

	limit := 100 // default
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	minPriority := 0.0
	if minPriorityStr != "" {
		if mp, err := strconv.ParseFloat(minPriorityStr, 64); err == nil {
			minPriority = mp
		}
	}

	var anomalies []models.Anomaly
	var err error

	if programIDStr != "" {
		programID, parseErr := strconv.Atoi(programIDStr)
		if parseErr != nil {
			respondError(w, http.StatusBadRequest, "invalid program_id")
			return
		}

		if unreviewedOnly {
			anomalies, err = h.anomalyRepo.ListUnreviewedByProgram(r.Context(), programID, minPriority, limit)
		} else {
			anomalies, err = h.anomalyRepo.ListByProgram(r.Context(), programID, limit)
		}
	} else {
		// Get top anomalies across all programs
		anomalies, err = h.anomalyRepo.GetTopAnomalies(r.Context(), limit)
	}

	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"anomalies": anomalies,
		"total":     len(anomalies),
	})
}

// GetAnomaly handles GET /api/v1/anomalies/{id}
func (h *AnomalyHandler) GetAnomaly(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid anomaly ID")
		return
	}

	anomaly, err := h.anomalyRepo.GetByID(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, anomaly)
}

// MarkAnomalyReviewed handles POST /api/v1/anomalies/{id}/review
func (h *AnomalyHandler) MarkAnomalyReviewed(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid anomaly ID")
		return
	}

	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req models.CreateAnomalyReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validator.Validate(&req); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.anomalyRepo.MarkAsReviewed(r.Context(), id, user.ID, req.Notes); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"message": "Anomaly marked as reviewed",
	})
}

// GetAnomalyCounts handles GET /api/v1/anomalies/counts
func (h *AnomalyHandler) GetAnomalyCounts(w http.ResponseWriter, r *http.Request) {
	programIDStr := r.URL.Query().Get("program_id")
	if programIDStr == "" {
		respondError(w, http.StatusBadRequest, "program_id is required")
		return
	}

	programID, err := strconv.Atoi(programIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid program_id")
		return
	}

	counts, err := h.anomalyRepo.CountByPriorityRange(r.Context(), programID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, counts)
}
