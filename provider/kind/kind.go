package kind

import (
	"context"
	_ "embed"
	"os"
	"os/exec"
	"strings"
	"time"

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

	kubeconfig := dir + "/.config"

	defer os.RemoveAll(dir)

	opts := []cluster.CreateOption{
		cluster.CreateWithRawConfig(config),
		cluster.CreateWithKubeconfigPath(kubeconfig),
		cluster.CreateWithWaitForReady(time.Duration(0)),
	}

	return k.provider.Create(name, opts...)
}

func (k *kind) Delete(ctx context.Context, name string) error {
	return k.provider.Delete(name, "")
}

func (k *kind) Start(ctx context.Context, name string) error {
	container := strings.ToLower(name + "-control-plane")

	exec.Command("docker", "start", container).Run()
	exec.Command("podman", "start", container).Run()
	exec.Command("nerdctl", "start", container).Run()

	return nil
}

func (k *kind) Stop(ctx context.Context, name string) error {
	container := strings.ToLower(name + "-control-plane")

	exec.Command("docker", "stop", container).Run()
	exec.Command("podman", "stop", container).Run()
	exec.Command("nerdctl", "stop", container).Run()

	return nil
}

func (k *kind) Config(ctx context.Context, name string) ([]byte, error) {
	data, err := k.provider.KubeConfig(name, false)

	if err != nil {
		return nil, err
	}

	return []byte(data), nil
}
