package dockerlib

import (
	"context"
	"errors"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"strings"
	"time"
)

type Controller interface {
	Start(ctx context.Context, c Container, ready string) (chan bool, error)
	Shutdown(ctx context.Context, c Container) error
	ShutdownAll(ctx context.Context) error
}

type DockerController struct {
	cli     *client.Client
	running map[string]Container
}

func NewDockerController() (*Controller, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		logger.Errorf("Unable to create Docker client: %v", err)
		return nil, DockerError{"unable to create Docker client", err}
	}

	running := make(map[string]Container, 5)

	var c Controller
	c = DockerController{cli: cli, running: running}
	return &c, nil
}

func (controller DockerController) Start(ctx context.Context, c Container, ready string) (chan bool, error) {
	logger := logger.Named(c.Name)

	portSet, portMap, err := c.PortBindings()
	if err != nil {
		logger.Errorf("Unable to get port bindings: %v", err)
		return nil, ContainerError{"unable to get port bindings for container", c.Name, err}
	}

	hostConfig := container.HostConfig{}
	hostConfig.Mounts = c.Mounts
	hostConfig.PortBindings = portMap

	containerConfig := container.Config{
		ExposedPorts: portSet,
		Tty:          false,
		Cmd:          c.Command,
		Image:        c.Image,
		Env:          c.Environment,
	}

	resp, err := controller.cli.ContainerCreate(ctx, &containerConfig, &hostConfig, nil, nil, c.Name)
	if err != nil {
		logger.Errorf("Unable to create container %s: %v", c, err)
		return nil, ContainerError{"unable to create container", c.Name, err}
	}

	err = controller.cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		logger.Errorf("Unable to start container %s: %v", c, err)
		return nil, ContainerError{"unable to start container", c.Name, err}
	}

	c.ID = resp.ID
	controller.running[c.Name] = c

	readyChan := make(chan bool)
	go controller.followLogs(resp.ID, c.Name, readyChan, ready)

	return readyChan, nil
}

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

func (controller DockerController) Shutdown(ctx context.Context, c Container) error {
	logger.Infof("Trying to shutdown %s...", c)

	timeout := 30 * time.Second
	err := controller.cli.ContainerStop(ctx, c.ID, &timeout)
	if err != nil {
		logger.Errorf("Unable to shutdown container %s: %v", c, err)
		return ContainerError{"unable to shutdown container", c.Name, err}
	}

	delete(controller.running, c.Name)

	err = controller.cli.ContainerRemove(ctx, c.ID, types.ContainerRemoveOptions{})
	if err != nil {
		logger.Errorf("Unable to remove container %s: %v", c, err)
		return ContainerError{"unable to remove container", c.Name, err}
	}

	return nil
}

func (controller DockerController) ShutdownAll(ctx context.Context) error {
	var allErrors []string
	for _, c := range controller.running {
		err := controller.Shutdown(ctx, c)
		if err != nil {
			allErrors = append(allErrors, err.Error())
		}
	}

	msg := strings.Join(allErrors, ",")
	if len(msg) > 0 {
		return errors.New("errors encountered when shutting down all containers: " + msg)
	}

	return nil
}
