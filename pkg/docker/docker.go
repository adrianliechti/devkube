package docker

import (
	"context"
	"errors"
	"os/exec"
	"runtime"
	"strings"

	"github.com/Masterminds/semver"
)

var (
	minimalVersion = semver.MustParse("19.0.0")

	errNotFound   = errors.New("docker not found. see https://docs.docker.com/get-docker/")
	errOutdated   = errors.New("docker is outdated. see https://docs.docker.com/get-docker/")
	errNotRunning = errors.New("docker seems not to be running")
)

func Info(ctx context.Context) (string, *semver.Version, error) {
	path, version, err := path(ctx)

	if err != nil {
		return path, version, err
	}

	cmd := exec.CommandContext(ctx, path, "info")

	if err := cmd.Run(); err == nil {
		return path, version, nil
	}

	return path, version, errNotRunning
}

func path(ctx context.Context) (string, *semver.Version, error) {
	name := "docker"

	if runtime.GOOS == "windows" {
		name = "docker.exe"
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
	output, _ := exec.CommandContext(ctx, path, "version", "--format", "{{.Client.Version}}").Output()

	parts := strings.Split(strings.TrimSuffix(string(output), "\n"), "\n")

	if len(parts) < 1 {
		return nil, errors.New("unable to get docker version")
	}

	version := strings.TrimSpace(parts[0])
	return semver.NewVersion(version)
}
