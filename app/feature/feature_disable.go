package feature

import (
	"errors"
	"strings"

	"github.com/adrianliechti/devkube/app"
	"github.com/adrianliechti/devkube/extension/falco"
	"github.com/adrianliechti/devkube/extension/trivy"
	"github.com/adrianliechti/devkube/extension/vault"
	"github.com/adrianliechti/devkube/pkg/cli"
	"github.com/adrianliechti/devkube/pkg/helm"
	"github.com/adrianliechti/devkube/pkg/kubectl"
)

func DisableCommand() *cli.Command {
	return &cli.Command{
		Name:  "disable",
		Usage: "Disable cluster feature",

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
			feature := c.Args().First()

			if feature == "" {
				return errors.New("feature name is required")
			}

			provider, cluster := app.MustCluster(c)

			kubeconfig, closer := app.MustClusterKubeconfig(c, provider, cluster)
			defer closer()

			switch strings.ToLower(feature) {
			case "trivy":
				if err := trivy.Uninstall(c.Context, kubeconfig, DefaultNamespace); err != nil {
					return err
				}

				return nil

			case "falco":
				if err := falco.Uninstall(c.Context, kubeconfig, DefaultNamespace); err != nil {
					return err
				}

				return nil

			case "vault":
				if err := vault.Uninstall(c.Context, kubeconfig, DefaultNamespace); err != nil {
					return err
				}

				return nil

			default:
				cli.Fatalf("inavlid feature: %s", feature)
			}

			return nil
		},
	}
}
