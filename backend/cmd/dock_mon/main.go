package main

import (
	"log"

	"github.com/KazikovAP/docker_monitoring/backend/internal/config"
	"github.com/KazikovAP/docker_monitoring/backend/internal/handler"
	"github.com/KazikovAP/docker_monitoring/backend/internal/repository"
	"github.com/KazikovAP/docker_monitoring/backend/internal/service"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.LoadConfig()

	dsn := "host=" + cfg.DBHost +
		" user=" + cfg.DBUser +
		" password=" + cfg.DBPassword +
		" dbname=" + cfg.DBName +
		" port=" + cfg.DBPort +
		" sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get DB instance: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Database not reachable: %v", err)
	}

	repo := repository.NewPingRepository(db)
	pingService := service.NewPingService(repo)
	pingHandler := handler.NewPingHandler(pingService)

	r := handler.NewRouter(pingHandler)

	log.Printf("Server running on port %s", cfg.ServerPort)

	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
