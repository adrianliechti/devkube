package kind

import (
	"bufio"
	"bytes"
	"context"
	"errors"
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

	// verify global tool
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

func List(ctx context.Context) ([]string, error) {
	var list []string

	tool, _, err := Info(ctx)

	if err != nil {
		return list, err
	}

	args := []string{
		"get", "clusters",
	}

	cmd := exec.CommandContext(ctx, tool, args...)
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

func Create(ctx context.Context, name string, config map[string]any, kubeconfig string) error {
	tool, _, err := Info(ctx)

	if err != nil {
		return err
	}

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
		args = append(args, "--config", "-")
	}

	cmd := exec.CommandContext(ctx, tool, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if len(config) > 0 {
		data, err := yaml.Marshal(config)

		if err != nil {
			return err
		}

		cmd.Stdin = bytes.NewReader(data)
	}

	return cmd.Run()
}

func Delete(ctx context.Context, name string) error {
	tool, _, err := Info(ctx)

	if err != nil {
		return err
	}

	args := []string{
		"delete", "cluster",
	}

	if name != "" {
		args = append(args, "--name", name)
	}

	cmd := exec.CommandContext(ctx, tool, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func ExportConfig(ctx context.Context, name, path string) error {
	tool, _, err := Info(ctx)

	if err != nil {
		return err
	}

	args := []string{
		"export", "kubeconfig", "--name", name, "--kubeconfig", path,
	}

	cmd := exec.CommandContext(ctx, tool, args...)
	return cmd.Run()
}
