package service_test

import (
	"net"
	"testing"

	"github.com/KazikovAP/docker_monitoring/pinger/service"
)

func TestPingService_Unreachable(t *testing.T) {
	pingService := service.NewPingService()
	result := pingService.Ping("192.0.2.1")

	if result.PingTime != 0 {
		t.Errorf("expected pingTime 0 for unreachable IP, got %d", result.PingTime)
	}

	if !result.LastSuccessDate.IsZero() {
		t.Errorf("expected LastSuccessDate to be zero for unreachable IP, got %v", result.LastSuccessDate)
	}
}

func TestPingService_InvalidIP(t *testing.T) {
	pingService := service.NewPingService()
	result := pingService.Ping("invalid-ip")

	if result.PingTime != 0 {
		t.Errorf("expected pingTime 0 for invalid IP, got %d", result.PingTime)
	}

	if !result.LastSuccessDate.IsZero() {
		t.Errorf("expected LastSuccessDate to be zero for invalid IP, got %v", result.LastSuccessDate)
	}
}

func TestPingService_Success(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:80")
	if err != nil {
		t.Skipf("Skipping success test: unable to bind to port 80: %v", err)
	}
	defer ln.Close()

	go func() {
		conn, err := ln.Accept()
		if err == nil {
			conn.Close()
		}
	}()

	pingService := service.NewPingService()
	result := pingService.Ping("127.0.0.1")

	if result.PingTime <= 0 {
		t.Errorf("expected pingTime > 0 for reachable IP, got %d", result.PingTime)
	}

	if result.LastSuccessDate.IsZero() {
		t.Errorf("expected LastSuccessDate to be set for reachable IP, got zero time")
	}
}
