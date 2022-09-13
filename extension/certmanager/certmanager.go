package certmanager

import (
	"context"
	_ "embed"
	"strings"

	"github.com/adrianliechti/devkube/pkg/kubectl"
)

//go:embed ca.yaml
var manifestCA string

//go:embed ca-trust.yaml
var manifestDaemonSet string

const (
	certmanagerVersion   = "v1.9.1"
	certmanagerNamespace = "cert-manager"
)

func Install(ctx context.Context, kubeconfig, namespace string) error {
	if namespace == "" {
		namespace = "default"
	}

	url := "https://github.com/cert-manager/cert-manager/releases/download/" + certmanagerVersion + "/cert-manager.yaml"

	if err := kubectl.Invoke(ctx, []string{"apply", "-f", url, "--validate=false", "--server-side=true", "--overwrite=true"}, kubectl.WithKubeconfig(kubeconfig), kubectl.WithDefaultOutput()); err != nil {
		return err
	}

	if err := kubectl.Invoke(ctx, []string{"wait", "deployments", "-n", certmanagerNamespace, "--all", "--for", "condition=Available", "--timeout", "300s"}, kubectl.WithKubeconfig(kubeconfig), kubectl.WithDefaultOutput()); err != nil {
		return err
	}

	if err := kubectl.Invoke(ctx, []string{"apply", "--validate=false", "-f", "-"}, kubectl.WithKubeconfig(kubeconfig), kubectl.WithInput(strings.NewReader(manifestCA)), kubectl.WithDefaultOutput()); err != nil {
		return err
	}

	if err := kubectl.Invoke(ctx, []string{"apply", "--validate=false", "-f", "-"}, kubectl.WithKubeconfig(kubeconfig), kubectl.WithInput(strings.NewReader(manifestDaemonSet)), kubectl.WithDefaultOutput()); err != nil {
		return err
	}

	return nil
}

func Uninstall(ctx context.Context, kubeconfig, namespace string) error {
	if namespace == "" {
		namespace = "default"
	}

	url := "https://github.com/cert-manager/cert-manager/releases/download/" + certmanagerVersion + "/cert-manager.yaml"

	if err := kubectl.Invoke(ctx, []string{"delete", "-f", url}, kubectl.WithKubeconfig(kubeconfig), kubectl.WithDefaultOutput()); err != nil {
		return err
	}

	return nil
}
