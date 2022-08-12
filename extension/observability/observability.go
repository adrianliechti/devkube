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

var (
	Images = []string{
		// "grafana/loki:2.5.0",

		// "grafana/tempo:1.4.0",
		// "grafana/tempo-query:1.4.0",

		// "grafana/promtail:2.5.0",

		// "quay.io/prometheus/prometheus:v2.34.0",
		// "quay.io/prometheus/node-exporter:v1.3.0",
		// "k8s.gcr.io/kube-state-metrics/kube-state-metrics:v2.4.1",

		// "grafana/grafana:8.5.0",
	}
)

func InstallCRD(ctx context.Context, kubeconfig, namespace string) error {
	crds := []string{
		"crd-alertmanagerconfigs.yaml",
		"crd-alertmanagers.yaml",
		"crd-podmonitors.yaml",
		"crd-probes.yaml",
		"crd-prometheuses.yaml",
		"crd-prometheusrules.yaml",
		"crd-servicemonitors.yaml",
		"crd-thanosrulers.yaml",
	}

	for _, crd := range crds {
		url := "https://raw.githubusercontent.com/prometheus-community/helm-charts/kube-prometheus-stack-" + prometheusVersion + "/charts/kube-prometheus-stack/crds/" + crd

		if err := kubectl.Invoke(ctx, kubeconfig, "apply", "-n", namespace, "-f", url, "--validate=false", "--server-side=true", "--overwrite=true"); err != nil {
			return err
		}
	}

	return nil
}

func Install(ctx context.Context, kubeconfig, namespace string) error {
	if namespace == "" {
		namespace = "default"
	}

	if err := installPrometheus(ctx, kubeconfig, namespace); err != nil {
		return err
	}

	if err := installLoki(ctx, kubeconfig, namespace); err != nil {
		return err
	}

	if err := installTempo(ctx, kubeconfig, namespace); err != nil {
		return err
	}

	if err := installPromtail(ctx, kubeconfig, namespace); err != nil {
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
		//return err
	}

	if err := uninstallPrometheus(ctx, kubeconfig, namespace); err != nil {
		//return err
	}

	if err := uninstallPromtail(ctx, kubeconfig, namespace); err != nil {
		//return err
	}

	if err := uninstallTempo(ctx, kubeconfig, namespace); err != nil {
		//return err
	}

	if err := uninstallLoki(ctx, kubeconfig, namespace); err != nil {
		//return err
	}

	return nil
}
