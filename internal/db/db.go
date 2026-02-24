package db

import (
	"fmt"
	"log"
	"os"

	"github.com/Stckrz/villageApi/internal/db/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ConnectDb() (*gorm.DB, error) {
	//try to get dbPath from env, otherwise next to main.go
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		// dbPath = "/data/app.db"
		dbPath = "data/app.db"
	}

	dsn := fmt.Sprintf(
		"%s?_pragma=journal_mode(WAL)&_pragma=busy_timeout=5000",
		dbPath,
	)
	// dsn := "app.db?_pragma=journal_mode(WAL)&_pragma=busy_timeout=5000"

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
		return nil, err
	}
	if err := db.AutoMigrate(
		&models.Building{},
		&models.BuildingCategory{},
		&models.Task{},
	); err != nil {
		return nil, err
	}
	// if err := db.Exec(`
	// 	CREATE UNIQUE INDEX IF NOT EXISTS unique_open_round_per_dealer
	// 	ON betting_rounds (dealer_id)
	// 	WHERE status = 'open'
	// `).Error; err != nil {
	// 		return nil, err
	// }
	if os.Getenv("ENVIRONMENT") == "dev" {
	}

	return db, nil
}

