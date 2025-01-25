package certmanager

import (
	"context"
	_ "embed"
	"strings"

	"github.com/adrianliechti/devkube/pkg/helm"
	"github.com/adrianliechti/loop/pkg/kubernetes"
)

var (
	//go:embed manifest.yaml
	manifest string
)

const (
	name      = "cert-manager"
	namespace = "cert-manager"

	// https://artifacthub.io/packages/helm/cert-manager/cert-manager
	repoURL      = "https://charts.jetstack.io"
	chartName    = "cert-manager"
	chartVersion = "1.16.3"
)

func Ensure(ctx context.Context, client kubernetes.Client) error {
	values := map[string]any{
		"crds": map[string]any{
			"enabled": true,
		},
	}

	if err := helm.Ensure(ctx, client, namespace, name, repoURL, chartName, chartVersion, values); err != nil {
		return err
	}

	if err := client.Apply(ctx, namespace, strings.NewReader(manifest)); err != nil {
		return err
	}

	return nil
}
