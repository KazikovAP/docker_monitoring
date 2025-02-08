package service_test

import (
	"errors"
	"testing"

	"github.com/KazikovAP/docker_monitoring/backend/internal/model"
	"github.com/KazikovAP/docker_monitoring/backend/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPingRepository struct {
	mock.Mock
}

func (m *MockPingRepository) GetAllPings() ([]model.Ping, error) {
	args := m.Called()
	return args.Get(0).([]model.Ping), args.Error(1)
}

func (m *MockPingRepository) CreatePing(ping *model.Ping) error {
	args := m.Called(ping)
	return args.Error(0)
}

func (m *MockPingRepository) IsIPExists(ip string) (bool, error) {
	args := m.Called(ip)
	return args.Bool(0), args.Error(1)
}

func TestPingService_GetAllPings(t *testing.T) {
	repo := new(MockPingRepository)
	servicePing := service.NewPingService(repo)

	expectedPings := []model.Ping{
		{ID: 1, IPAddress: "192.168.0.1", PingTime: 50},
		{ID: 2, IPAddress: "192.168.0.2", PingTime: 100},
	}

	repo.On("GetAllPings").Return(expectedPings, nil)

	result, err := servicePing.GetAllPings()

	assert.NoError(t, err)
	assert.Equal(t, expectedPings, result)
	repo.AssertExpectations(t)
}

func TestPingService_AddPing_Success(t *testing.T) {
	repo := new(MockPingRepository)
	servicePing := service.NewPingService(repo)

	newPing := &model.Ping{IPAddress: "192.168.0.3", PingTime: 30}

	repo.On("IsIPExists", newPing.IPAddress).Return(false, nil)
	repo.On("CreatePing", newPing).Return(nil)

	err := servicePing.AddPing(newPing)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestPingService_AddPing_IPAlreadyExists(t *testing.T) {
	repo := new(MockPingRepository)
	servicePing := service.NewPingService(repo)

	existingPing := &model.Ping{IPAddress: "192.168.0.4", PingTime: 60}

	repo.On("IsIPExists", existingPing.IPAddress).Return(true, nil)

	err := servicePing.AddPing(existingPing)

	assert.Error(t, err)
	assert.EqualError(t, err, "IP address 192.168.0.4 already exists")
	repo.AssertExpectations(t)
}

func TestPingService_AddPing_RepoError(t *testing.T) {
	repo := new(MockPingRepository)
	servicePing := service.NewPingService(repo)

	newPing := &model.Ping{IPAddress: "192.168.0.5", PingTime: 70}

	repo.On("IsIPExists", newPing.IPAddress).Return(false, nil)
	repo.On("CreatePing", newPing).Return(errors.New("DB error"))

	err := servicePing.AddPing(newPing)

	assert.Error(t, err)
	assert.EqualError(t, err, "DB error")
	repo.AssertExpectations(t)
}
