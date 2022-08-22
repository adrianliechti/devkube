package feature

import (
	"strings"

	"github.com/adrianliechti/devkube/app"
	"github.com/adrianliechti/devkube/extension/trivy"
	"github.com/adrianliechti/devkube/pkg/cli"
	"github.com/adrianliechti/devkube/pkg/helm"
	"github.com/adrianliechti/devkube/pkg/kubectl"
)

func EnableCommand() *cli.Command {
	return &cli.Command{
		Name:  "enable",
		Usage: "Enable cluster feature",

		Flags: []cli.Flag{
			app.ProviderFlag,
			app.ClusterFlag,
		},

		Before: func(c *cli.Context) error {
			if _, _, err := helm.Info(c.Context); err != nil {
				return err
			}

			if _, _, err := kubectl.Info(c.Context); err != nil {
				return err
			}

			return nil
		},

		Action: func(c *cli.Context) error {
			provider, cluster := app.MustCluster(c)

			kubeconfig, closer := app.MustClusterKubeconfig(c, provider, cluster)
			defer closer()

			feature := c.Args().First()

			switch strings.ToLower(feature) {
			case "trivy":
				if err := trivy.Install(c.Context, kubeconfig, DefaultNamespace); err != nil {
					return err
				}

				return nil
			default:
				cli.Fatal("inavlid feature: %s", feature)
			}

			return nil
		},
	}
}
