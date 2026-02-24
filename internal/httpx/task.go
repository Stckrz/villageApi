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

type TaskHandler struct {
	service services.TaskService
	db      *gorm.DB
}

func NewTaskHandler(db *gorm.DB, service services.TaskService) *TaskHandler {
	return &TaskHandler{
		db:      db,
		service: service,
	}
}

type CreateTaskRequest struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	BuildingId  uint            `json:"building_id"`
	IsCompleted bool            `json:"is_completed"`
}

type UpdateTaskRequest struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	BuildingId  uint            `json:"building_id"`
	IsCompleted bool            `json:"is_completed"`
}

// GetTasks godoc
// @Summary Get tasks
// @Tags tasks
// @Produce json
// @Success 200 {array} models.Task
// @Router /tasks [get]
func (h *TaskHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.ListTasks()
	if err != nil {
		http.Error(w, "failed to fetch users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// @CreateTask godoc
// @Summary Create new task
// @Tags tasks
// @Produce application/json
// @Param request body CreateTaskRequest true "Create task payload"
// @Success 200 {object} map[string]string
// @Failure 500 {string} string "Internal Service Error"
// @Router /tasks [post]
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var body CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	task := models.Task{
		Name:        body.Name,
		Description: body.Description,
		BuildingId:  body.BuildingId,
		IsCompleted: body.IsCompleted,
	}

	task, err := h.service.CreateTask(task)
	if err != nil {
		switch err {
		default:
			http.Error(w, "failed to create task", http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(task)
	return
}

// @DeleteTask godoc
// @Summary Delete a task
// @Tags tasks
// @Produce application/json
// @Param id path int true "Task ID"
// @Success 200 {object} map[string]string
// @Failure 500 {string} string "Internal Service Error"
// @Router /tasks/{id} [delete]
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	idInt, err := strconv.Atoi(idParam)
	if err != nil || idInt <= 0 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteTask(uint(idInt)); err != nil {
		http.Error(w, "failed to delete task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// @UpdateTask godoc
// @Summary Update a task
// @Tags tasks
// @Produce application/json
// @Param id path int true "Task ID"
// @Param request body UpdateTaskRequest true "Update task payload"
// @Success 200 {object} map[string]string
// @Failure 500 {string} string "Internal Service Error"
// @Router /tasks/{id} [put]
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	idInt, err := strconv.Atoi(idParam)
	if err != nil || idInt <= 0 {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var body UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	// categories := make([]models.BuildingCategory, 0, len(body.Categories))

	// for _, text := range body.Categories {
	// 	categories = append(categories, models.BuildingCategory{
	// 		Text: text,
	// 	})
	// }

	task := models.Task{
		Name:          body.Name,
		Description:   body.Description,
		BuildingId:    body.BuildingId,
		IsCompleted:     body.IsCompleted,
	}

	if err := h.service.UpdateTask(task, uint(idInt)); err != nil {
		http.Error(w, "failed to update task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
