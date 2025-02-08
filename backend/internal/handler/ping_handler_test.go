package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/KazikovAP/docker_monitoring/backend/internal/handler"
	"github.com/KazikovAP/docker_monitoring/backend/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPingService struct {
	mock.Mock
}

func (m *MockPingService) GetAllPings() ([]model.Ping, error) {
	args := m.Called()
	return args.Get(0).([]model.Ping), args.Error(1)
}

func (m *MockPingService) AddPing(ping *model.Ping) error {
	args := m.Called(ping)
	return args.Error(0)
}

func setupRouter(pingHandler *handler.PingHandler) *gin.Engine {
	r := gin.Default()
	r.GET("/pings", pingHandler.GetAllPings)
	r.POST("/pings", pingHandler.AddPing)

	return r
}

func TestPingHandler_GetAllPings(t *testing.T) {
	mockService := new(MockPingService)
	pingHandler := handler.NewPingHandler(mockService)
	r := setupRouter(pingHandler)

	fixedTime := time.Date(2024, 2, 4, 19, 35, 52, 0, time.UTC)

	expectedPings := []model.Ping{
		{ID: 1, IPAddress: "192.168.0.1", PingTime: 20, LastSuccessDate: fixedTime},
	}
	mockService.On("GetAllPings").Return(expectedPings, nil)

	req, _ := http.NewRequest(http.MethodGet, "/pings", http.NoBody)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var actualPings []model.Ping
	_ = json.Unmarshal(w.Body.Bytes(), &actualPings)

	assert.Equal(t, expectedPings, actualPings)
	mockService.AssertExpectations(t)
}

func TestPingHandler_AddPing_Success(t *testing.T) {
	mockService := new(MockPingService)
	pingHandler := handler.NewPingHandler(mockService)
	r := setupRouter(pingHandler)

	newPing := model.Ping{IPAddress: "192.168.0.2", PingTime: 15}
	mockService.On("AddPing", &newPing).Return(nil)

	jsonData, _ := json.Marshal(newPing)
	req, _ := http.NewRequest(http.MethodPost, "/pings", bytes.NewBuffer(jsonData))

	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Ping added successfully")
	mockService.AssertExpectations(t)
}

func TestPingHandler_AddPing_Conflict(t *testing.T) {
	mockService := new(MockPingService)
	pingHandler := handler.NewPingHandler(mockService)
	r := setupRouter(pingHandler)

	conflictPing := model.Ping{IPAddress: "192.168.0.3"}
	err := errors.New("IP address 192.168.0.3 already exists")
	mockService.On("AddPing", &conflictPing).Return(err)

	jsonData, _ := json.Marshal(conflictPing)
	req, _ := http.NewRequest(http.MethodPost, "/pings", bytes.NewBuffer(jsonData))

	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusConflict, w.Code)
	assert.Contains(t, w.Body.String(), err.Error())
	mockService.AssertExpectations(t)
}

func TestPingHandler_AddPing_BadRequest(t *testing.T) {
	mockService := new(MockPingService)
	pingHandler := handler.NewPingHandler(mockService)
	r := setupRouter(pingHandler)

	invalidJSON := []byte(`{"invalid": "data"}`)

	req, _ := http.NewRequest(http.MethodPost, "/pings", bytes.NewBuffer(invalidJSON))

	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Key: 'Ping.IPAddress' Error")
	mockService.AssertExpectations(t)
}
