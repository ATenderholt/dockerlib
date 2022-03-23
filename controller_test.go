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

const TestImage = "mlupin/docker-lambda:python3.8"

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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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
				Target:      "/var/task",
				Type:        mount.TypeBind,
				ReadOnly:    true,
				Consistency: mount.ConsistencyDelegated,
			},
		},
		Ports:       nil,
		Command:     []string{"basic.handler"},
		Environment: nil,
		Network:     nil,
	}

	ready, err := controller.Start(ctx, container, "Got event")
	if err != nil {
		t.Error(err)
	}
	defer controller.ShutdownAll(ctx)

	for {
		select {
		case <-ready:
			return
		case <-ctx.Done():
			t.Error("Test timeout - container didn't start.")
			t.FailNow()
		}
	}
}
