package metrics

import (
	"context"

	"github.com/adrianliechti/devkube/pkg/helm"
	"github.com/adrianliechti/loop/pkg/kubernetes"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	name      = "metrics-server"
	namespace = "kube-system"

	repoURL      = "https://kubernetes-sigs.github.io/metrics-server"
	chartName    = "metrics-server"
	chartVersion = "3.12.1"
)

func Ensure(ctx context.Context, client kubernetes.Client) error {
	if _, err := client.RbacV1().ClusterRoles().Get(ctx, "system:metrics-server", metav1.GetOptions{}); err == nil {
		return nil
	}

	values := map[string]any{
		"args": []string{
			"--kubelet-insecure-tls",
			"--kubelet-preferred-address-types=InternalIP",
		},
	}

	if err := helm.Ensure(ctx, client, namespace, name, repoURL, chartName, chartVersion, values); err != nil {
		return err
	}

	return nil
}
