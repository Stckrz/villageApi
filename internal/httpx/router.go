package httpx

import (
	"net/http"

	"github.com/Stckrz/villageApi/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"

	"gorm.io/gorm"
)

type RouterDeps struct {
	DB *gorm.DB
}

func BuildRouter(deps RouterDeps) *chi.Mux {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"http://127.0.0.1:5173",
			"http://localhost:5173",
			"http://localhost:8080",
		},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Authorization"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	buildingService := services.NewBuildingService(deps.DB)
	taskService := services.NewTaskService(deps.DB)

	buildings := NewBuildingHandler(deps.DB, buildingService)
	tasks := NewTaskHandler(deps.DB, taskService)

	// Health Check godoc
	// @Summary Health Check
	// @Tags health
	// @Produce text/plain
	// @Success 200 {string} string "ok"
	// @Router /api/healthz [get]
	r.Get("/api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	// Building Endpoints
	r.Get("/api/buildings", buildings.ListBuildings)
	r.Get("/api/buildings/{id}", buildings.GetBuilding)
	r.Post("/api/buildings", buildings.CreateBuilding)
	r.Delete("/api/buildings/{id}", buildings.DeleteBuilding)
	r.Put("/api/buildings/{id}", buildings.UpdateBuilding)

	//Task Endpoints
	r.Get("/api/tasks", tasks.ListTasks)
	r.Post("/api/tasks", tasks.CreateTask)
	r.Delete("/api/tasks/{id}", tasks.DeleteTask)
	r.Put("/api/tasks/{id}", tasks.UpdateTask)

	//Swagger
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	return r
}
