package service

import (
	"fmt"

	"github.com/KazikovAP/docker_monitoring/backend/internal/model"
	"github.com/KazikovAP/docker_monitoring/backend/internal/repository"
)

type PingService interface {
	GetAllPings() ([]model.Ping, error)
	AddPing(ping *model.Ping) error
}

type pingService struct {
	repo repository.PingRepository
}

func NewPingService(repo repository.PingRepository) PingService {
	return &pingService{repo: repo}
}

func (s *pingService) GetAllPings() ([]model.Ping, error) {
	return s.repo.GetAllPings()
}

func (s *pingService) AddPing(ping *model.Ping) error {
	exists, err := s.repo.IsIPExists(ping.IPAddress)
	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf("IP address %s already exists", ping.IPAddress)
	}

	return s.repo.CreatePing(ping)
}
