package dockerlib_test

import (
	"context"
	"github.com/ATenderholt/dockerlib"
	"github.com/docker/docker/api/types/mount"
	"os"
	"path/filepath"
	"testing"
	"time"
)

const TestImage = "python:3.9.11-alpine3.14"

func TestEnsureImage(t *testing.T) {
	controller, err := dockerlib.NewDockerController()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	err = (*controller).EnsureImage(context.Background(), TestImage)
	if err != nil {
		t.Error(err)
	}
}

func TestBasicStartImage(t *testing.T) {
	controller, err := dockerlib.NewDockerController()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cwd, err := os.Getwd()
	if err != nil {
		t.Errorf("Unable to get cwd: %v", err)
		t.FailNow()
	}

	container := dockerlib.Container{
		Name:  "dockerlib-test",
		Image: TestImage,
		Mounts: []mount.Mount{
			{
				Source:      filepath.Join(cwd, "testdata"),
				Target:      "/scripts",
				Type:        mount.TypeBind,
				ReadOnly:    true,
				Consistency: mount.ConsistencyDelegated,
			},
		},
		Ports:       nil,
		Command:     []string{"python", "/scripts/basic.py"},
		Environment: nil,
		Network:     nil,
	}

	ready, err := controller.Start(ctx, container, "Hello")
	if err != nil {
		t.Error(err)
	}
	defer controller.ShutdownAll(context.Background())

	select {
	case <-ready:
		return
	case <-ctx.Done():
		t.Error("Test timeout - container didn't start.")
		t.FailNow()
	}
}

func TestStartWithNetwork(t *testing.T) {
	controller, err := dockerlib.NewDockerController()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	network := "dockerlib"
	err = controller.EnsureNetwork(ctx, network)
	if err != nil {
		t.Errorf("Unable to create network: %v", err)
	}
	defer controller.CleanupNetworks(context.Background())

	cwd, err := os.Getwd()
	if err != nil {
		t.Errorf("Unable to get cwd: %v", err)
		t.FailNow()
	}

	server := dockerlib.Container{
		Name:  "dockerlib-test-server",
		Image: TestImage,
		Mounts: []mount.Mount{
			{
				Source:      filepath.Join(cwd, "testdata", "server.py"),
				Target:      "/scripts/server.py",
				Type:        mount.TypeBind,
				ReadOnly:    true,
				Consistency: mount.ConsistencyDelegated,
			},
			{
				Source:      filepath.Join(cwd, "testdata", "hello.txt"),
				Target:      "/site/hello.txt",
				Type:        mount.TypeBind,
				ReadOnly:    true,
				Consistency: mount.ConsistencyDelegated,
			},
		},
		Ports:       nil,
		Command:     []string{"python", "/scripts/server.py", "/site"},
		Environment: nil,
		Network:     []string{network},
	}

	client := dockerlib.Container{
		Name:  "dockerlib-test-client",
		Image: TestImage,
		Mounts: []mount.Mount{
			{
				Source:      filepath.Join(cwd, "testdata", "client.py"),
				Target:      "/scripts/client.py",
				Type:        mount.TypeBind,
				ReadOnly:    true,
				Consistency: mount.ConsistencyDelegated,
			},
		},
		Ports:       nil,
		Command:     []string{"python", "/scripts/client.py", "http://dockerlib-test-server:8000"},
		Environment: nil,
		Network:     []string{network},
	}

	ready, err := controller.Start(ctx, server, "Server started on port")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer controller.Shutdown(context.Background(), server)

	select {
	case <-ready:
	case <-ctx.Done():
		t.Error("Test timeout - server didn't start.")
		t.FailNow()
	}

	status, err := controller.Start(ctx, client, "Status")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer controller.ShutdownAll(context.Background())

	select {
	case <-status:
		return
	case <-ctx.Done():
		t.Error("Test timeout - client didn't start.")
		t.FailNow()
	}
}
