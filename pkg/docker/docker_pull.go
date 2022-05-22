package docker

import (
	"context"
	"os"
	"os/exec"
)

func Pull(ctx context.Context, image string) error {
	tool, _, err := Tool(ctx)

	if err != nil {
		return err
	}

	cmd := exec.CommandContext(ctx, tool, "pull", image)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
