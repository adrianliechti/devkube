package helm

import (
	"context"
	"time"

	"github.com/adrianliechti/loop/pkg/kubernetes"

	"helm.sh/helm/v3/pkg/action"
)

func Install(ctx context.Context, client kubernetes.Client, namespace, name, repoURL, chartName, chartVersion string, values map[string]any) error {
	chart, err := loadChart(repoURL, chartName, chartVersion)

	if err != nil {
		return err
	}

	config := new(action.Configuration)
	logger := func(format string, v ...interface{}) {}

	if err := config.Init(NewClientGetter(client), namespace, "", logger); err != nil {
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

	if _, err := a.Run(chart, values); err != nil {
		return err
	}

	return nil
}
