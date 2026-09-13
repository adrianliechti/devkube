package kind

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/adrianliechti/devkube/provider"

	"sigs.k8s.io/kind/pkg/cluster"
	"sigs.k8s.io/kind/pkg/log"
)

type kind struct {
	provider *cluster.Provider
}

var (
	//go:embed config.yaml
	config []byte
)

func New() (provider.Provider, error) {
	logger := log.NoopLogger{}

	opts := []cluster.ProviderOption{
		cluster.ProviderWithLogger(logger),
	}

	if o, err := cluster.DetectNodeProvider(); err == nil {
		opts = append(opts, o)
	}

	return &kind{
		provider: cluster.NewProvider(opts...),
	}, nil
}

func (k *kind) List(ctx context.Context) ([]string, error) {
	return k.provider.List()
}

func (k *kind) Create(ctx context.Context, name string) error {
	dir, err := os.MkdirTemp("", "kubeconfig-")

	if err != nil {
		return err
	}

	defer os.RemoveAll(dir)

	// kind writes the kubeconfig itself - keep it out of the user's config,
	// merging it in is handled by the setup command
	kubeconfig := filepath.Join(dir, "config")

	opts := []cluster.CreateOption{
		cluster.CreateWithRawConfig(config),
		cluster.CreateWithKubeconfigPath(kubeconfig),
		cluster.CreateWithWaitForReady(0),
	}

	return k.provider.Create(name, opts...)
}

func (k *kind) Delete(ctx context.Context, name string) error {
	return k.provider.Delete(name, "")
}

func (k *kind) Start(ctx context.Context, name string) error {
	return control(ctx, "start", name)
}

func (k *kind) Stop(ctx context.Context, name string) error {
	return control(ctx, "stop", name)
}

// kind offers no API to start/stop an existing cluster, so drive the
// container runtime directly - using the first one that accepts the node.
func control(ctx context.Context, action, name string) error {
	container := strings.ToLower(name + "-control-plane")

	for _, runtime := range []string{"docker", "podman", "nerdctl"} {
		if _, err := exec.LookPath(runtime); err != nil {
			continue
		}

		if err := exec.CommandContext(ctx, runtime, action, container).Run(); err == nil {
			return nil
		}
	}

	return fmt.Errorf("unable to %s cluster %s", action, name)
}

func (k *kind) Config(ctx context.Context, name string) ([]byte, error) {
	data, err := k.provider.KubeConfig(name, false)

	if err != nil {
		return nil, err
	}

	return []byte(data), nil
}
