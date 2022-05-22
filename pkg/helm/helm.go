package helm

import (
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
	minimalVersion = semver.MustParse("3.8.2")

	errNotFound = errors.New("helm not found. see https://helm.sh/docs/intro/install/")
	errOutdated = errors.New("helm is outdated. see https://helm.sh/docs/intro/install/")
)

func Tool(ctx context.Context) (string, *semver.Version, error) {
	if path, version, err := Path(ctx); err == nil {
		return path, version, err
	}

	return "", nil, errNotFound
}

func Path(ctx context.Context) (string, *semver.Version, error) {
	name := "helm"

	if runtime.GOOS == "windows" {
		name = "helm.exe"
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
		return nil, errors.New("invalid helm version")
	}

	return semver.NewVersion(lines[0])
}

func Install(ctx context.Context, kubeconfig, namespace, release, repo, chart, version string, values map[string]interface{}) error {
	tool, _, err := Tool(ctx)

	if err != nil {
		return err
	}

	args := []string{
		"upgrade", "--install", "--create-namespace",
		release,
		chart,
	}

	if kubeconfig != "" {
		args = append(args, "--kubeconfig", kubeconfig)
	}

	if repo != "" {
		args = append(args, "--repo", repo)
	}

	if version != "" {
		args = append(args, "--version", version)
	}

	if namespace != "" {
		args = append(args, "--namespace", namespace)
	}

	if len(values) > 0 {
		args = append(args, "-f", "-")
	}

	cmd := exec.CommandContext(ctx, tool, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if len(values) > 0 {
		data, err := yaml.Marshal(values)

		if err != nil {
			return err
		}

		cmd.Stdin = bytes.NewReader(data)
	}

	return cmd.Run()
}

func Uninstall(ctx context.Context, kubeconfig, namespace, release string) error {
	tool, _, err := Tool(ctx)

	if err != nil {
		return err
	}

	args := []string{
		"uninstall",
		release,
	}

	if kubeconfig != "" {
		args = append(args, "--kubeconfig", kubeconfig)
	}

	if namespace != "" {
		args = append(args, "--namespace", namespace)
	}

	cmd := exec.CommandContext(ctx, tool, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
