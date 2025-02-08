package service_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/KazikovAP/docker_monitoring/pinger/models"
	"github.com/KazikovAP/docker_monitoring/pinger/service"
)

func TestSendPingResult_Success(t *testing.T) {
	var receivedBody []byte

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/ping" {
			t.Errorf("unexpected URL: got %s, want %s", r.URL.Path, "/ping")
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("failed to read request body: %v", err)
		}

		receivedBody = body

		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	apiService := service.NewAPIService(ts.URL)
	pingResult := &models.PingResult{
		ContainerID:     "container123",
		IPAddress:       "192.168.0.1",
		PingTime:        150,
		LastSuccessDate: time.Now(),
		Timestamp:       time.Now(),
	}

	err := apiService.SendPingResult(pingResult)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	expectedJSON, err := json.Marshal(pingResult)
	if err != nil {
		t.Fatalf("error marshaling expected ping result: %v", err)
	}

	var expected, actual map[string]interface{}

	if err := json.Unmarshal(expectedJSON, &expected); err != nil {
		t.Fatalf("error unmarshaling expected JSON: %v", err)
	}

	if err := json.Unmarshal(receivedBody, &actual); err != nil {
		t.Fatalf("error unmarshaling actual JSON: %v", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("received JSON does not match expected.\nExpected: %v\nGot: %v", expected, actual)
	}
}

func TestSendPingResult_HTTPPostError(t *testing.T) {
	apiService := service.NewAPIService("http://invalid.url")
	pingResult := &models.PingResult{
		ContainerID:     "container123",
		IPAddress:       "192.168.0.1",
		PingTime:        150,
		LastSuccessDate: time.Now(),
		Timestamp:       time.Now(),
	}

	err := apiService.SendPingResult(pingResult)
	if err == nil {
		t.Error("expected error due to invalid URL, got nil")
	}
}

func TestSendPingResult_NonOKResponse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	apiService := service.NewAPIService(ts.URL)
	pingResult := &models.PingResult{
		ContainerID:     "container123",
		IPAddress:       "192.168.0.1",
		PingTime:        150,
		LastSuccessDate: time.Now(),
		Timestamp:       time.Now(),
	}

	err := apiService.SendPingResult(pingResult)
	if err != nil {
		t.Errorf("expected no error even with non-OK response, got: %v", err)
	}
}
