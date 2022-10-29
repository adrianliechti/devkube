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

func EnableCommand() *cli.Command {
	return &cli.Command{
		Name:  "enable",
		Usage: "Enable cluster feature",

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
				if err := falco.Install(c.Context, kubeconfig, app.DefaultNamespace); err != nil {
					return err
				}

				return nil

			case "kyverno":
				if err := kyverno.Install(c.Context, kubeconfig, app.DefaultNamespace); err != nil {
					return err
				}

				return nil

			case "linkerd":
				if err := linkerd.Install(c.Context, kubeconfig); err != nil {
					return err
				}

				return nil

			case "trivy":
				if err := trivy.Install(c.Context, kubeconfig, app.DefaultNamespace); err != nil {
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
