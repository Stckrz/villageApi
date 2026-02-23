package httpx

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Stckrz/villageApi/internal/db/models"
	"github.com/Stckrz/villageApi/internal/services"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type BuildingHandler struct {
	service services.BuildingService
	db      *gorm.DB
}

func NewBuildingHandler(db *gorm.DB, service services.BuildingService) *BuildingHandler {
	return &BuildingHandler{
		db:      db,
		service: service,
	}
}

type CreateBuildingRequest struct {
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	Categories    []string `json:"categories"`
	ThumbnailPath string   `json:"thumbnailPath"`
	ImagePath     string   `json:"imagePath"`
}

type UpdateBuildingRequest struct {
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	Categories    []string `json:"categories"`
	ThumbnailPath string   `json:"thumbnailPath"`
	ImagePath     string   `json:"imagePath"`
}

// GetBuildingById godoc
// @Summary Get building by id
// @Tags buildings
// @Produce json
// @Param id path int true "Building ID"
// @Success 200 {Building} models.Building
// @Router /buildings/{id} [get]
func (h *BuildingHandler) GetBuilding(w http.ResponseWriter, r *http.Request) {

	idParam := chi.URLParam(r, "id")

	idInt, err := strconv.Atoi(idParam)
	if err != nil || idInt <= 0 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	users, err := h.service.GetBuildingByID(uint(idInt))
	if err != nil {
		http.Error(w, "failed to fetch building", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// GetBuildings godoc
// @Summary Get buildings
// @Tags buildings
// @Produce json
// @Success 200 {array} models.Building
// @Router /buildings [get]
func (h *BuildingHandler) ListBuildings(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.ListBuildings()
	if err != nil {
		http.Error(w, "failed to fetch users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// @CreateBuilding godoc
// @Summary Create new building
// @Tags buildings
// @Produce application/json
// @Param request body CreateBuildingRequest true "Create building payload"
// @Success 200 {object} map[string]string
// @Failure 500 {string} string "Internal Service Error"
// @Router /buildings [post]
func (h *BuildingHandler) CreateBuilding(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var body CreateBuildingRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	categories := make([]models.BuildingCategory, 0, len(body.Categories))

	for _, text := range body.Categories {
		categories = append(categories, models.BuildingCategory{
			Text: text,
		})
	}
	building := models.Building{
		Name:          body.Name,
		Description:   body.Description,
		Categories:    categories,
		ImagePath:     body.ImagePath,
		ThumbnailPath: body.ThumbnailPath,
	}

	building, err := h.service.CreateBuilding(building)
	if err != nil {
		switch err {
		default:
			http.Error(w, "failed to create building", http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(building)
	return
}

// @DeleteBuilding godoc
// @Summary Delete a building
// @Tags buildings
// @Produce application/json
// @Param id path int true "Building ID"
// @Success 200 {object} map[string]string
// @Failure 500 {string} string "Internal Service Error"
// @Router /buildings/{id} [delete]
func (h *BuildingHandler) DeleteBuilding(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	idInt, err := strconv.Atoi(idParam)
	if err != nil || idInt <= 0 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteBuilding(uint(idInt)); err != nil {
		http.Error(w, "failed to delete building", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @UpdateBuilding godoc
// @Summary Update a building
// @Tags buildings
// @Produce application/json
// @Param id path int true "Building ID"
// @Param request body UpdateBuildingRequest true "Update building payload"
// @Success 200 {object} map[string]string
// @Failure 500 {string} string "Internal Service Error"
// @Router /buildings/{id} [put]
func (h *BuildingHandler) UpdateBuilding(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	idInt, err := strconv.Atoi(idParam)
	if err != nil || idInt <= 0 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var body UpdateBuildingRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	categories := make([]models.BuildingCategory, 0, len(body.Categories))

	for _, text := range body.Categories {
		categories = append(categories, models.BuildingCategory{
			Text: text,
		})
	}

	building := models.Building{
		Name:          body.Name,
		Description:   body.Description,
		Categories:    categories,
		ImagePath:     body.ImagePath,
		ThumbnailPath: body.ThumbnailPath,
	}

	if err := h.service.UpdateBuilding(building, uint(idInt)); err != nil {
		http.Error(w, "failed to delete building", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
