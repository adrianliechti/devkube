package cluster

import (
	"github.com/adrianliechti/devkube/app"
	"github.com/adrianliechti/devkube/extension/prometheus"
	"github.com/adrianliechti/devkube/pkg/cli"
)

func CreateCommand() *cli.Command {
	return &cli.Command{
		Name:  "create",
		Usage: "Create cluster",

		Action: func(c *cli.Context) error {
			provider := app.MustProvider(c)
			cluster := "devkube"

			if err := provider.Create(c.Context, cluster); err != nil {
				//return err
			}

			client := app.MustClient(c)

			// if err := certmanager.Ensure(c.Context, client); err != nil {
			// 	//return err
			// }

			if err := prometheus.Ensure(c.Context, client); err != nil {
				//return err
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
