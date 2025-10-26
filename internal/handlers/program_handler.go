package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/presstronic/recontronic-server/internal/models"
	"github.com/presstronic/recontronic-server/internal/services"
)

// ProgramHandler handles program management endpoints
type ProgramHandler struct {
	programService *services.ProgramService
	validator      Validator
}

// NewProgramHandler creates a new program handler
func NewProgramHandler(programService *services.ProgramService, validator Validator) *ProgramHandler {
	return &ProgramHandler{
		programService: programService,
		validator:      validator,
	}
}

// CreateProgram handles POST /api/v1/programs
func (h *ProgramHandler) CreateProgram(w http.ResponseWriter, r *http.Request) {
	var req models.CreateProgramRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validator.Validate(&req); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	program, err := h.programService.CreateProgram(r.Context(), &req)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, program)
}

// GetProgram handles GET /api/v1/programs/{id}
func (h *ProgramHandler) GetProgram(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid program ID")
		return
	}

	program, err := h.programService.GetProgram(r.Context(), id)
	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, program)
}

// ListPrograms handles GET /api/v1/programs
func (h *ProgramHandler) ListPrograms(w http.ResponseWriter, r *http.Request) {
	// Check for active_only query parameter
	activeOnly := r.URL.Query().Get("active_only") == "true"

	programs, err := h.programService.ListPrograms(r.Context(), activeOnly)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"programs": programs,
		"total":    len(programs),
	})
}

// UpdateProgram handles PATCH /api/v1/programs/{id}
func (h *ProgramHandler) UpdateProgram(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid program ID")
		return
	}

	var req models.UpdateProgramRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validator.Validate(&req); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	program, err := h.programService.UpdateProgram(r.Context(), id, &req)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, program)
}

// DeleteProgram handles DELETE /api/v1/programs/{id}
func (h *ProgramHandler) DeleteProgram(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid program ID")
		return
	}

	if err := h.programService.DeleteProgram(r.Context(), id); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"message": "Program deleted successfully",
	})
}
