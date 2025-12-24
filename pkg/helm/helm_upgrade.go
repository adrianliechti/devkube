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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Upgrade(ctx context.Context, client kubernetes.Client, namespace, name, repoURL, chartName, chartVersion string, values map[string]any) error {
	settings := cli.New()

	config := new(action.Configuration)
	config.LogHolder.SetLogger(slog.DiscardHandler)

	if err := config.Init(NewClientGetter(client, namespace), namespace, ""); err != nil {
		return err
	}

	client.CoreV1().Namespaces().Create(ctx, &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}, metav1.CreateOptions{})

	a := action.NewUpgrade(config)

	a.Namespace = namespace

	a.RepoURL = repoURL
	a.Version = chartVersion

	a.ReuseValues = false
	a.ResetValues = true

	//a.Wait = true
	a.Devel = true

	a.CleanupOnFail = true

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

	if _, err := a.Run(name, chart, values); err != nil {
		return err
	}

	return nil
}
