package kind

import (
	"context"

	"github.com/adrianliechti/devkube/pkg/kind"
	"github.com/adrianliechti/devkube/provider"
)

type Provider struct {
}

func New() provider.Provider {
	return new(Provider)
}

func (p *Provider) List(ctx context.Context) ([]string, error) {
	return kind.List(ctx)
}

func (p *Provider) Create(ctx context.Context, name string, kubeconfig string) error {
	var config = map[string]any{
		"kind":       "Cluster",
		"apiVersion": "kind.x-k8s.io/v1alpha4",

		"kubeadmConfigPatches": []string{
			`kind: ClusterConfiguration
controllerManager:
  extraArgs:
    bind-address: 0.0.0.0
scheduler:
  extraArgs:
    bind-address: 0.0.0.0
etcd:
  local:
    extraArgs:
      listen-metrics-urls: http://0.0.0.0:2381
`,
			`kind: KubeProxyConfiguration
metricsBindAddress: 0.0.0.0
`,
		},
	}

	if err := kind.Create(ctx, name, config, kubeconfig, kind.WithDefaultOutput()); err != nil {
		return err
	}

	return nil
}

func (p *Provider) Delete(ctx context.Context, name string) error {
	return kind.Delete(ctx, name)
}

func (p *Provider) ExportConfig(ctx context.Context, name, path string) error {
	return kind.ExportConfig(ctx, name, path)
}
