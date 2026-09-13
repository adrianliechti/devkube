package kube

import (
	"context"

	"github.com/adrianliechti/loop/pkg/kubernetes"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EnsureNamespace creates the given namespace unless it already exists.
func EnsureNamespace(ctx context.Context, client kubernetes.Client, name string) error {
	obj := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}

	if _, err := client.CoreV1().Namespaces().Create(ctx, obj, metav1.CreateOptions{}); err != nil {
		if !kubernetes.IsAlreadyExists(err) {
			return err
		}
	}

	return nil
}
