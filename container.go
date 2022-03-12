package dockerlib

import (
	"fmt"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/go-connections/nat"
)

// Container represents a simplified interface for starting a Docker container
type Container struct {
	Name        string
	Image       string
	ID          string
	Mounts      []mount.Mount
	Ports       map[int]int
	Command     []string
	Environment []string
}

// Returns a simplified string representation
func (c Container) String() string {
	return fmt.Sprintf("%s (%s/%s)", c.Name, c.Image, c.ID)
}

// PortBindings Helper method to return the structs required to start a Docker container, or any error
func (c Container) PortBindings() (map[nat.Port]struct{}, map[nat.Port][]nat.PortBinding, error) {
	specs := make([]string, len(c.Ports))
	i := 0
	for from, to := range c.Ports {
		// ip:public:private/proto
		specs[i] = fmt.Sprintf("0.0.0.0:%d:%d/tcp", to, from)
		i += 1
	}

	return nat.ParsePortSpecs(specs)
}
