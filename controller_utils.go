package dockerlib

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
	"strings"
)

// Helper method to follow logs of running container.
func (controller *DockerController) followLogs(containerID string, containerName string, readyChan chan<- bool, readyText string) {
	logOptions := types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true, Follow: true}

	// logs need to be in background context so they aren't canceled before container.
	reader, err := controller.cli.ContainerLogs(context.Background(), containerID, logOptions)
	if err != nil {
		logger.Errorf("Unable to follow logs for container %s: %v", containerName, err)
		return
	}
	defer reader.Close()

	cLogger := logger.Named(containerName)
	lines := ReadLinesAsBytes(reader)
	for line := range lines {
		text := string(line)
		cLogger.Info(text)
		if len(readyText) > 0 && strings.Contains(text, readyText) {
			readyChan <- true
			close(readyChan)
		}
	}

	logger.Infof("Logs finished for container %s", containerName)
}

func (controller *DockerController) attachNetworks(ctx context.Context, container Container) error {
	networks, err := controller.cli.NetworkList(ctx, types.NetworkListOptions{})
	if err != nil {
		logger.Errorf("Unable to list networks: %v", err)
		return DockerError{"unable to list networks", err}
	}

	// put names in map to easily find and keep track of whether added
	names := make(map[string]bool)
	for _, name := range container.Network {
		names[name] = false
	}

	for _, nw := range networks {
		_, exists := names[nw.Name]
		if !exists {
			continue
		}

		logger.Infof("Attaching network %+v to container %s", nw, container.Name)
		err := controller.cli.NetworkConnect(ctx, nw.ID, container.ID, &network.EndpointSettings{})
		if err != nil {
			logger.Errorf("Unable to attach network %s to container %s: %v", nw.Name, container.Name, err)
			return ContainerError{
				msg:           "unable to attach network " + nw.Name + " to container",
				containerName: container.Name,
				baseError:     err,
			}
		}

		names[nw.Name] = true
	}

	// now check how many networks were not attached
	var notFound []string
	for key, value := range names {
		if !value {
			notFound = append(notFound, key)
		}
	}

	if len(notFound) > 0 {
		logger.Errorf("Unable to find networks %v to attach to container %s", notFound, container.Name)
		return fmt.Errorf("unable to find networks %v to attach to container %s", notFound, container.Name)
	}

	return nil
}
