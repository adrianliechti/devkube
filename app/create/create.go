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

			cli.Info("📦 Creating Cluster...")

			if err := provider.Create(c.Context, cluster); err != nil {
				return err
			}

			client := app.MustClient(c)

			cli.Info("📦 Installing Cert Manager...")

			if err := certmanager.Ensure(c.Context, client); err != nil {
				return err
			}

			cli.Info("📦 Installing Metrics Server...")

			if err := metrics.Ensure(c.Context, client); err != nil {
				return err
			}

			cli.Info("📦 Installing Prometheus Operator...")

			if err := monitoring.Ensure(c.Context, client); err != nil {
				return err
			}

			cli.Info("📦 Installing Loki...")

			if err := loki.Ensure(c.Context, client); err != nil {
				return err
			}

			cli.Info("📦 Installing Promtail...")

			if err := promtail.Ensure(c.Context, client); err != nil {
				return err
			}

			cli.Info("📦 Installing Tempo...")

			if err := tempo.Ensure(c.Context, client); err != nil {
				return err
			}

			cli.Info("📦 Installing Grafana...")

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
