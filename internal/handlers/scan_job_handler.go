package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/presstronic/recontronic-server/internal/models"
	"github.com/presstronic/recontronic-server/internal/services"
)

// ScanJobHandler handles scan job endpoints
type ScanJobHandler struct {
	scanJobService *services.ScanJobService
	validator      Validator
}

// NewScanJobHandler creates a new scan job handler
func NewScanJobHandler(scanJobService *services.ScanJobService, validator Validator) *ScanJobHandler {
	return &ScanJobHandler{
		scanJobService: scanJobService,
		validator:      validator,
	}
}

// CreateScanJob handles POST /api/v1/scans
func (h *ScanJobHandler) CreateScanJob(w http.ResponseWriter, r *http.Request) {
	var req models.CreateScanJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validator.Validate(&req); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	scanJob, err := h.scanJobService.CreateScanJob(r.Context(), &req)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusAccepted, scanJob)
}

// GetScanJob handles GET /api/v1/scans/{id}
func (h *ScanJobHandler) GetScanJob(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid scan job ID")
		return
	}

	scanJob, err := h.scanJobService.GetScanJob(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, scanJob)
}

// ListScanJobs handles GET /api/v1/scans
func (h *ScanJobHandler) ListScanJobs(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	programIDStr := r.URL.Query().Get("program_id")
	status := r.URL.Query().Get("status")
	limitStr := r.URL.Query().Get("limit")

	limit := 50 // default
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	var scanJobs []models.ScanJob
	var err error

	// Filter by program ID if provided
	if programIDStr != "" {
		programID, parseErr := strconv.Atoi(programIDStr)
		if parseErr != nil {
			respondError(w, http.StatusBadRequest, "invalid program_id")
			return
		}
		scanJobs, err = h.scanJobService.ListScanJobsByProgram(r.Context(), programID, limit)
	} else if status != "" {
		// Filter by status
		scanJobs, err = h.scanJobService.ListScanJobsByStatus(r.Context(), status, limit)
	} else {
		// Get pending jobs by default
		scanJobs, err = h.scanJobService.GetPendingJobs(r.Context(), limit)
	}

	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"scans": scanJobs,
		"total": len(scanJobs),
	})
}
