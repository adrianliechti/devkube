package helm

import (
	"context"

	"github.com/adrianliechti/devkube/pkg/kube"
	"github.com/adrianliechti/loop/pkg/kubernetes"

	"helm.sh/helm/v4/pkg/action"
	helmkube "helm.sh/helm/v4/pkg/kube"
)

func Upgrade(ctx context.Context, client kubernetes.Client, namespace, name, repoURL, chartName, chartVersion string, values map[string]any) error {
	config, err := newConfiguration(client, namespace)

	if err != nil {
		return err
	}

	// unlike install, upgrade does not create the namespace itself
	if err := kube.EnsureNamespace(ctx, client, namespace); err != nil {
		return err
	}

	a := action.NewUpgrade(config)

	a.Namespace = namespace

	a.RepoURL = repoURL
	a.Version = chartVersion

	a.ReuseValues = false
	a.ResetValues = true

	a.Devel = true

	a.CleanupOnFail = true

	a.Timeout = timeout
	a.WaitStrategy = helmkube.StatusWatcherStrategy

	setRegistryClient(a.SetRegistryClient)

	chart, err := loadChart(&a.ChartPathOptions, chartName)

	if err != nil {
		return err
	}

	if _, err := a.Run(name, chart, values); err != nil {
		return err
	}

	return nil
}
