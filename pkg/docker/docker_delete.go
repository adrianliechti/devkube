package docker

import (
	"context"
	"os/exec"
)

type DeleteOptions struct {
	Force   bool
	Volumes bool
}

func deleteArgs(container string, options DeleteOptions) []string {
	args := []string{
		"rm",
	}

	if options.Force {
		args = append(args, "--force")
	}

	if options.Volumes {
		args = append(args, "--volumes")
	}

	args = append(args, container)

	return args
}

func Delete(ctx context.Context, container string, options DeleteOptions) error {
	tool, _, err := Tool(ctx)

	if err != nil {
		return err
	}

	stop := exec.CommandContext(ctx, tool, "stop", container)

	if err := stop.Run(); err != nil {
		return err
	}

	rm := exec.CommandContext(ctx, tool, deleteArgs(container, options)...)

	return rm.Run()
}
