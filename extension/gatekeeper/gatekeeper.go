package gatekeeper

import (
	"context"

	"github.com/adrianliechti/devkube/pkg/apply"
	"github.com/adrianliechti/loop/pkg/kubernetes"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	namespace = "gatekeeper-system"

	// https://github.com/open-policy-agent/gatekeeper/releases
	version = "v3.16.3"
)

func Ensure(ctx context.Context, client kubernetes.Client) error {
	if err := apply.ApplyURL(ctx, client, namespace, "https://raw.githubusercontent.com/open-policy-agent/gatekeeper/"+version+"/deploy/gatekeeper.yaml"); err != nil {
		return err
	}

	list, _ := client.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})

	for _, pod := range list.Items {
		client.WaitForPod(ctx, pod.Namespace, pod.Name)
	}

	return nil
}
