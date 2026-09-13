package helm

import (
	"context"

	"github.com/adrianliechti/loop/pkg/kubernetes"

	"helm.sh/helm/v4/pkg/action"
	helmkube "helm.sh/helm/v4/pkg/kube"
)

func Install(ctx context.Context, client kubernetes.Client, namespace, name, repoURL, chartName, chartVersion string, values map[string]any) error {
	config, err := newConfiguration(client, namespace)

	if err != nil {
		return err
	}

	a := action.NewInstall(config)

	a.ReleaseName = name

	a.CreateNamespace = true
	a.Namespace = namespace

	a.RepoURL = repoURL
	a.Version = chartVersion

	a.Devel = true

	a.Timeout = timeout
	a.WaitStrategy = helmkube.StatusWatcherStrategy

	setRegistryClient(a.SetRegistryClient)

	chart, err := loadChart(&a.ChartPathOptions, chartName)

	if err != nil {
		return err
	}

	if _, err := a.Run(chart, values); err != nil {
		return err
	}

	return nil
}
