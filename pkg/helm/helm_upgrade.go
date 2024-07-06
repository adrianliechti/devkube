package helm

import (
	"context"

	"github.com/adrianliechti/loop/pkg/kubernetes"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"helm.sh/helm/v3/pkg/action"
)

func Upgrade(ctx context.Context, client kubernetes.Client, namespace, name, repoURL, chartName, chartVersion string, values map[string]any) error {
	chart, err := loadChart(repoURL, chartName, chartVersion)

	if err != nil {
		return err
	}

	config := new(action.Configuration)
	logger := func(format string, v ...interface{}) {}

	if err := config.Init(NewClientGetter(client), namespace, "", logger); err != nil {
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

	if _, err := a.Run(name, chart, values); err != nil {
		return err
	}

	return nil
}
