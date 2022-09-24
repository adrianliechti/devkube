package cluster

import (
	"fmt"
	"runtime"

	"github.com/adrianliechti/devkube/app"
	"github.com/adrianliechti/devkube/pkg/cli"
	"github.com/adrianliechti/devkube/pkg/kubectl"
)

func RegistryCommand() *cli.Command {
	return &cli.Command{
		Name:  "registry",
		Usage: "Connect Registry",

		Flags: []cli.Flag{
			app.ProviderFlag,
			app.ClusterFlag,
			app.PortFlag,
		},

		Before: func(c *cli.Context) error {
			if _, _, err := kubectl.Info(c.Context); err != nil {
				return err
			}

			return nil
		},

		Action: func(c *cli.Context) error {
			provider, cluster := app.MustCluster(c)

			kubeconfig, closer := app.MustClusterKubeconfig(c, provider, cluster)
			defer closer()

			port := 5000

			if runtime.GOOS == "darwin" {
				port = 5001
			}

			port = app.MustPortOrRandom(c, port)

			cli.Info("Configure Docker to use this registry")
			cli.Info("  {")
			cli.Info("    ...")
			cli.Info("    \"insecure-registries\": [")
			cli.Infof("      \"localhost:%d\",", port)
			cli.Infof("      \"host.docker.internal:%d\"", port)
			cli.Info("    ]")
			cli.Info("    ...")
			cli.Info("  }")
			cli.Info("  (see https://docs.docker.com/registry/insecure/#deploy-a-plain-http-registry)")
			cli.Info()
			cli.Info()
			cli.Info("Push an image")
			cli.Infof("  docker tag my-image host.docker.internal:%d/my-image", port)
			cli.Infof("  docker push host.docker.internal:%d/my-image", port)
			cli.Info()
			cli.Info("Push an image (not using BuildKit)")
			cli.Infof("  docker tag my-image localhost:%d/my-image", port)
			cli.Infof("  docker push localhost:%d/my-image", port)
			cli.Info()
			cli.Info()

			if err := kubectl.Invoke(c.Context, []string{"port-forward", "service/registry", fmt.Sprintf("%d:80", port)}, kubectl.WithKubeconfig(kubeconfig), kubectl.WithNamespace(DefaultNamespace), kubectl.WithDefaultOutput()); err != nil {
				return err
			}

			return nil
		},
	}
}
