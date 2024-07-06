package helm

import (
	"context"
	"errors"

	"github.com/adrianliechti/loop/pkg/kubernetes"
	"helm.sh/helm/v3/pkg/storage/driver"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Ensure(ctx context.Context, client kubernetes.Client, namespace, name, repoURL, chartName, chartVersion string, values map[string]any) error {
	client.CoreV1().Namespaces().Create(ctx, &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}, metav1.CreateOptions{})

	err := Upgrade(ctx, client, namespace, name, repoURL, chartName, chartVersion, values)

	if errors.Is(err, driver.ErrNoDeployedReleases) {
		err = Install(ctx, client, namespace, name, repoURL, chartName, chartVersion, values)
	}

	return err
}
