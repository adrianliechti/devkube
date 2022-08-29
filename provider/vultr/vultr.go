package vultr

import (
	"context"
	"encoding/base64"
	"errors"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/vultr/govultr/v2"
	"golang.org/x/oauth2"

	"github.com/adrianliechti/devkube/provider"
)

type Provider struct {
	client *govultr.Client
}

func New(token string) provider.Provider {
	client := &http.Client{
		Transport: &oauth2.Transport{
			Source: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}),
		},
	}

	c := govultr.NewClient(client)

	return &Provider{
		client: c,
	}
}

func (p *Provider) List(ctx context.Context) ([]string, error) {
	var list []string

	cluster, _, err := p.client.Kubernetes.ListClusters(ctx, nil)

	if err != nil {
		return list, err
	}

	for _, c := range cluster {
		list = append(list, c.Label)
	}

	return list, nil
}

func (p *Provider) Create(ctx context.Context, name string, kubeconfig string) error {
	versions, err := p.client.Kubernetes.GetVersions(ctx)

	if err != nil {
		return err
	}

	if len(versions.Versions) == 0 {
		return errors.New("no versions found")
	}

	version := versions.Versions[0]

	cluster, err := p.client.Kubernetes.CreateCluster(ctx, &govultr.ClusterReq{
		Label: name,

		Region:  "ams",
		Version: version,

		NodePools: []govultr.NodePoolReq{
			{
				NodeQuantity: 3,

				Label: name + "-pool",
				Plan:  "voc-c-4c-8gb-75s-amd",
			},
		},
	})

	if err != nil {
		return err
	}

	for {
		cluster, err := p.client.Kubernetes.GetCluster(ctx, cluster.ID)

		if err != nil {
			return err
		}

		if cluster.Status == "active" {
			break
		}

		println(cluster.Status)

		time.Sleep(20 * time.Second)
	}

	config, err := p.client.Kubernetes.GetKubeConfig(ctx, cluster.ID)

	if err != nil {
		return err
	}

	return writeKubeconfig(kubeconfig, config)
}

func (p *Provider) Delete(ctx context.Context, name string) error {
	id, err := p.clusterID(ctx, name)

	if err != nil {
		return err
	}

	return p.client.Kubernetes.DeleteCluster(ctx, id)
}

func (p *Provider) Export(ctx context.Context, name, kubeconfig string) error {
	id, err := p.clusterID(ctx, name)

	if err != nil {
		return err
	}

	config, err := p.client.Kubernetes.GetKubeConfig(ctx, id)

	if err != nil {
		return err
	}

	return writeKubeconfig(kubeconfig, config)
}

func (p *Provider) clusterID(ctx context.Context, name string) (string, error) {
	cluster, _, err := p.client.Kubernetes.ListClusters(ctx, nil)

	if err != nil {
		return "", err
	}

	for _, c := range cluster {
		if !strings.EqualFold(c.Label, name) {
			continue
		}

		return c.ID, nil
	}

	return "", errors.New("cluster not found")
}

func writeKubeconfig(kubeconfig string, config *govultr.KubeConfig) error {
	if kubeconfig == "" {
		home, err := os.UserHomeDir()

		if err != nil {
			return err
		}

		dir := path.Join(home, ".kube")

		if err := os.MkdirAll(dir, 0700); err != nil {
			return err
		}

		kubeconfig = path.Join(home, ".kube", "config")
	}

	data, err := base64.StdEncoding.DecodeString(config.KubeConfig)

	if err != nil {
		return err
	}

	return os.WriteFile(kubeconfig, data, 0600)
}
