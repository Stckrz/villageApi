package application

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Stckrz/villageApi/internal/db"
	"github.com/Stckrz/villageApi/internal/httpx"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

// defined struct for config items
type Config struct {
	Port        string
	jwtSecret   string
}

// Our APP item which will have our config, database, router, and server. Eventually, the websocket hub too.
type App struct {
	Cfg    Config
	Db     *gorm.DB
	Router *chi.Mux
	// Hub *ws.Hub
	srv *http.Server
}

func New() (*App, error) {
	//create the cfg object
	cfg := Config{
		Port:        env("APP_PORT", ":8080"),
		jwtSecret:   env("JWT_SECRET", "devsecret"),
	}

	//create the DB connection
	database, err := db.ConnectDb()
	if err != nil {
		return nil, fmt.Errorf("Connect db returned nil without error")
	}

	// hub := ws.NewHub(database)

	//create our app object, and setup the server.
	app := &App{
		Cfg: cfg,
		Db:  database,
		// hub: hub,
	}
	app.Router = httpx.BuildRouter(httpx.RouterDeps{
		DB: app.Db,
	})
	app.srv = &http.Server{
		Addr:         cfg.Port,
		Handler:      app.Router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	return app, nil
}

func (a *App) Start(ctx context.Context) error {
	//create buffered channel to capture errors. Buffer size is 1.
	errCh := make(chan error, 1)

	//start http server in goroutine, so that we can keep listening for shutdown signals or fatal server error.
	go func() {
		//blocks until the server is shut down gracefully, or server error
		log.Printf("listening on %s\n", a.Cfg.Port)
		errCh <- a.srv.ListenAndServe()
	}()

	//graceful shutdown on ctx.done(). This is waiting for either the context to be cancelled, or an error from the errorchannel errCh
	select {
	//this case is when the parent cancels the context, like a ctrlC. Then, it gives the existingn requests 10 seconds to finish.
	//If they do not finish, they are cancelled.
	case <-ctx.Done(): 
		shutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := a.srv.Shutdown(shutCtx); err != nil {
			return err
		}
		return nil
	//This case is if listenandserve fails. We just return the error, so it can be logged.
	case err := <-errCh:
		return err
	}
}

func env(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
func mustEnv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("%s not set", k)
	}
	return v
}


