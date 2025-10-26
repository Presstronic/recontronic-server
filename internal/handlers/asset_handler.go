package handlers

import (
	"net/http"
	"strconv"

	"github.com/presstronic/recontronic-server/internal/repository"
)

// AssetHandler handles asset query endpoints
type AssetHandler struct {
	assetRepo *repository.AssetRepository
}

// NewAssetHandler creates a new asset handler
func NewAssetHandler(assetRepo *repository.AssetRepository) *AssetHandler {
	return &AssetHandler{
		assetRepo: assetRepo,
	}
}

// ListAssets handles GET /api/v1/assets
func (h *AssetHandler) ListAssets(w http.ResponseWriter, r *http.Request) {
	programIDStr := r.URL.Query().Get("program_id")
	assetType := r.URL.Query().Get("type")
	liveOnly := r.URL.Query().Get("live_only") == "true"

	if programIDStr == "" {
		respondError(w, http.StatusBadRequest, "program_id is required")
		return
	}

	programID, err := strconv.Atoi(programIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid program_id")
		return
	}

	var assets interface{}

	if liveOnly {
		assets, err = h.assetRepo.GetLiveAssets(r.Context(), programID)
	} else if assetType != "" {
		assets, err = h.assetRepo.GetLatestByProgramAndType(r.Context(), programID, assetType, 1000)
	} else {
		assets, err = h.assetRepo.GetLatestUniqueAssets(r.Context(), programID)
	}

	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"assets": assets,
	})
}

// GetAssetCounts handles GET /api/v1/assets/counts
func (h *AssetHandler) GetAssetCounts(w http.ResponseWriter, r *http.Request) {
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

	counts, err := h.assetRepo.CountAssetsByType(r.Context(), programID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, counts)
}
