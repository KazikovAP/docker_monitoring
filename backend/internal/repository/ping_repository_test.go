package repository_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/KazikovAP/docker_monitoring/backend/internal/model"
	"github.com/KazikovAP/docker_monitoring/backend/internal/repository"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func getTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open in-memory sqlite db: %v", err)
	}

	db.Exec("DELETE FROM pings")

	return db
}

func TestPingRepository_CreatePing(t *testing.T) {
	db := getTestDB(t)
	repo := repository.NewPingRepository(db)

	ping := &model.Ping{
		IPAddress:       "127.0.0.1",
		PingTime:        100,
		LastSuccessDate: time.Now(),
		CreatedAt:       time.Now(),
	}

	err := repo.CreatePing(ping)
	assert.NoError(t, err)

	var count int64
	err = db.Model(&model.Ping{}).Count(&count).Error
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)
}

func TestPingRepository_GetAllPings(t *testing.T) {
	db := getTestDB(t)
	repo := repository.NewPingRepository(db)

	pings := []model.Ping{
		{
			IPAddress:       fmt.Sprintf("127.0.0.%d", time.Now().UnixNano()),
			PingTime:        100,
			LastSuccessDate: time.Now(),
			CreatedAt:       time.Now(),
		},
		{
			IPAddress:       fmt.Sprintf("10.0.0.%d", time.Now().UnixNano()),
			PingTime:        200,
			LastSuccessDate: time.Now(),
			CreatedAt:       time.Now(),
		},
	}

	for i := range pings {
		err := repo.CreatePing(&pings[i])
		assert.NoError(t, err)
	}

	all, err := repo.GetAllPings()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(all))
}

func TestPingRepository_IsIPExists(t *testing.T) {
	db := getTestDB(t)
	repo := repository.NewPingRepository(db)

	exists, err := repo.IsIPExists("10.0.0.1")
	assert.NoError(t, err)
	assert.False(t, exists)

	ping := &model.Ping{
		IPAddress:       "10.0.0.1",
		PingTime:        150,
		LastSuccessDate: time.Now(),
		CreatedAt:       time.Now(),
	}
	err = repo.CreatePing(ping)
	assert.NoError(t, err)

	exists, err = repo.IsIPExists("10.0.0.1")
	assert.NoError(t, err)
	assert.True(t, exists)
}
