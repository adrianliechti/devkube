package docker

import (
	"context"
	"io"
	"os"
	"os/exec"
)

type ExecOptions struct {
	Privileged bool

	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer

	TTY         bool
	Interactive bool

	Dir  string
	User string

	Env map[string]string
}

func execArgs(container string, options ExecOptions, command string, arg ...string) []string {
	args := []string{
		"exec",
	}

	if options.User != "" {
		args = append(args, "--user", options.User)
	}

	if options.Privileged {
		args = append(args, "--privileged")
	}

	if options.Interactive {
		args = append(args, "--interactive")
	}

	if options.TTY {
		args = append(args, "--tty")
	}

	if options.Dir != "" {
		args = append(args, "--workdir", options.Dir)
	}

	for key, value := range options.Env {
		args = append(args, "--env", key+"="+value)
	}

	args = append(args, container, command)
	args = append(args, arg...)

	return args
}

func Exec(ctx context.Context, container string, options ExecOptions, command string, args ...string) error {
	tool, _, err := Tool(ctx)

	if err != nil {
		return err
	}

	run := exec.CommandContext(ctx, tool, execArgs(container, options, command, args...)...)
	run.Stdin = options.Stdin
	run.Stdout = options.Stdout
	run.Stderr = options.Stderr

	return run.Run()
}

func ExecInteractive(ctx context.Context, container string, options ExecOptions, shell string, args ...string) error {
	tool, _, err := Tool(ctx)

	if err != nil {
		return err
	}

	if options.Stdin == nil {
		options.Stdin = os.Stdin
	}

	if options.Stdout == nil {
		options.Stdout = os.Stdout
	}

	if options.Stderr == nil {
		options.Stderr = os.Stderr
	}

	options.TTY = true
	options.Interactive = true

	run := exec.CommandContext(ctx, tool, execArgs(container, options, shell, args...)...)
	run.Stdin = options.Stdin
	run.Stdout = options.Stdout
	run.Stderr = options.Stderr

	return run.Run()
}
