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
	grafanaRepo = "https://grafana.github.io/helm-charts"

	crds        = "linkerd-crds"
	crdsChart   = "linkerd-crds"
	crdsVersion = "1.4.0"

	linkerd          = "linkerd-control-plane"
	linkerdChart     = "linkerd-control-plane"
	linkerdVersion   = "1.9.3"
	linkerdNamespace = "linkerd"

	smi          = "linkerd-smi"
	smiChart     = "linkerd-smi"
	smiVersion   = "1.0.0"
	smiNamespace = "linkerd-smi"

	viz          = "linkerd-viz"
	vizChart     = "linkerd-viz"
	vizVersion   = "30.3.3"
	vizNamespace = "linkerd-viz"

	jaeger          = "linkerd-jaeger"
	jaegerChart     = "linkerd-jaeger"
	jaegerVersion   = "30.4.3"
	jaegerNamespace = "linkerd-jaeger"

	grafana        = "grafana"
	grafanaChart   = "grafana"
	grafanaVersion = "6.40.0"
)

var (
	//go:embed linkerd.yaml
	manifest string
)

func Install(ctx context.Context, kubeconfig string) error {
	kubectl.Invoke(ctx, []string{"create", "namespace", linkerdNamespace}, kubectl.WithKubeconfig(kubeconfig))

	if err := installCRDs(ctx, kubeconfig); err != nil {
		return err
	}

	if err := installCA(ctx, kubeconfig); err != nil {
		return err
	}

	if err := installLinkerd(ctx, kubeconfig); err != nil {
		return err
	}

	if err := installSMI(ctx, kubeconfig); err != nil {
		return err
	}

	if err := installJaeger(ctx, kubeconfig); err != nil {
		return err
	}

	if err := installGrafana(ctx, kubeconfig); err != nil {
		return err
	}

	if err := installViz(ctx, kubeconfig); err != nil {
		return err
	}

	return nil
}

func Uninstall(ctx context.Context, kubeconfig string) error {
	if err := uninstallViz(ctx, kubeconfig); err != nil {
		// return err
	}

	if err := uninstallGrafana(ctx, kubeconfig); err != nil {
		// return err
	}

	if err := uninstallJaeger(ctx, kubeconfig); err != nil {
		// return err
	}

	if err := uninstallSMI(ctx, kubeconfig); err != nil {
		// return err
	}

	if err := uninstallLinkerd(ctx, kubeconfig); err != nil {
		// return err
	}

	if err := uninstallCA(ctx, kubeconfig); err != nil {
		// return err
	}

	if err := uninstallCRDs(ctx, kubeconfig); err != nil {
		// return err
	}

	return nil
}

func installCRDs(ctx context.Context, kubeconfig string) error {
	if err := helm.Install(ctx, crds, linkerdRepo, crdsChart, crdsVersion, nil, helm.WithKubeconfig(kubeconfig), helm.WithNamespace(linkerdNamespace), helm.WithDefaultOutput()); err != nil {
		return err
	}

	return nil
}

func uninstallCRDs(ctx context.Context, kubeconfig string) error {
	if err := helm.Uninstall(ctx, crds, helm.WithKubeconfig(kubeconfig), helm.WithNamespace(linkerdNamespace)); err != nil {
		// return err
	}

	return nil
}

func installCA(ctx context.Context, kubeconfig string) error {
	client, err := kubernetes.NewFromConfig(kubeconfig)

	if err != nil {
		return err
	}

	client.CoreV1().Secrets(linkerdNamespace).Delete(ctx, "linkerd-identity-issuer", metav1.DeleteOptions{})
	client.CoreV1().ConfigMaps(linkerdNamespace).Delete(ctx, "linkerd-identity-trust-roots", metav1.DeleteOptions{})

	if err := kubectl.Invoke(ctx, []string{"apply", "-f", "-"}, kubectl.WithKubeconfig(kubeconfig), kubectl.WithNamespace(linkerdNamespace), kubectl.WithInput(strings.NewReader(manifest)), kubectl.WithDefaultOutput()); err != nil {
		return err
	}

	if err := kubectl.Invoke(ctx, []string{"wait", "--for=condition=Ready", "certificate/linkerd-identity-issuer"}, kubectl.WithKubeconfig(kubeconfig), kubectl.WithNamespace(linkerdNamespace)); err != nil {
		return err
	}

	secret, err := client.CoreV1().Secrets(linkerdNamespace).Get(ctx, "linkerd-identity-issuer", metav1.GetOptions{})

	if err != nil {
		return err
	}

	if _, err := client.CoreV1().ConfigMaps(linkerdNamespace).Create(ctx, &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: "linkerd-identity-trust-roots",

			Labels: map[string]string{
				"linkerd.io/control-plane-component": "identity",
				"linkerd.io/control-plane-ns":        linkerdNamespace,
			},
		},
		Data: map[string]string{
			"ca-bundle.crt": string(secret.Data["ca.crt"]),
		},
	}, metav1.CreateOptions{}); err != nil {
		return err
	}

	return nil
}

func uninstallCA(ctx context.Context, kubeconfig string) error {
	client, err := kubernetes.NewFromConfig(kubeconfig)

	if err != nil {
		return err
	}

	if err := kubectl.Invoke(ctx, []string{"delete", "-f", "-"}, kubectl.WithKubeconfig(kubeconfig), kubectl.WithNamespace(linkerdNamespace), kubectl.WithInput(strings.NewReader(manifest))); err != nil {
		// return err
	}

	if err := client.CoreV1().ConfigMaps(linkerdNamespace).Delete(ctx, "linkerd-identity-trust-roots", metav1.DeleteOptions{}); err != nil {
		// return err
	}

	if err := client.CoreV1().Secrets(linkerdNamespace).Delete(ctx, "linkerd-identity-issuer", metav1.DeleteOptions{}); err != nil {
		// return err
	}

	return nil
}

func installLinkerd(ctx context.Context, kubeconfig string) error {
	values := map[string]any{
		"identity": map[string]any{
			"externalCA": true,

			"issuer": map[string]any{
				"scheme": "kubernetes.io/tls",
			},
		},
	}

	if err := helm.Install(ctx, linkerd, linkerdRepo, linkerdChart, linkerdVersion, values, helm.WithKubeconfig(kubeconfig), helm.WithNamespace(linkerdNamespace), helm.WithWait(true), helm.WithDefaultOutput()); err != nil {
		return err
	}

	return nil
}

func uninstallLinkerd(ctx context.Context, kubeconfig string) error {
	if err := helm.Uninstall(ctx, linkerd, helm.WithKubeconfig(kubeconfig), helm.WithNamespace(linkerdNamespace)); err != nil {
		// return err
	}

	kubectl.Invoke(ctx, []string{"delete", "namespace", linkerdNamespace}, kubectl.WithKubeconfig(kubeconfig))

	return nil
}

func installSMI(ctx context.Context, kubeconfig string) error {
	values := map[string]any{}

	if err := helm.Install(ctx, smi, linkerdRepo, smiChart, smiVersion, values, helm.WithKubeconfig(kubeconfig), helm.WithNamespace(smiNamespace), helm.WithDefaultOutput()); err != nil {
		return err
	}

	return nil
}

func uninstallSMI(ctx context.Context, kubeconfig string) error {
	if err := helm.Uninstall(ctx, smi, helm.WithKubeconfig(kubeconfig), helm.WithNamespace(smiNamespace)); err != nil {
		// return err
	}

	kubectl.Invoke(ctx, []string{"delete", "namespace", smiNamespace}, kubectl.WithKubeconfig(kubeconfig))

	return nil
}

func installJaeger(ctx context.Context, kubeconfig string) error {
	values := map[string]any{}

	if err := helm.Install(ctx, jaeger, linkerdRepo, jaegerChart, jaegerVersion, values, helm.WithKubeconfig(kubeconfig), helm.WithNamespace(jaegerNamespace), helm.WithDefaultOutput()); err != nil {
		return err
	}

	return nil
}

func uninstallJaeger(ctx context.Context, kubeconfig string) error {
	if err := helm.Uninstall(ctx, jaeger, helm.WithKubeconfig(kubeconfig), helm.WithNamespace(jaegerNamespace)); err != nil {
		// return err
	}

	kubectl.Invoke(ctx, []string{"delete", "namespace", jaegerNamespace}, kubectl.WithKubeconfig(kubeconfig))

	return nil
}

func installGrafana(ctx context.Context, kubeconfig string) error {
	h := helm.New(helm.WithKubeconfig(kubeconfig), helm.WithNamespace(vizNamespace), helm.WithDefaultOutput())

	args := []string{
		"upgrade", "--install", "--create-namespace", "--timeout", "10m0s",
		grafana,
		grafanaChart,
		"--repo", grafanaRepo,
		"--version", grafanaVersion,
		"--set", "rbac.create=false",
		"--set", "service.port=3000",
		"--set", "serviceAccount.create=false",
		"-f", "https://raw.githubusercontent.com/linkerd/linkerd2/main/grafana/values.yaml",
	}

	if err := h.Invoke(ctx, args...); err != nil {
		return err
	}

	return nil
}

func uninstallGrafana(ctx context.Context, kubeconfig string) error {
	if err := helm.Uninstall(ctx, grafana, helm.WithKubeconfig(kubeconfig), helm.WithNamespace(vizNamespace)); err != nil {
		// return err
	}

	return nil
}

func installViz(ctx context.Context, kubeconfig string) error {
	values := map[string]any{
		"jaegerUrl": "jaeger." + jaegerNamespace + ":16686",

		"grafana": map[string]any{
			"url": "grafana." + vizNamespace + ":3000",
		},
	}

	if err := helm.Install(ctx, viz, linkerdRepo, vizChart, vizVersion, values, helm.WithKubeconfig(kubeconfig), helm.WithNamespace(vizNamespace), helm.WithDefaultOutput()); err != nil {
		return err
	}

	return nil
}

func uninstallViz(ctx context.Context, kubeconfig string) error {
	if err := helm.Uninstall(ctx, viz, helm.WithKubeconfig(kubeconfig), helm.WithNamespace(vizNamespace)); err != nil {
		// return err
	}

	kubectl.Invoke(ctx, []string{"delete", "namespace", vizNamespace}, kubectl.WithKubeconfig(kubeconfig))

	return nil
}
