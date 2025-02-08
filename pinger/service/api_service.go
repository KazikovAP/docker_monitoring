package service

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/KazikovAP/docker_monitoring/pinger/models"
)

type APIService interface {
	SendPingResult(result *models.PingResult) error
}

type DefaultAPIService struct {
	BaseURL string
}

func NewAPIService(baseURL string) *DefaultAPIService {
	return &DefaultAPIService{BaseURL: baseURL}
}

func (api *DefaultAPIService) SendPingResult(result *models.PingResult) error {
	jsonData, err := json.Marshal(result)
	if err != nil {
		return err
	}

	resp, err := http.Post(api.BaseURL+"/ping", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to send ping result, status code: %d", resp.StatusCode)
	}

	return nil
}
