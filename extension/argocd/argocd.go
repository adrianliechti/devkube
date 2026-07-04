package argocd

import (
	"context"

	"github.com/adrianliechti/devkube/pkg/helm"
	"github.com/adrianliechti/loop/pkg/kubernetes"
)

const (
	name      = "argocd"
	namespace = "argocd"

	// https://artifacthub.io/packages/helm/argo/argo-cd
	repoURL      = "https://argoproj.github.io/argo-helm"
	chartName    = "argo-cd"
	chartVersion = "10.1.2"
)

func Ensure(ctx context.Context, client kubernetes.Client) error {
	values := map[string]any{
		// since chart 10.0.0 network policies are created by default,
		// blocking cross-namespace access in this dev cluster
		"global": map[string]any{
			"networkPolicy": map[string]any{
				"create": false,
			},
		},

		"crds": map[string]any{
			"install": true,
			"keep":    true,
		},

		"dex": map[string]any{
			"enabled": false,
		},

		"configs": map[string]any{
			"cm": map[string]any{
				"exec.enabled": true,
			},

			"params": map[string]any{
				"server.insecure":              true,
				"server.disable.auth":          true,
				"server.repo.server.plaintext": true,

				"application.namespaces": "*",

				"dexserver.disable.tls":       true,
				"server.dex.server.plaintext": true,

				"reposerver.disable.tls":           true,
				"controller.repo.server.plaintext": true,
			},
		},
	}

	if err := helm.Ensure(ctx, client, namespace, name, repoURL, chartName, chartVersion, values); err != nil {
		return err
	}

	return nil
}
