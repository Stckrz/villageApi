package httpx

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/gorm"
)

type RouterDeps struct {
	DB  *gorm.DB
}

func BuildRouter(deps RouterDeps) *chi.Mux{
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{
			"http://127.0.0.1:5173",
			"http://localhost:5173",
			"http://localhost:8080",
		},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Authorization"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

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

	//Swagger
	r.Get("/swagger/*", httpSwagger.WrapHandler)

	return r
}

