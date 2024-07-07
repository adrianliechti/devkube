package create

import (
	"errors"

	"github.com/adrianliechti/devkube/app"
	"github.com/adrianliechti/devkube/app/setup"
	"github.com/adrianliechti/devkube/pkg/cli"

	"github.com/adrianliechti/devkube/extension/certmanager"
	"github.com/adrianliechti/devkube/extension/crossplane"
	"github.com/adrianliechti/devkube/extension/gatekeeper"
	"github.com/adrianliechti/devkube/extension/grafana"
	"github.com/adrianliechti/devkube/extension/loki"
	"github.com/adrianliechti/devkube/extension/metrics"
	"github.com/adrianliechti/devkube/extension/monitoring"
	"github.com/adrianliechti/devkube/extension/otel"
	"github.com/adrianliechti/devkube/extension/promtail"
	"github.com/adrianliechti/devkube/extension/tempo"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:  "create",
		Usage: "create cluster",

		Action: func(c *cli.Context) error {
			provider, cluster, _ := app.Cluster(c)

			if provider == nil {
				return errors.New("invalid provider specified")
			}

			if cluster == "" {
				return errors.New("invalid cluster specified")
			}

			cli.Info("★ installing Kubernetes Cluster...")

			if err := provider.Create(c.Context, cluster); err != nil {
				//return err
			}

			client := app.MustClient(c)

			cli.Info("★ installing Cert-Manager...")

			if err := certmanager.Ensure(c.Context, client); err != nil {
				return err
			}

			cli.Info("★ installing Gatekeeper...")

			if err := gatekeeper.Ensure(c.Context, client); err != nil {
				return err
			}

			cli.Info("★ installing Crossplane...")

			if err := crossplane.Ensure(c.Context, client); err != nil {
				return err
			}

			cli.Info("★ installing Prometheus...")

			if err := metrics.Ensure(c.Context, client); err != nil {
				return err
			}

			if err := monitoring.Ensure(c.Context, client); err != nil {
				return err
			}

			cli.Info("★ installing Grafana Loki...")

			if err := loki.Ensure(c.Context, client); err != nil {
				return err
			}

			cli.Info("★ installing Grafana Tempo...")

			if err := tempo.Ensure(c.Context, client); err != nil {
				return err
			}

			cli.Info("★ installing Grafana...")

			if err := grafana.Ensure(c.Context, client); err != nil {
				return err
			}

			cli.Info("★ installing Promtail...")

			if err := promtail.Ensure(c.Context, client); err != nil {
				return err
			}

			cli.Info("★ installing OpenTelemetry...")

			if err := otel.Ensure(c.Context, client); err != nil {
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
