package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	application "github.com/Stckrz/villageApi/internal/app"
	"github.com/joho/godotenv"
)

// @title Village Api
// @version 1.0
// @description Service for managing Village UI Application
// @BasePath /api
func main() {
    godotenv.Load(".env") // loads env vars
	app, err := application.New()
	if err != nil {log.Fatal(err)}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := app.Start(ctx); err != nil {
		log.Println("server exited:", err)
	}
}
