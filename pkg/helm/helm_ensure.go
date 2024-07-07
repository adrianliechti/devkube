package helm

import (
	"context"

	"github.com/adrianliechti/loop/pkg/kubernetes"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Ensure(ctx context.Context, client kubernetes.Client, namespace, name, repoURL, chartName, chartVersion string, values map[string]any) error {
	if _, err := client.CoreV1().Namespaces().Get(ctx, namespace, metav1.GetOptions{}); err != nil {
		if !kubernetes.IsNotFound(err) {
			return err
		}

		obj := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: namespace,
			},
		}

		if _, err := client.CoreV1().Namespaces().Create(ctx, obj, metav1.CreateOptions{}); err != nil {
			return err
		}
	}

	err := Install(ctx, client, namespace, name, repoURL, chartName, chartVersion, values)

	if err != nil {
		err = Upgrade(ctx, client, namespace, name, repoURL, chartName, chartVersion, values)
	}

	return err
}
