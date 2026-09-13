package helm

import (
	"log/slog"
	"time"

	"github.com/adrianliechti/loop/pkg/kubernetes"

	"helm.sh/helm/v4/pkg/action"
	"helm.sh/helm/v4/pkg/chart"
	"helm.sh/helm/v4/pkg/chart/loader"
	helmcli "helm.sh/helm/v4/pkg/cli"
	"helm.sh/helm/v4/pkg/registry"
	"helm.sh/helm/v4/pkg/storage/driver"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

const timeout = 15 * time.Minute

var ErrNoDeployedReleases = driver.ErrNoDeployedReleases

func newConfiguration(client kubernetes.Client, namespace string) (*action.Configuration, error) {
	config := new(action.Configuration)
	config.LogHolder.SetLogger(slog.DiscardHandler)

	if err := config.Init(NewClientGetter(client, namespace), namespace, ""); err != nil {
		return nil, err
	}

	return config, nil
}

func loadChart(opts *action.ChartPathOptions, chartName string) (chart.Charter, error) {
	path, err := opts.LocateChart(chartName, helmcli.New())

	if err != nil {
		return nil, err
	}

	return loader.Load(path)
}

// setRegistryClient wires up OCI support - charts served over plain HTTP(S)
// work without it, so a failure here is not fatal.
func setRegistryClient(set func(*registry.Client)) {
	if client, err := registry.NewClient(); err == nil {
		set(client)
	}
}

func NewClientGetter(client kubernetes.Client, namespace string) genericclioptions.RESTClientGetter {
	return &clientGetter{
		client:    client,
		namespace: namespace,
	}
}

type clientGetter struct {
	client    kubernetes.Client
	namespace string

	mapper    meta.RESTMapper
	discovery discovery.CachedDiscoveryInterface
}

func (c *clientGetter) ToRESTConfig() (*rest.Config, error) {
	return c.client.Config(), nil
}

func (c *clientGetter) ToDiscoveryClient() (discovery.CachedDiscoveryInterface, error) {
	if c.discovery == nil {
		client, err := discovery.NewDiscoveryClientForConfig(c.client.Config())

		if err != nil {
			return nil, err
		}

		c.discovery = memory.NewMemCacheClient(client)
	}

	return c.discovery, nil
}

func (c *clientGetter) ToRESTMapper() (meta.RESTMapper, error) {
	if c.mapper == nil {
		dc, err := c.ToDiscoveryClient()

		if err != nil {
			return nil, err
		}

		c.mapper = restmapper.NewDeferredDiscoveryRESTMapper(dc)
	}

	return c.mapper, nil
}

func (c *clientGetter) ToRawKubeConfigLoader() clientcmd.ClientConfig {
	return &clientConfig{client: c.client, namespace: c.namespace}
}

type clientConfig struct {
	client    kubernetes.Client
	namespace string
}

func (c *clientConfig) ClientConfig() (*rest.Config, error) {
	return c.client.Config(), nil
}

func (c *clientConfig) Namespace() (string, bool, error) {
	if c.namespace != "" {
		return c.namespace, true, nil
	}

	if val := c.client.Namespace(); val != "" {
		return val, true, nil
	}

	return "default", true, nil
}

func (c *clientConfig) RawConfig() (clientcmdapi.Config, error) {
	panic("not implemented")
}

func (c *clientConfig) ConfigAccess() clientcmd.ConfigAccess {
	panic("not implemented")
}
