package registry

import (
	"context"
	_ "embed"
	"strings"

	"github.com/adrianliechti/loop/pkg/kubernetes"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	namespace = "platform"
)

var (
	//go:embed manifest.yaml
	manifest string
)

func Ensure(ctx context.Context, client kubernetes.Client) error {
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

	if err := client.Apply(ctx, namespace, strings.NewReader(manifest)); err != nil {
		return err
	}

	return nil
}
