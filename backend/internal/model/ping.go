package model

import "time"

type Ping struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	IPAddress       string    `gorm:"uniqueIndex" json:"ip_address" binding:"required"`
	PingTime        int64     `json:"ping_time"`
	LastSuccessDate time.Time `json:"last_success_date"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
