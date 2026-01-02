package bridge

import (
	"context"
	"fmt"
	"time"

	"github.com/adrianliechti/devkube/app"
	"github.com/adrianliechti/go-cli"
	"github.com/adrianliechti/loop/pkg/bridge"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:  "bridge",
		Usage: "open Bridge Kubernetes dashboard",

		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := app.MustClient(ctx, cmd)

			platform := &bridge.PlatformConfig{
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

			port := app.MustPortOrRandom(ctx, cmd, 8888)

			srv, err := bridge.New(client, platform)

			if err != nil {
				return err
			}

			url := fmt.Sprintf("http://localhost:%d", port)
			addr := fmt.Sprintf("localhost:%d", port)

			time.AfterFunc(500*time.Millisecond, func() {
				cli.Infof("Bridge on %s", url)
				cli.OpenURL(url)
			})

			return srv.ListenAndServe(ctx, addr)
		},
	}
}
