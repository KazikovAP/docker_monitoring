package main

import (
	"log"
	"time"

	"github.com/KazikovAP/docker_monitoring/pinger/service"
)

func main() {
	containerService, err := service.NewDockerContainerService()
	if err != nil {
		log.Fatalf("Failed to create container service: %v", err)
	}

	pingService := service.NewPingService()
	apiService := service.NewAPIService("http://localhost:8080")

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		containers, err := containerService.GetContainers()
		if err != nil {
			log.Printf("Error getting containers: %v", err)
			continue
		}

		for _, container := range containers {
			pingResult := pingService.Ping(container.IPAddress)
			pingResult.ContainerID = container.ID

			if err := apiService.SendPingResult(&pingResult); err != nil {
				log.Printf("Error sending ping result: %v", err)
			}
		}
	}
}
