package create

import (
	"github.com/adrianliechti/devkube/app"
	"github.com/adrianliechti/devkube/app/setup"
	"github.com/adrianliechti/devkube/pkg/cli"

	"github.com/adrianliechti/devkube/extension/certmanager"
	"github.com/adrianliechti/devkube/extension/grafana"
	"github.com/adrianliechti/devkube/extension/loki"
	"github.com/adrianliechti/devkube/extension/metrics"
	"github.com/adrianliechti/devkube/extension/monitoring"
	"github.com/adrianliechti/devkube/extension/promtail"
	"github.com/adrianliechti/devkube/extension/tempo"
)

func Command() *cli.Command {
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

			if err := certmanager.Ensure(c.Context, client); err != nil {
				return err
			}

			if err := metrics.Ensure(c.Context, client); err != nil {
				return err
			}

			if err := monitoring.Ensure(c.Context, client); err != nil {
				return err
			}

			if err := loki.Ensure(c.Context, client); err != nil {
				return err
			}

			if err := promtail.Ensure(c.Context, client); err != nil {
				return err
			}

			if err := tempo.Ensure(c.Context, client); err != nil {
				return err
			}

			if err := grafana.Ensure(c.Context, client); err != nil {
				return err
			}

			// if err := dashboard.Install(c.Context, kubeconfig, app.DefaultNamespace); err != nil {
			// 	return err
			// }

			// if err := registry.Install(c.Context, kubeconfig, app.DefaultNamespace); err != nil {
			// 	return err
			// }

			// if err := ingress.Install(c.Context, kubeconfig, app.DefaultNamespace); err != nil {
			// 	return err
			// }

			return setup.Export(c.Context, provider, cluster, "")
		},
	}
}
