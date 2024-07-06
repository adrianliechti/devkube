package kind

import (
	"context"

	"sigs.k8s.io/kind/pkg/cluster"
	"sigs.k8s.io/kind/pkg/log"
)

func Config(ctx context.Context, name string, opt ...Option) ([]byte, error) {
	logger := log.NoopLogger{}

	provider := cluster.NewProvider(
		cluster.ProviderWithLogger(logger),
		//runtime.GetDefault(logger),
	)

	data, err := provider.KubeConfig(name, false)

	if err != nil {
		return nil, err
	}

	return []byte(data), nil
}
