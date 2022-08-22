package helm

import (
	"context"
	"fmt"
	"io"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/kube"

	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func (h *Helm) Install(ctx context.Context, release, repo, chart, version string, values map[string]interface{}) error {
	namespace := h.namespace

	if namespace == "" {
		namespace = "default"
	}

	if version == "" {
		version = ">0.0.0-0"
	}

	config := new(action.Configuration)
	settings := cli.New()

	log := func(format string, v ...interface{}) {
		//log.Printf(format, v)
	}

	if err := config.Init(kube.GetConfig(h.kubeconfig, h.context, namespace), namespace, "", log); err != nil {
		return err
	}

	client := action.NewInstall(config)

	client.ReleaseName = release

	client.Namespace = namespace
	client.CreateNamespace = true

	client.RepoURL = repo
	client.Version = version

	chartPath, err := client.LocateChart(chart, settings)

	if err != nil {
		return err
	}

	chartRequested, err := loader.Load(chartPath)

	if err != nil {
		return err
	}

	if chartRequested.Metadata.Dependencies != nil {
		if err := action.CheckDependencies(chartRequested, chartRequested.Metadata.Dependencies); err != nil {
			downloader := &downloader.Manager{
				Out:              io.Discard,
				ChartPath:        chartPath,
				Keyring:          client.ChartPathOptions.Keyring,
				SkipUpdate:       false,
				Getters:          getter.All(settings),
				RepositoryConfig: settings.RepositoryConfig,
				RepositoryCache:  settings.RepositoryCache,
				Debug:            settings.Debug,
			}

			if err := downloader.Update(); err != nil {
				return err
			}

			if chartRequested, err = loader.Load(chartPath); err != nil {
				return err
			}
		}
	}

	result, err := client.Run(chartRequested, values)

	if err != nil {
		return err
	}

	fmt.Println("Successfully installed release", result.Name)

	return nil
}
