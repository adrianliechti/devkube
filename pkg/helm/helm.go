package helm

import (
	"github.com/adrianliechti/loop/pkg/kubernetes"

	"helm.sh/helm/v3/pkg/storage/driver"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

var ErrReleaseExists = driver.ErrReleaseExists
var ErrReleaseNotFound = driver.ErrReleaseNotFound
var ErrNoDeployedReleases = driver.ErrNoDeployedReleases

func NewClientGetter(client kubernetes.Client) genericclioptions.RESTClientGetter {
	return &clientGetter{
		client: client,
	}
}

type clientGetter struct {
	client kubernetes.Client

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
	return &clientConfig{client: c.client}
}

type clientConfig struct {
	client kubernetes.Client
}

func (c *clientConfig) ClientConfig() (*rest.Config, error) {
	return c.client.Config(), nil
}

func (c *clientConfig) Namespace() (string, bool, error) {
	return c.client.Namespace(), true, nil
}

func (c *clientConfig) RawConfig() (clientcmdapi.Config, error) {
	panic("not implemented")
}

func (c *clientConfig) ConfigAccess() clientcmd.ConfigAccess {
	panic("not implemented")
}
