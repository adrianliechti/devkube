package logs

import (
	"context"
	"os"

	"github.com/adrianliechti/devkube/app"
	"github.com/adrianliechti/devkube/pkg/cli"
	"github.com/adrianliechti/loop/pkg/kubernetes"
	"github.com/adrianliechti/loop/pkg/kubernetes/resource"

	corev1 "k8s.io/api/core/v1"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:  "logs",
		Usage: "stream application logs",

		Action: func(c *cli.Context) error {
			client := app.MustClient(c)

			name := ""
			namespace := ""

			if name == "" {
				app := app.MustApplication(c.Context, client, namespace)

				name = app.Name
				namespace = app.Namespace
			}

			return Stream(c.Context, client, namespace, name)
		},
	}
}

func Stream(ctx context.Context, client kubernetes.Client, namespace, name string) error {
	app, err := resource.App(ctx, client, namespace, name)

	if err != nil {
		return err
	}

	for _, r := range app.Resources {
		if pod, ok := r.Object.(corev1.Pod); ok {
			for _, container := range pod.Spec.Containers {
				go client.PodLogs(ctx, pod.Namespace, pod.Name, container.Name, os.Stdout, true)
			}
		}
	}

	<-ctx.Done()
	return nil
}
