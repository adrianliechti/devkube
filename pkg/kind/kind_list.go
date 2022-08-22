package kind

import (
	"context"

	"sigs.k8s.io/kind/pkg/cluster"
	"sigs.k8s.io/kind/pkg/log"
)

func List(ctx context.Context) ([]string, error) {
	logger := log.NoopLogger{}

	provider := cluster.NewProvider(
		cluster.ProviderWithLogger(logger),
		//runtime.GetDefault(logger),
	)

	return provider.List()
}
