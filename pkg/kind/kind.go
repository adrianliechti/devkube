package kind

import (
	"io"
	"os"
)

type Option func(h *Kind)

type Kind struct {
	stdout io.Writer
	stderr io.Writer
}

func New(options ...Option) *Kind {
	k := &Kind{}

	for _, option := range options {
		option(k)
	}

	return k
}

func WithOutput(stdout, stderr io.Writer) Option {
	return func(k *Kind) {
		k.stdout = stdout
		k.stderr = stderr
	}
}

func WithDefaultOutput() Option {
	return WithOutput(os.Stdout, os.Stderr)
}
