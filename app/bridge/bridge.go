package bridge

import (
	"context"
	"fmt"
	"time"

	"github.com/adrianliechti/devkube/app"
	"github.com/adrianliechti/go-cli"

	"github.com/adrianliechti/bridge/pkg/config"
	"github.com/adrianliechti/bridge/pkg/server"
)

// namespaces bridge should present as platform infrastructure
// instead of user workload
var platformNamespaces = []string{
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
}

func Command() *cli.Command {
	return &cli.Command{
		Name:  "bridge",
		Usage: "open Bridge Kubernetes dashboard",

		Flags: []cli.Flag{
			app.PortFlag,
		},

		Action: func(ctx context.Context, cmd *cli.Command) error {
			// bridge reads the kubeconfig itself and serves every context in it
			cfg, err := config.New()

			if err != nil {
				return err
			}

			if cfg.Kubernetes != nil {
				cfg.Kubernetes.PlatformNamespaces = platformNamespaces
			}

			srv, err := server.New(cfg)

			if err != nil {
				return err
			}

			port := app.MustPortOrRandom(ctx, cmd, 8888)

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
