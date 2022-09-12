package linkerd

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/adrianliechti/devkube/pkg/cli"
	"github.com/adrianliechti/devkube/pkg/kubectl"
)

func Install(ctx context.Context, kubeconfig, namespace string) error {
	if namespace == "" {
		namespace = "default"
	}

	cli.Info("Downloading Linkerd CLI...")
	cli, closer, err := downloadCLI(ctx)

	if err != nil {
		return err
	}

	defer closer()

	crdManifest, err := exec.Command(cli, "install", "--crds", "--kubeconfig", kubeconfig).Output()

	if err != nil {
		return err
	}

	if err := kubectl.Invoke(ctx, []string{"apply", "-f", "-"}, kubectl.WithKubeconfig(kubeconfig), kubectl.WithInput(bytes.NewReader(crdManifest)), kubectl.WithDefaultOutput()); err != nil {
		return err
	}

	installManifest, err := exec.Command(cli, "install", "--kubeconfig", kubeconfig).Output()

	if err != nil {
		return err
	}

	if err := kubectl.Invoke(ctx, []string{"apply", "-f", "-"}, kubectl.WithKubeconfig(kubeconfig), kubectl.WithInput(bytes.NewReader(installManifest)), kubectl.WithDefaultOutput()); err != nil {
		return err
	}

	return nil
}

func Uninstall(ctx context.Context, kubeconfig, namespace string) error {
	if namespace == "" {
		namespace = "default"
	}

	cli.Info("Downloading Linkerd CLI...")
	cli, closer, err := downloadCLI(ctx)

	if err != nil {
		return err
	}

	defer closer()

	uninstallManifest, err := exec.Command(cli, "uninstall", "--kubeconfig", kubeconfig).Output()

	if err != nil {
		return err
	}

	if err := kubectl.Invoke(ctx, []string{"delete", "-f", "-"}, kubectl.WithKubeconfig(kubeconfig), kubectl.WithInput(bytes.NewReader(uninstallManifest)), kubectl.WithDefaultOutput()); err != nil {
		return err
	}

	return nil
}

func downloadCLI(ctx context.Context) (string, func(), error) {
	var name string

	if runtime.GOOS == "windows" {
		name = "linkerd2-cli-stable-2.12.0-windows.exe"
	}

	if runtime.GOOS == "darwin" {
		name = "linkerd2-cli-stable-2.12.0-darwin"

		if runtime.GOARCH == "arm64" {
			name = "linkerd2-cli-stable-2.12.0-darwin-arm64"
		}
	}

	if runtime.GOOS == "linux" {
		if runtime.GOARCH == "amd64" {
			name = "linkerd2-cli-stable-2.12.0-linux-amd64"
		}

		if runtime.GOARCH == "arm" {
			name = "linkerd2-cli-stable-2.12.0-linux-arm"
		}

		if runtime.GOARCH == "arm64" {
			name = "linkerd2-cli-stable-2.12.0-linux-arm64"
		}
	}

	if name == "" {
		return "", nil, errors.New("unsupported platform")
	}

	dir, err := os.MkdirTemp("", "devkube")

	if err != nil {
		return "", nil, err
	}

	closer := func() {
		os.RemoveAll(dir)
	}

	resp, err := http.Get("https://github.com/linkerd/linkerd2/releases/download/stable-2.12.0/" + name)

	if err != nil {
		closer()
		return "", nil, err
	}

	defer resp.Body.Close()

	path := filepath.Join(dir, name)

	out, err := os.Create(path)

	if err != nil {
		closer()
		return "", nil, err
	}

	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		closer()
		return "", nil, err
	}

	if err := os.Chmod(path, 0755); err != nil {
		closer()
		return "", nil, err
	}

	return path, closer, nil
}
