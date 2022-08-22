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

	data, err := yaml.Marshal(config)

	if err != nil {
		return err
	}

	return provider.Create(
		name,
		cluster.CreateWithRawConfig(data),
		cluster.CreateWithWaitForReady(time.Duration(0)),
		cluster.CreateWithKubeconfigPath(kubeconfig),
		//cluster.CreateWithDisplayUsage(true),
		//cluster.CreateWithDisplaySalutation(true),
	)
}
