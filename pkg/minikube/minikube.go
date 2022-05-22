package minikube

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/Masterminds/semver"
)

var (
	minimalVersion = semver.MustParse("1.25.2")

	errNotFound = errors.New("minikube not found. see https://minikube.sigs.k8s.io/docs/start/")
	errOutdated = errors.New("minikube is outdated. see https://minikube.sigs.k8s.io/docs/start/")
)

func Tool(ctx context.Context) (string, *semver.Version, error) {
	if path, version, err := Path(ctx); err == nil {
		return path, version, err
	}

	return "", nil, errNotFound
}

func Path(ctx context.Context) (string, *semver.Version, error) {
	name := "minikube"

	if runtime.GOOS == "windows" {
		name = "minikube.exe"
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
	cmd := exec.CommandContext(ctx, path, "version", "--short")
	data, err := cmd.Output()

	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")

	if len(lines) == 0 {
		return nil, errors.New("invalid minikube version")
	}

	return semver.NewVersion(lines[0])
}

func Create(ctx context.Context, profile string) error {
	tool, _, err := Tool(ctx)

	if err != nil {
		return err
	}

	args := []string{
		"start",
	}

	if profile != "" {
		args = append(args, "-p", profile)
	}

	cmd := exec.CommandContext(ctx, tool, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func Delete(ctx context.Context, profile string) error {
	tool, _, err := Tool(ctx)

	if err != nil {
		return err
	}

	args := []string{
		"delete",
	}

	if profile != "" {
		args = append(args, "-p", profile)
	}

	cmd := exec.CommandContext(ctx, tool, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
