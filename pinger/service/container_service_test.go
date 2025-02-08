package service_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/KazikovAP/docker_monitoring/pinger/service"
	docker "github.com/fsouza/go-dockerclient"
)

const (
	containersJSONPath      = "/containers/json"
	containerABC123JSONPath = "/containers/abc123/json"
	containerDEF456JSONPath = "/containers/def456/json"
)

func setupFakeDockerServer() *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case containersJSONPath:
			containers := []map[string]interface{}{
				{
					"ID":      "abc123",
					"Names":   []string{"/container1"},
					"Created": float64(1680000000),
					"Status":  "Up 5 minutes",
				},
				{
					"ID":      "def456",
					"Names":   []string{"/container2"},
					"Created": float64(1680001000),
					"Status":  "Exited",
				},
			}

			w.Header().Set("Content-Type", "application/json")

			if err := json.NewEncoder(w).Encode(containers); err != nil {
				panic("failed to encode containers: " + err.Error())
			}
		case containerABC123JSONPath:
			response := map[string]interface{}{
				"NetworkSettings": map[string]interface{}{
					"Networks": map[string]interface{}{
						"bridge": map[string]interface{}{
							"IPAddress": "192.168.1.10",
						},
					},
				},
			}

			w.Header().Set("Content-Type", "application/json")

			if err := json.NewEncoder(w).Encode(response); err != nil {
				panic("failed to encode container abc123 response: " + err.Error())
			}
		case containerDEF456JSONPath:
			response := map[string]interface{}{
				"NetworkSettings": map[string]interface{}{
					"Networks": map[string]interface{}{},
				},
			}

			w.Header().Set("Content-Type", "application/json")

			if err := json.NewEncoder(w).Encode(response); err != nil {
				panic("failed to encode container def456 response: " + err.Error())
			}
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))

	return server
}

func TestGetContainers_Success(t *testing.T) {
	server := setupFakeDockerServer()
	defer server.Close()

	client, err := docker.NewClient(server.URL)
	if err != nil {
		t.Fatalf("failed to create docker client: %v", err)
	}

	svc := service.NewDockerContainerServiceWithClient(client)

	containers, err := svc.GetContainers()
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if len(containers) != 1 {
		t.Fatalf("expected 1 container, got %d", len(containers))
	}

	c := containers[0]
	if c.ID != "abc123" {
		t.Errorf("expected container ID 'abc123', got '%s'", c.ID)
	}

	if c.Name != "/container1" {
		t.Errorf("expected container name '/container1', got '%s'", c.Name)
	}

	if c.IPAddress != "192.168.1.10" {
		t.Errorf("expected IPAddress '192.168.1.10', got '%s'", c.IPAddress)
	}

	if c.Status != "Up 5 minutes" {
		t.Errorf("expected Status 'Up 5 minutes', got '%s'", c.Status)
	}

	expectedCreated := time.Unix(1680000000, 0)
	if !c.CreatedAt.Equal(expectedCreated) {
		t.Errorf("expected CreatedAt %v, got %v", expectedCreated, c.CreatedAt)
	}
}

func TestGetContainers_ListError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == containersJSONPath {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client, err := docker.NewClient(server.URL)
	if err != nil {
		t.Fatalf("failed to create docker client: %v", err)
	}

	svc := service.NewDockerContainerServiceWithClient(client)

	_, err = svc.GetContainers()
	if err == nil {
		t.Error("expected error when ListContainers fails, got nil")
	}
}

func TestGetContainers_InspectError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == containersJSONPath {
			containers := []map[string]interface{}{
				{"ID": "abc123", "Names": []string{"/container1"}, "Created": float64(1680000000), "Status": "Up 5 minutes"},
			}

			w.Header().Set("Content-Type", "application/json")

			if err := json.NewEncoder(w).Encode(containers); err != nil {
				panic("failed to encode containers: " + err.Error())
			}

			return
		}

		if r.URL.Path == containerABC123JSONPath {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client, err := docker.NewClient(server.URL)
	if err != nil {
		t.Fatalf("failed to create docker client: %v", err)
	}

	svc := service.NewDockerContainerServiceWithClient(client)

	_, err = svc.GetContainers()
	if err == nil {
		t.Error("expected error from InspectContainerWithOptions, got nil")
	}
}
