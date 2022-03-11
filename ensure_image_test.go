package dockerlib

import (
	"context"
	"github.com/docker/docker/client"
	"testing"
)

const TestImage = "alpine"

func TestEnsureImage(t *testing.T) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	err = EnsureImage(context.Background(), cli, TestImage)
	if err != nil {
		t.Error(err)
	}
}
