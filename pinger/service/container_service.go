package service

import (
	"time"

	"github.com/KazikovAP/docker_monitoring/pinger/models"
	docker "github.com/fsouza/go-dockerclient"
)

type ContainerService interface {
	GetContainers() ([]models.Container, error)
}

type DockerContainerService struct {
	client *docker.Client
}

func NewDockerContainerService() (*DockerContainerService, error) {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		return nil, err
	}

	return &DockerContainerService{client: client}, nil
}

func (d *DockerContainerService) GetContainers() ([]models.Container, error) {
	containers, err := d.client.ListContainers(docker.ListContainersOptions{All: true})
	if err != nil {
		return nil, err
	}

	result := make([]models.Container, 0, len(containers))

	for i := range containers {
		c := &containers[i]

		inspectData, err := d.client.InspectContainerWithOptions(docker.InspectContainerOptions{ID: c.ID})
		if err != nil {
			return nil, err
		}

		ipAddress := ""

		if len(inspectData.NetworkSettings.Networks) > 0 {
			for networkName := range inspectData.NetworkSettings.Networks {
				ipAddress = inspectData.NetworkSettings.Networks[networkName].IPAddress
				break
			}
		}

		if ipAddress == "" {
			continue
		}

		createdAt := time.Unix(c.Created, 0)

		result = append(result, models.Container{
			ID:        c.ID,
			Name:      c.Names[0],
			IPAddress: ipAddress,
			Status:    c.Status,
			CreatedAt: createdAt,
		})
	}

	return result, nil
}
