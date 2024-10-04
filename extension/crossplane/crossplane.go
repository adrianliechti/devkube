package crossplane

import (
	"context"

	"github.com/adrianliechti/devkube/pkg/helm"
	"github.com/adrianliechti/loop/pkg/kubernetes"
)

const (
	name      = "crossplane"
	namespace = "crossplane-system"

	//
	// https://artifacthub.io/packages/helm/crossplane/crossplane
	// https://github.com/crossplane/crossplane/releases
	repoURL      = "https://charts.crossplane.io/stable"
	chartName    = "crossplane"
	chartVersion = "1.17.1"
)

func Ensure(ctx context.Context, client kubernetes.Client) error {
	values := map[string]any{
		"args": []string{
			"--enable-environment-configs",
		},
	}

	if err := helm.Ensure(ctx, client, namespace, name, repoURL, chartName, chartVersion, values); err != nil {
		return err
	}

	return nil
}
