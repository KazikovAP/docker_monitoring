package models

import "time"

type Container struct {
	ID        string
	Name      string
	IPAddress string
	Status    string
	CreatedAt time.Time
}
