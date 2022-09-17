package helm

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/Masterminds/semver"
	"gopkg.in/yaml.v3"
)

var (
	minimalVersion = semver.MustParse("3.8.2")

	errNotFound = errors.New("helm not found. see https://helm.sh/docs/intro/install/")
	errOutdated = errors.New("helm is outdated. see https://helm.sh/docs/intro/install/")
)

func Info(ctx context.Context) (string, *semver.Version, error) {
	return path(ctx)
}

func path(ctx context.Context) (string, *semver.Version, error) {
	name := "helm"

	if runtime.GOOS == "windows" {
		name = "helm.exe"
	}

	if path, err := exec.LookPath(name); err == nil {
		if version, err := version(ctx, path); err == nil {
			if !version.LessThan(minimalVersion) {
				return path, version, nil
			}

			return path, version, errOutdated
		}

		return path, nil, errOutdated
	}

	return "", nil, errNotFound
}

func version(ctx context.Context, path string) (*semver.Version, error) {
	cmd := exec.CommandContext(ctx, path, "version", "--short")
	data, err := cmd.Output()

	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")

	if len(lines) == 0 {
		return nil, errors.New("invalid helm version")
	}

	return semver.NewVersion(lines[0])
}

type Option func(h *Helm)

type Helm struct {
	kubeconfig string

	context   string
	namespace string

	wait bool

	stdin  io.Reader
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

func WithWait(wait bool) Option {
	return func(h *Helm) {
		h.wait = wait
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

func (h *Helm) Invoke(ctx context.Context, arg ...string) error {
	path, _, err := Info(ctx)

	if err != nil {
		return err
	}

	if h.kubeconfig != "" {
		arg = append(arg, "--kubeconfig", h.kubeconfig)
	}

	if h.context != "" {
		arg = append(arg, "--kube-context", h.context)
	}

	if h.namespace != "" {
		arg = append(arg, "--namespace", h.namespace)
	}

	cmd := exec.CommandContext(ctx, path, arg...)
	cmd.Stdin = h.stdin
	cmd.Stdout = h.stdout
	cmd.Stderr = h.stderr

	return cmd.Run()
}

func Install(ctx context.Context, release, repo, chart, version string, values map[string]interface{}, opt ...Option) error {
	h := New(opt...)

	args := []string{
		"upgrade", "--install", "--create-namespace", "--timeout", "10m0s",
		release,
		chart,
	}

	if repo != "" {
		args = append(args, "--repo", repo)
	}

	if version != "" {
		args = append(args, "--version", version)
	}

	if h.wait {
		args = append(args, "--wait")
	}

	if len(values) > 0 {
		args = append(args, "-f", "-")

		data, err := yaml.Marshal(values)

		if err != nil {
			return err
		}

		h.stdin = bytes.NewReader(data)
	}

	return h.Invoke(ctx, args...)
}

func Uninstall(ctx context.Context, release string, opt ...Option) error {
	h := New(opt...)

	args := []string{
		"uninstall",
		release,
	}

	if h.wait {
		args = append(args, "--wait")
	}

	return h.Invoke(ctx, args...)
}
