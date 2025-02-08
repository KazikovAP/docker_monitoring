package models

import "time"

type PingResult struct {
	ContainerID     string
	IPAddress       string
	PingTime        int64
	LastSuccessDate time.Time
	Timestamp       time.Time
}
