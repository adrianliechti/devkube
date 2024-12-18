package gatekeeper

import (
	"context"

	"github.com/adrianliechti/devkube/pkg/helm"
	"github.com/adrianliechti/loop/pkg/kubernetes"
)

const (
	name      = "gatekeeper"
	namespace = "gatekeeper-system"

	// https://artifacthub.io/packages/helm/gatekeeper/gatekeeper
	repoURL      = "https://open-policy-agent.github.io/gatekeeper/charts"
	chartName    = "gatekeeper"
	chartVersion = "3.18.1"
)

func Ensure(ctx context.Context, client kubernetes.Client) error {
	values := map[string]any{}

	if err := helm.Ensure(ctx, client, namespace, name, repoURL, chartName, chartVersion, values); err != nil {
		return err
	}

	return nil
}
