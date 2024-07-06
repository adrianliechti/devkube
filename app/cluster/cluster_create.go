package cluster

import (
	"github.com/adrianliechti/devkube/app"
	"github.com/adrianliechti/devkube/extension/certmanager"
	"github.com/adrianliechti/devkube/pkg/cli"
)

func CreateCommand() *cli.Command {
	return &cli.Command{
		Name:  "create",
		Usage: "Create cluster",

		// Before: func(c *cli.Context) error {
		// 	if _, _, err := helm.Info(c.Context); err != nil {
		// 		return err
		// 	}

		// 	if _, _, err := kubectl.Info(c.Context); err != nil {
		// 		return err
		// 	}

		// 	return nil
		// },

		Action: func(c *cli.Context) error {
			provider := app.MustProvider(c)
			cluster := "devkube"

			if err := provider.Create(c.Context, cluster); err != nil {
				return err
			}

			client := app.MustClient(c)

			// kubectl.Invoke(c.Context, []string{"create", "namespace", app.DefaultNamespace}, kubectl.WithKubeconfig(kubeconfig))

			// if err := observability.InstallCRD(c.Context, kubeconfig, app.DefaultNamespace); err != nil {
			// 	return err
			// }

			if err := certmanager.Ensure(c.Context, client); err != nil {
				return err
			}

			// if err := metrics.Install(c.Context, kubeconfig, app.DefaultNamespace); err != nil {
			// 	return err
			// }

			// if err := dashboard.Install(c.Context, kubeconfig, app.DefaultNamespace); err != nil {
			// 	return err
			// }

			// if err := registry.Install(c.Context, kubeconfig, app.DefaultNamespace); err != nil {
			// 	return err
			// }

			// if err := ingress.Install(c.Context, kubeconfig, app.DefaultNamespace); err != nil {
			// 	return err
			// }

			// if err := observability.Install(c.Context, kubeconfig, app.DefaultNamespace); err != nil {
			// 	return err
			// }

			return ExportConfig(c.Context, provider, cluster, "")
		},
	}
}
