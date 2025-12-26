package dashboard

import (
	"context"
	"fmt"
	"time"

	"github.com/adrianliechti/devkube/app"
	"github.com/adrianliechti/go-cli"
	"github.com/adrianliechti/loop/pkg/dashboard"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:  "dashboard",
		Usage: "open Dashboard in Browser",

		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := app.MustClient(ctx, cmd)

			port := app.MustPortOrRandom(ctx, cmd, 8888)
			url := fmt.Sprintf("http://127.0.0.1:%d", port)

			// ready := make(chan struct{})

			// go func() {
			// 	<-ready
			// 	cli.OpenURL(url)
			// }()

			time.AfterFunc(2*time.Second, func() {
				cli.OpenURL(url)
			})

			options := &dashboard.DashboardOptions{
				Port: port,

				PlatformNamespaces: []string{
					"kube-public",
					"kube-system",
					"kube-node-lease",
					"local-path-storage",

					"cert-manager",
					"crossplane-system",
					"gatekeeper-system",

					"argocd",
					"tekton-pipelines",
					"tekton-pipelines-resolvers",

					"platform",
				},
			}

			return dashboard.Run(ctx, client, options)
		},
	}
}
