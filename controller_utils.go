package dockerlib

import (
	"context"
	"github.com/docker/docker/api/types"
	"strings"
)

// Helper method to follow logs of running container.
func (controller DockerController) followLogs(containerID string, containerName string, readyChan chan<- bool, readyText string) {
	logOptions := types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true, Follow: true}

	// logs need to be in background context so they aren't canceled before container.
	reader, err := controller.cli.ContainerLogs(context.Background(), containerID, logOptions)
	if err != nil {
		logger.Errorf("Unable to follow logs for container %s: %v", containerName, err)
		return
	}
	defer reader.Close()

	cLogger := logger.Named(containerName)
	lines := readLinesAsBytes(reader)
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
