package observability

import (
	"context"

	"github.com/adrianliechti/devkube/pkg/kubectl"
)

const (
	falcoRepo      = "https://falcosecurity.github.io/charts"
	grafanaRepo    = "https://grafana.github.io/helm-charts"
	prometheusRepo = "https://prometheus-community.github.io/helm-charts"
)

func InstallCRD(ctx context.Context, kubeconfig, namespace string) error {
	baseURL := "https://raw.githubusercontent.com/prometheus-community/helm-charts/kube-prometheus-stack-" + prometheusVersion + "/charts/kube-prometheus-stack/crds/"

	crds := []string{
		baseURL + "crd-alertmanagerconfigs.yaml",
		baseURL + "crd-alertmanagers.yaml",
		baseURL + "crd-podmonitors.yaml",
		baseURL + "crd-probes.yaml",
		baseURL + "crd-prometheuses.yaml",
		baseURL + "crd-prometheusrules.yaml",
		baseURL + "crd-servicemonitors.yaml",
		baseURL + "crd-thanosrulers.yaml",
	}

	for _, crd := range crds {
		if err := kubectl.Invoke(ctx, []string{"apply", "-f", crd, "--validate=false", "--server-side=true", "--overwrite=true"}, kubectl.WithKubeconfig(kubeconfig), kubectl.WithNamespace(namespace), kubectl.WithDefaultOutput()); err != nil {
			return err
		}
	}

	return nil
}

func Install(ctx context.Context, kubeconfig, namespace string) error {
	if namespace == "" {
		namespace = "default"
	}

	if err := installLoki(ctx, kubeconfig, namespace); err != nil {
		return err
	}

	if err := installPromtail(ctx, kubeconfig, namespace); err != nil {
		return err
	}

	if err := installTempo(ctx, kubeconfig, namespace); err != nil {
		return err
	}

	if err := installPrometheus(ctx, kubeconfig, namespace); err != nil {
		return err
	}

	if err := installGrafana(ctx, kubeconfig, namespace); err != nil {
		return err
	}

	return nil
}

func Uninstall(ctx context.Context, kubeconfig, namespace string) error {
	if namespace == "" {
		namespace = "default"
	}

	if err := uninstallGrafana(ctx, kubeconfig, namespace); err != nil {
		// return err
	}

	if err := uninstallPrometheus(ctx, kubeconfig, namespace); err != nil {
		// return err
	}

	if err := uninstallTempo(ctx, kubeconfig, namespace); err != nil {
		// return err
	}

	if err := uninstallPromtail(ctx, kubeconfig, namespace); err != nil {
		// return err
	}

	if err := uninstallLoki(ctx, kubeconfig, namespace); err != nil {
		// return err
	}

	return nil
}
