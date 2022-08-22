package observability

import (
	"context"
)

const (
	falcoRepo      = "https://falcosecurity.github.io/charts"
	grafanaRepo    = "https://grafana.github.io/helm-charts"
	prometheusRepo = "https://prometheus-community.github.io/helm-charts"
)

func manifests() []string {
	result := []string{
		"crd-alertmanagerconfigs.yaml",
		"crd-alertmanagers.yaml",
		"crd-podmonitors.yaml",
		"crd-probes.yaml",
		"crd-prometheuses.yaml",
		"crd-prometheusrules.yaml",
		"crd-servicemonitors.yaml",
		"crd-thanosrulers.yaml",
	}

	for i := range result {
		name := result[i]
		result[i] = "https://raw.githubusercontent.com/prometheus-community/helm-charts/kube-prometheus-stack-" + prometheusVersion + "/charts/kube-prometheus-stack/crds/" + name
	}

	return result
}

// func InstallCRD(ctx context.Context, kubeconfig, namespace string) error {
// 	for _, manifest := range manifests() {
// 		if err := kubectl.Invoke(ctx, []string{"apply", "-f", manifest, "--validate=false", "--server-side=true", "--overwrite=true"}, kubectl.WithKubeconfig(kubeconfig), kubectl.WithNamespace(namespace), kubectl.WithDefaultOutput()); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// func UninstallCRD(ctx context.Context, kubeconfig, namespace string) error {
// 	for _, manifest := range manifests() {
// 		if err := kubectl.Invoke(ctx, []string{"delete", "-f", manifest}, kubectl.WithKubeconfig(kubeconfig)); err != nil {
// 			// return err
// 		}
// 	}

// 	return nil
// }

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
		//return err
	}

	if err := uninstallPrometheus(ctx, kubeconfig, namespace); err != nil {
		//return err
	}

	if err := uninstallTempo(ctx, kubeconfig, namespace); err != nil {
		//return err
	}

	if err := uninstallPromtail(ctx, kubeconfig, namespace); err != nil {
		//return err
	}

	if err := uninstallLoki(ctx, kubeconfig, namespace); err != nil {
		//return err
	}

	return nil
}
