package kind

import (
	"bufio"
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
	minimalVersion = semver.MustParse("0.13.0")

	errNotFound = errors.New("kind not found. see https://kind.sigs.k8s.io/docs/user/quick-start/")
	errOutdated = errors.New("kind is outdated. see https://kind.sigs.k8s.io/docs/user/quick-start/")
)

func Info(ctx context.Context) (string, *semver.Version, error) {
	return path(ctx)
}

func path(ctx context.Context) (string, *semver.Version, error) {
	name := "kind"

	if runtime.GOOS == "windows" {
		name = "kind.exe"
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
	cmd := exec.CommandContext(ctx, path, "version", "-q")
	data, err := cmd.Output()

	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")

	if len(lines) == 0 {
		return nil, errors.New("invalid kind version")
	}

	return semver.NewVersion(lines[0])
}

type Option func(h *Kind)

type Kind struct {
	stdin  io.Reader
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

func WithInput(stdout, stdin io.Reader) Option {
	return func(k *Kind) {
		k.stdin = stdin
	}
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

func (k *Kind) Invoke(ctx context.Context, arg ...string) error {
	path, _, err := Info(ctx)

	if err != nil {
		return err
	}

	cmd := exec.CommandContext(ctx, path, arg...)
	cmd.Stdin = k.stdin
	cmd.Stdout = k.stdout
	cmd.Stderr = k.stderr

	return cmd.Run()
}

func List(ctx context.Context) ([]string, error) {
	path, _, err := Info(ctx)

	if err != nil {
		return nil, err
	}

	var list []string

	args := []string{
		"get", "clusters",
	}

	cmd := exec.CommandContext(ctx, path, args...)
	data, err := cmd.Output()

	if err != nil {
		return list, err
	}

	scanner := bufio.NewScanner(bytes.NewReader(data))

	for scanner.Scan() {
		name := scanner.Text()
		list = append(list, name)
	}

	return list, nil
}

func Create(ctx context.Context, name string, config map[string]any, kubeconfig string, opt ...Option) error {
	k := New(opt...)

	args := []string{
		"create", "cluster",
	}

	if name != "" {
		args = append(args, "--name", name)
	}

	if kubeconfig != "" {
		args = append(args, "--kubeconfig", kubeconfig)
	}

	if len(config) > 0 {
		data, err := yaml.Marshal(config)

		if err != nil {
			return err
		}

		k.stdin = bytes.NewReader(data)
		args = append(args, "--config", "-")
	}

	return k.Invoke(ctx, args...)
}

func Delete(ctx context.Context, name string, opt ...Option) error {
	k := New(opt...)

	args := []string{
		"delete", "cluster",
	}

	if name != "" {
		args = append(args, "--name", name)
	}

	return k.Invoke(ctx, args...)
}

func ExportConfig(ctx context.Context, name, kubeconfig string, opt ...Option) error {
	k := New(opt...)

	args := []string{
		"export", "kubeconfig", "--name", name, "--kubeconfig", kubeconfig,
	}

	return k.Invoke(ctx, args...)
}
