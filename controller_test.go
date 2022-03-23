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
