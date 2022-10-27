package feature

import (
	"errors"
	"strings"

	"github.com/adrianliechti/devkube/app"

	"github.com/adrianliechti/devkube/pkg/cli"
	"github.com/adrianliechti/devkube/pkg/helm"
	"github.com/adrianliechti/devkube/pkg/kubectl"

	"github.com/adrianliechti/devkube/extension/falco"
	"github.com/adrianliechti/devkube/extension/kyverno"
	"github.com/adrianliechti/devkube/extension/linkerd"
	"github.com/adrianliechti/devkube/extension/trivy"
)

func DisableCommand() *cli.Command {
	return &cli.Command{
		Name:  "disable",
		Usage: "Disable cluster feature",

		Category: app.FeaturesCategory,

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

			case "falco":
				if err := falco.Uninstall(c.Context, kubeconfig, app.DefaultNamespace); err != nil {
					return err
				}

				return nil

			case "kyverno":
				if err := kyverno.Uninstall(c.Context, kubeconfig, app.DefaultNamespace); err != nil {
					return err
				}

				return nil

			case "linkerd":
				if err := linkerd.Uninstall(c.Context, kubeconfig); err != nil {
					return err
				}

				return nil

			case "trivy":
				if err := trivy.Uninstall(c.Context, kubeconfig, app.DefaultNamespace); err != nil {
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
