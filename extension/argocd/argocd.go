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
	chartVersion = "7.5.0"
)

func Ensure(ctx context.Context, client kubernetes.Client) error {
	values := map[string]any{
		"crds": map[string]any{
			"enabled": true,
			"keep":    true,
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
