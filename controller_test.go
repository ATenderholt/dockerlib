package dockerlib_test

import (
	"context"
	"github.com/ATenderholt/dockerlib"
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

func TestStartImage(t *testing.T) {
	controller, err := dockerlib.NewDockerController()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	container := dockerlib.Container{
		Name:        "dockerlib-test",
		Image:       TestImage,
		Mounts:      nil,
		Ports:       nil,
		Command:     nil,
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
