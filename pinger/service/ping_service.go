package service

import (
	"net"
	"time"

	"github.com/KazikovAP/docker_monitoring/pinger/models"
)

type PingService interface {
	Ping(ip string) models.PingResult
}

type DefaultPingService struct{}

func NewPingService() *DefaultPingService {
	return &DefaultPingService{}
}

func (p *DefaultPingService) Ping(ip string) models.PingResult {
	start := time.Now()
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(ip, "80"), 2*time.Second)
	pingTime := int64(0)

	var lastSuccessDate time.Time

	if err == nil {
		pingTime = time.Since(start).Milliseconds()
		lastSuccessDate = time.Now()

		conn.Close()
	}

	return models.PingResult{
		IPAddress:       ip,
		PingTime:        pingTime,
		LastSuccessDate: lastSuccessDate,
		Timestamp:       time.Now(),
	}
}
