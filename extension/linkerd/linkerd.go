package linkerd

import (
	"context"
	_ "embed"
	"strings"

	"github.com/adrianliechti/devkube/pkg/helm"
	"github.com/adrianliechti/devkube/pkg/kubectl"
	"github.com/adrianliechti/devkube/pkg/kubernetes"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	linkerdRepo = "https://helm.linkerd.io/stable"

	crds        = "linkerd-crds"
	crdsChart   = "linkerd-crds"
	crdsVersion = "1.4.0"

	linkerd        = "linkerd-control-plane"
	linkerdChart   = "linkerd-control-plane"
	linkerdVersion = "1.9.3"

	viz        = "linkerd-viz"
	vizChart   = "linkerd-viz"
	vizVersion = "30.3.3"
)

var (
	//go:embed linkerd.yaml
	manifest string
)

func Install(ctx context.Context, kubeconfig, namespace string) error {
	if namespace == "" {
		namespace = "default"
	}

	client, err := kubernetes.NewFromConfig(kubeconfig)

	if err != nil {
		return err
	}

	if err := kubectl.Invoke(ctx, []string{"apply", "-f", "-"}, kubectl.WithKubeconfig(kubeconfig), kubectl.WithNamespace(namespace), kubectl.WithInput(strings.NewReader(manifest)), kubectl.WithDefaultOutput()); err != nil {
		return err
	}

	secret, err := client.CoreV1().Secrets(namespace).Get(ctx, "linkerd-identity-issuer", metav1.GetOptions{})

	if err != nil {
		return err
	}

	configmap := corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: "linkerd-identity-trust-roots",

			Labels: map[string]string{
				"linkerd.io/control-plane-component": "identity",
				"linkerd.io/control-plane-ns":        namespace,
			},
		},
		Data: map[string]string{
			"ca-bundle.crt": string(secret.Data["ca.crt"]),
		},
	}

	client.CoreV1().ConfigMaps(namespace).Delete(ctx, "linkerd-identity-trust-roots", metav1.DeleteOptions{})

	if _, err := client.CoreV1().ConfigMaps(namespace).Create(ctx, &configmap, metav1.CreateOptions{}); err != nil {
		return err
	}

	if err := helm.Install(ctx, crds, linkerdRepo, crdsChart, crdsVersion, nil, helm.WithKubeconfig(kubeconfig), helm.WithNamespace(namespace), helm.WithDefaultOutput()); err != nil {
		return err
	}

	cpvalues := map[string]any{
		"identity": map[string]any{
			"externalCA": true,

			"issuer": map[string]any{
				"scheme": "kubernetes.io/tls",
			},
		},
	}

	if err := helm.Install(ctx, linkerd, linkerdRepo, linkerdChart, linkerdVersion, cpvalues, helm.WithKubeconfig(kubeconfig), helm.WithNamespace(namespace), helm.WithDefaultOutput()); err != nil {
		return err
	}

	vizvalues := map[string]any{
		"linkerdNamespace": namespace,

		"jaegerUrl":     "loki.loop:16686",
		"prometheusUrl": "monitoring-prometheus.loop:9090",

		"grafana": map[string]any{
			"url": "grafana.loop",
		},

		"prometheus": map[string]any{
			"enabled": false,
		},
	}

	if err := helm.Install(ctx, viz, linkerdRepo, vizChart, vizVersion, vizvalues, helm.WithKubeconfig(kubeconfig), helm.WithNamespace(namespace), helm.WithDefaultOutput()); err != nil {
		return err
	}

	return nil
}

func Uninstall(ctx context.Context, kubeconfig, namespace string) error {
	if namespace == "" {
		namespace = "default"
	}

	client, err := kubernetes.NewFromConfig(kubeconfig)

	if err != nil {
		return err
	}

	if err := helm.Uninstall(ctx, linkerd, helm.WithKubeconfig(kubeconfig), helm.WithNamespace(namespace), helm.WithDefaultOutput()); err != nil {
		// return err
	}

	if err := helm.Uninstall(ctx, crds, helm.WithKubeconfig(kubeconfig), helm.WithNamespace(namespace), helm.WithDefaultOutput()); err != nil {
		// return err
	}

	if err := client.CoreV1().ConfigMaps(namespace).Delete(ctx, "linkerd-identity-trust-roots", metav1.DeleteOptions{}); err != nil {
		// return err
	}

	if err := kubectl.Invoke(ctx, []string{"delete", "-f", "-"}, kubectl.WithKubeconfig(kubeconfig), kubectl.WithNamespace(namespace), kubectl.WithInput(strings.NewReader(manifest)), kubectl.WithDefaultOutput()); err != nil {
		// return err
	}

	return nil
}
