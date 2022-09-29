package kind

import (
	"context"
	"time"

	"sigs.k8s.io/kind/pkg/cluster"
	"sigs.k8s.io/kind/pkg/log"

	"gopkg.in/yaml.v3"
)

func Create(ctx context.Context, name string, config map[string]any, kubeconfig string, opt ...Option) error {
	logger := log.NoopLogger{}

	provider := cluster.NewProvider(
		cluster.ProviderWithLogger(logger),
		//runtime.GetDefault(logger),
	)

	opts := []cluster.CreateOption{
		cluster.CreateWithWaitForReady(time.Duration(0)),
		cluster.CreateWithKubeconfigPath(kubeconfig),
		//cluster.CreateWithDisplayUsage(true),
		//cluster.CreateWithDisplaySalutation(true),
	}

	if config != nil {
		data, err := yaml.Marshal(config)

		if err != nil {
			return err
		}

		opts = append(opts, cluster.CreateWithRawConfig(data))

	}

	return provider.Create(name, opts...)
}
