package certmanager

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
	certmanagerRepo      = "https://charts.jetstack.io"
	certmanagerNamespace = "cert-manager"

	certmanager        = "cert-manager"
	certmanagerChart   = "cert-manager"
	certmanagerVersion = "v1.9.1"
)

var (
	//go:embed manifest.yaml
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

	nodes, err := client.CoreV1().Nodes().List(ctx, metav1.ListOptions{})

	if err != nil {
		return err
	}

	values := map[string]any{
		"installCRDs": true,

		"prometheus": map[string]any{
			"servicemonitor": map[string]any{
				"enabled": true,
			},
		},
	}

	if isAWS(nodes.Items) {
		values["webhook"] = map[string]any{
			"securePort":  10260,
			"hostNetwork": true,
		}
	}

	if err := helm.Install(ctx, certmanager, certmanagerRepo, certmanagerChart, certmanagerVersion, values, helm.WithKubeconfig(kubeconfig), helm.WithNamespace(certmanagerNamespace), helm.WithWait(true), helm.WithDefaultOutput()); err != nil {
		return err
	}

	if err := kubectl.Invoke(ctx, []string{"apply", "-f", "-"}, kubectl.WithKubeconfig(kubeconfig), kubectl.WithNamespace(namespace), kubectl.WithInput(strings.NewReader(manifest)), kubectl.WithDefaultOutput()); err != nil {
		return err
	}

	return nil
}

func Uninstall(ctx context.Context, kubeconfig, namespace string) error {
	if namespace == "" {
		namespace = "default"
	}

	if err := kubectl.Invoke(ctx, []string{"delete", "-f", "-"}, kubectl.WithKubeconfig(kubeconfig), kubectl.WithNamespace(namespace), kubectl.WithInput(strings.NewReader(manifest)), kubectl.WithDefaultOutput()); err != nil {
		// return err
	}

	if err := helm.Uninstall(ctx, certmanager, helm.WithKubeconfig(kubeconfig), helm.WithNamespace(certmanagerNamespace)); err != nil {
		// return err
	}

	return nil
}

func isAWS(nodes []corev1.Node) bool {
	for _, node := range nodes {
		if strings.HasPrefix(node.Spec.ProviderID, "aws://") {
			return true
		}
	}

	return false
}
