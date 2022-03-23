package dockerlib_test

import (
	"github.com/ATenderholt/dockerlib"
	"github.com/docker/go-connections/nat"
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestContainerPortBindings(t *testing.T) {
	container := dockerlib.Container{Ports: map[int]int{
		123: 234,
	}}

	portSet, portMap, err := container.PortBindings()
	if err != nil {
		t.Errorf("unexpected error when getting port bindings: %v", err)
	}

	if _, ok := portSet["123/tcp"]; !ok {
		t.Errorf("port 123 not in port set")
	}

	value, ok := portMap["123/tcp"]
	if !ok {
		t.Errorf("port 123 not in port map")
	}

	expected := []nat.PortBinding{
		{
			HostIP:   "0.0.0.0",
			HostPort: "234",
		},
	}

	if !cmp.Equal(value, expected) {
		t.Errorf("bindings on port 123 not correct: %+v", value)
	}
}
