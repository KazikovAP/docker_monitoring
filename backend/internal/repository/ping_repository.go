package repository

import (
	"log"

	"github.com/KazikovAP/docker_monitoring/backend/internal/model"
	"gorm.io/gorm"
)

type PingRepository interface {
	GetAllPings() ([]model.Ping, error)
	CreatePing(ping *model.Ping) error
	IsIPExists(ip string) (bool, error)
}

type pingRepository struct {
	db *gorm.DB
}

func NewPingRepository(db *gorm.DB) PingRepository {
	// Автоматическая миграция для создания таблицы, если её нет
	if err := db.AutoMigrate(&model.Ping{}); err != nil {
		log.Fatalf("Error during migration: %v", err)
	}

	return &pingRepository{db: db}
}

func (r *pingRepository) GetAllPings() ([]model.Ping, error) {
	var pings []model.Ping
	err := r.db.Find(&pings).Error

	return pings, err
}

func (r *pingRepository) CreatePing(ping *model.Ping) error {
	return r.db.Create(ping).Error
}

func (r *pingRepository) IsIPExists(ip string) (bool, error) {
	var count int64
	err := r.db.Model(&model.Ping{}).Where("ip_address = ?", ip).Count(&count).Error

	return count > 0, err
}
