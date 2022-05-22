package docker

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
)

type RunOptions struct {
	Name   string
	Labels map[string]string

	Platform string

	Temporary  bool
	Privileged bool

	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer

	Attach      bool
	TTY         bool
	Interactive bool

	Dir  string
	User string

	Env     map[string]string
	Ports   map[int]int
	Volumes map[string]string
}

func runArgs(image string, options RunOptions, arg ...string) []string {
	args := []string{
		"run",
	}

	if options.Name != "" {
		args = append(args, "--name", options.Name)
	}

	for k, v := range options.Labels {
		args = append(args, "--label", k+"="+v)
	}

	if options.User != "" {
		args = append(args, "--user", options.User)
	}

	if options.Platform != "" {
		args = append(args, "--platform", options.Platform)
	}

	if options.Temporary {
		args = append(args, "--rm")
	}

	if options.Privileged {
		args = append(args, "--privileged")
	}

	if !options.Attach {
		args = append(args, "--detach")
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

	for source, target := range options.Ports {
		args = append(args, "--publish", fmt.Sprintf("127.0.0.1:%d:%d", source, target))
	}

	for source, target := range options.Volumes {
		args = append(args, "--volume", source+":"+target)
	}

	args = append(args, image)
	args = append(args, arg...)

	return args
}

func Run(ctx context.Context, image string, options RunOptions, args ...string) error {
	tool, _, err := Tool(ctx)

	if err != nil {
		return err
	}

	run := exec.CommandContext(ctx, tool, runArgs(image, options, args...)...)
	run.Stdin = options.Stdin
	run.Stdout = options.Stdout
	run.Stderr = options.Stderr

	return run.Run()
}

func RunInteractive(ctx context.Context, image string, options RunOptions, args ...string) error {
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

	options.Temporary = true

	options.TTY = true
	options.Attach = true
	options.Interactive = true

	run := exec.CommandContext(ctx, tool, runArgs(image, options, args...)...)
	run.Stdin = options.Stdin
	run.Stdout = options.Stdout
	run.Stderr = options.Stderr

	return run.Run()
}
