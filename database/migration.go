package database

import (
	"log"

	"github.com/shubhamgupta32873/todo-app/config"
	"github.com/shubhamgupta32873/todo-app/models"
)

// RunMigrations runs all database migrations
func RunMigrations() {
	db := config.GetDB()

	// Auto-migrate will create tables if they don't exist
	err := db.AutoMigrate(
		&models.User{},
		&models.Task{},
		&models.Reminder{},
	)

	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("✅ Database migrations completed successfully")
}
