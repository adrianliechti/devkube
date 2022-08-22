package helm

import (
	"context"
	"io"
	"os"
)

type Option func(h *Helm)

type Helm struct {
	kubeconfig string

	context   string
	namespace string

	stdout io.Writer
	stderr io.Writer
}

func New(options ...Option) *Helm {
	h := &Helm{}

	for _, option := range options {
		option(h)
	}

	return h
}

func WithKubeconfig(kubeconfig string) Option {
	return func(h *Helm) {
		h.kubeconfig = kubeconfig
	}
}

func WithContext(context string) Option {
	return func(h *Helm) {
		h.context = context
	}
}

func WithNamespace(namespace string) Option {
	return func(h *Helm) {
		h.namespace = namespace
	}
}

func WithOutput(stdout, stderr io.Writer) Option {
	return func(h *Helm) {
		h.stdout = stdout
		h.stderr = stderr
	}
}

func WithDefaultOutput() Option {
	return WithOutput(os.Stdout, os.Stderr)
}

func Install(ctx context.Context, release, repo, chart, version string, values map[string]interface{}, opt ...Option) error {
	h := New(opt...)
	return h.Install(ctx, release, repo, chart, version, values)
}

func Uninstall(ctx context.Context, release string, opt ...Option) error {
	h := New(opt...)
	return h.Uninstall(ctx, release)
}
