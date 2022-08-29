package kind

import (
	"context"

	"sigs.k8s.io/kind/pkg/cluster"
	"sigs.k8s.io/kind/pkg/log"
)

func Export(ctx context.Context, name, kubeconfig string, opt ...Option) error {
	logger := log.NoopLogger{}

	provider := cluster.NewProvider(
		cluster.ProviderWithLogger(logger),
		//runtime.GetDefault(logger),
	)

	return provider.ExportKubeConfig(name, kubeconfig, false)
}
