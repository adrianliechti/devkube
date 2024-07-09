package helm

import (
	"context"
	"errors"

	"github.com/adrianliechti/loop/pkg/kubernetes"
)

func Ensure(ctx context.Context, client kubernetes.Client, namespace, name, repoURL, chartName, chartVersion string, values map[string]any) error {
	err := Upgrade(ctx, client, namespace, name, repoURL, chartName, chartVersion, values)

	if errors.Is(err, ErrNoDeployedReleases) {
		err = Install(ctx, client, namespace, name, repoURL, chartName, chartVersion, values)
	}

	return err
}
