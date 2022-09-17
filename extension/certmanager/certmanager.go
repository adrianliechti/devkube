package certmanager

import (
	"context"
	_ "embed"
	"strings"

	"github.com/adrianliechti/devkube/pkg/helm"
	"github.com/adrianliechti/devkube/pkg/kubectl"
)

const (
	certmanagerRepo       = "https://charts.jetstack.io"
	certmanagerNamespace1 = "cert-manager"

	certmanager        = "cert-manager"
	certmanagerChart   = "cert-manager"
	certmanagerVersion = "v1.9.1"
)

var (
	//go:embed manifest.yaml
	manifest string
)

func Install(ctx context.Context, kubeconfig, namespace string) error {
	// if namespace == "" {
	// 	namespace = "default"
	// }

	namespace = certmanagerNamespace1

	values := map[string]any{
		"installCRDs": true,

		"prometheus": map[string]any{
			"servicemonitor": map[string]any{
				"enabled": true,
			},
		},
	}

	if err := helm.Install(ctx, certmanager, certmanagerRepo, certmanagerChart, certmanagerVersion, values, helm.WithKubeconfig(kubeconfig), helm.WithNamespace(namespace), helm.WithWait(true), helm.WithDefaultOutput()); err != nil {
		return err
	}

	if err := kubectl.Invoke(ctx, []string{"apply", "-f", "-"}, kubectl.WithKubeconfig(kubeconfig), kubectl.WithNamespace(namespace), kubectl.WithInput(strings.NewReader(manifest)), kubectl.WithDefaultOutput()); err != nil {
		return err
	}

	return nil
}

func Uninstall(ctx context.Context, kubeconfig, namespace string) error {
	// if namespace == "" {
	// 	namespace = "default"
	// }

	namespace = certmanagerNamespace1

	if err := kubectl.Invoke(ctx, []string{"delete", "-f", "-"}, kubectl.WithKubeconfig(kubeconfig), kubectl.WithNamespace(namespace), kubectl.WithInput(strings.NewReader(manifest)), kubectl.WithDefaultOutput()); err != nil {
		// return err
	}

	if err := helm.Uninstall(ctx, certmanager, helm.WithKubeconfig(kubeconfig), helm.WithNamespace(namespace)); err != nil {
		//return err
	}

	return nil
}
