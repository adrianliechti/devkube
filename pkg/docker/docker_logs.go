package docker

import (
	"context"
	"io"
	"os"
	"os/exec"
)

type LogsOptions struct {
	Follow bool

	Stdout io.Writer
	Stderr io.Writer
}

func logsArgs(container string, options LogsOptions) []string {
	args := []string{
		"logs",
	}

	if options.Follow {
		args = append(args, "--follow")
	}

	args = append(args, container)

	return args
}

func Logs(ctx context.Context, container string, options LogsOptions) error {
	tool, _, err := Tool(ctx)

	if err != nil {
		return err
	}

	if options.Stdout == nil {
		options.Stdout = os.Stdout
	}

	if options.Stderr == nil {
		options.Stderr = os.Stderr
	}

	logs := exec.CommandContext(ctx, tool, logsArgs(container, options)...)
	logs.Stdout = options.Stdout
	logs.Stderr = options.Stderr

	return logs.Run()
}
