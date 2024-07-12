package extension

import (
	"context"

	"github.com/adrianliechti/loop/pkg/kubernetes"

	"github.com/adrianliechti/devkube/extension/argocd"
	"github.com/adrianliechti/devkube/extension/certmanager"
	"github.com/adrianliechti/devkube/extension/crossplane"
	"github.com/adrianliechti/devkube/extension/dashboard"
	"github.com/adrianliechti/devkube/extension/gatekeeper"
	"github.com/adrianliechti/devkube/extension/grafana"
	"github.com/adrianliechti/devkube/extension/loki"
	"github.com/adrianliechti/devkube/extension/metrics"
	"github.com/adrianliechti/devkube/extension/monitoring"
	"github.com/adrianliechti/devkube/extension/otel"
	"github.com/adrianliechti/devkube/extension/promtail"
	"github.com/adrianliechti/devkube/extension/registry"
	"github.com/adrianliechti/devkube/extension/tekton"
	"github.com/adrianliechti/devkube/extension/tempo"
)

type Extension struct {
	Name string

	Title  string
	Ensure EnsureFunc
}

type EnsureFunc = func(ctx context.Context, client kubernetes.Client) error

var Default []Extension = []Extension{
	{
		Name:   "certmanager",
		Title:  "Cert-Manager",
		Ensure: certmanager.Ensure,
	},
	{
		Name:   "gatekeeper",
		Title:  "Gatekeeper",
		Ensure: gatekeeper.Ensure,
	},
	{
		Name:   "crossplane",
		Title:  "Crossplane",
		Ensure: crossplane.Ensure,
	},
	{
		Name:   "metrics",
		Title:  "Metrics",
		Ensure: metrics.Ensure,
	},
	{
		Name:   "registry",
		Title:  "Registry",
		Ensure: registry.Ensure,
	},
	{
		Name:   "dashboard",
		Title:  "Dashboard",
		Ensure: dashboard.Ensure,
	},
	{
		Name:   "prometheus",
		Title:  "Prometheus",
		Ensure: monitoring.Ensure,
	},
	{
		Name:   "loki",
		Title:  "Grafana Loki",
		Ensure: loki.Ensure,
	},
	{
		Name:   "tempo",
		Title:  "Grafana Tempo",
		Ensure: tempo.Ensure,
	},
	{
		Name:   "grafana",
		Title:  "Grafana",
		Ensure: grafana.Ensure,
	},
	{
		Name:   "promtail",
		Title:  "Promtail",
		Ensure: promtail.Ensure,
	},
	{
		Name:   "otel",
		Title:  "OpenTelemetry",
		Ensure: otel.Ensure,
	},
}

var Optional []Extension = []Extension{
	{
		Name:   "argocd",
		Title:  "Argo CD",
		Ensure: argocd.Ensure,
	},
	{
		Name:   "tekton",
		Title:  "Tekton",
		Ensure: tekton.Ensure,
	},
}
