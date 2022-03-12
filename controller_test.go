package dockerlib

import (
	"context"
	"testing"
)

const TestImage = "alpine"

func TestEnsureImage(t *testing.T) {
	controller, err := NewDockerController()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	err = (*controller).EnsureImage(context.Background(), TestImage)
	if err != nil {
		t.Error(err)
	}
}
