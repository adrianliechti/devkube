package helm

import (
	"context"
	"log/slog"
	"time"

	"github.com/adrianliechti/loop/pkg/kubernetes"

	"helm.sh/helm/v4/pkg/action"
	"helm.sh/helm/v4/pkg/chart/loader"
	"helm.sh/helm/v4/pkg/cli"
	"helm.sh/helm/v4/pkg/kube"
	"helm.sh/helm/v4/pkg/registry"
)

func Install(ctx context.Context, client kubernetes.Client, namespace, name, repoURL, chartName, chartVersion string, values map[string]any) error {
	settings := cli.New()

	config := new(action.Configuration)
	config.LogHolder.SetLogger(slog.DiscardHandler)

	if err := config.Init(NewClientGetter(client, namespace), namespace, ""); err != nil {
		return err
	}

	a := action.NewInstall(config)

	a.ReleaseName = name

	a.CreateNamespace = true
	a.Namespace = namespace

	a.RepoURL = repoURL
	a.Version = chartVersion

	//a.Wait = true
	a.Devel = true

	a.Timeout = 15 * time.Minute
	a.WaitStrategy = kube.StatusWatcherStrategy

	if client, err := registry.NewClient(); err == nil {
		a.SetRegistryClient(client)
	}

	path, err := a.ChartPathOptions.LocateChart(chartName, settings)

	if err != nil {
		return err
	}

	chart, err := loader.Load(path)

	if err != nil {
		return err
	}

	if _, err := a.Run(chart, values); err != nil {
		return err
	}

	return nil
}
