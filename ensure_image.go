package dockerlib

import (
	"context"
	"encoding/json"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type EnsureImageProgressDetail struct {
	Current int
	Total   int
}

type EnsureImageProgress struct {
	Status         string
	ProgressDetail EnsureImageProgressDetail
	Progress       string
	ID             string
}

func (p EnsureImageProgress) String() string {
	if len(p.ID) > 0 {
		return p.ID + " " + p.Status + " " + p.Progress
	} else {
		return p.Status
	}
}

func EnsureImage(ctx context.Context, cli *client.Client, image string) error {
	reader, err := cli.ImagePull(ctx, image, types.ImagePullOptions{})
	if err != nil {
		logger.Errorf("Unable to ensure image %s exists: %v", image, err)
		return err
	}

	defer reader.Close()
	lines := readLinesAsBytes(reader)
	for line := range lines {
		var progress EnsureImageProgress
		err := json.Unmarshal(line, &progress)
		if err != nil {
			logger.Error(line)
			continue
		}

		logger.Info(progress)
	}

	return nil
}
