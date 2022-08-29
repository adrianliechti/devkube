package digitalocean

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/digitalocean/godo"

	"github.com/adrianliechti/devkube/provider"
)

type Provider struct {
	client *godo.Client
}

func New(token string) provider.Provider {
	client := godo.NewFromToken(token)

	return &Provider{
		client: client,
	}
}

func NewFromEnvironment() (provider.Provider, error) {
	token := os.Getenv("DIGITALOCEAN_TOKEN")

	if token == "" {
		return nil, fmt.Errorf("DIGITALOCEAN_TOKEN is not set")
	}

	return New(token), nil
}

func (p *Provider) List(ctx context.Context) ([]string, error) {
	var list []string

	opt := &godo.ListOptions{}

	for {
		clusters, resp, err := p.client.Kubernetes.List(ctx, opt)

		if err != nil {
			return nil, err
		}

		for _, cluster := range clusters {
			list = append(list, cluster.Name)
		}

		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()

		if err != nil {
			return nil, err
		}

		opt.Page = page + 1
	}

	return list, nil
}

func (p *Provider) Create(ctx context.Context, name string, kubeconfig string) error {
	options, _, err := p.client.Kubernetes.GetOptions(ctx)

	if err != nil {
		return err
	}

	if len(options.Versions) == 0 {
		return errors.New("no versions found")
	}

	version := options.Versions[0].Slug

	cluster, _, err := p.client.Kubernetes.Create(ctx, &godo.KubernetesClusterCreateRequest{
		Name: name,

		RegionSlug:  "fra1",
		VersionSlug: version,

		NodePools: []*godo.KubernetesNodePoolCreateRequest{
			{
				Name: "default",

				Size:  "s-4vcpu-8gb",
				Count: 2,
			},
		},

		AutoUpgrade:  true,
		SurgeUpgrade: true,
	})

	if err != nil {
		return err
	}

	for {
		println("waiting for cluster to be ready")

		time.Sleep(10 * time.Second)

		cluster, _, err := p.client.Kubernetes.Get(ctx, cluster.ID)

		if err != nil {
			return err
		}

		if cluster.Status.State == godo.KubernetesClusterStatusError {
			return errors.New("cluster creation failed")
		}

		if cluster.Status.State == godo.KubernetesClusterStatusInvalid {
			return errors.New("cluster creation failed")
		}

		if cluster.Status.State == godo.KubernetesClusterStatusRunning {
			break
		}
	}

	return p.Export(ctx, name, kubeconfig)
}

func (p *Provider) Delete(ctx context.Context, name string) error {
	id, err := p.clusterID(ctx, name)

	if err != nil {
		return err
	}

	_, err = p.client.Kubernetes.DeleteDangerous(ctx, id)
	return err
}

func (p *Provider) Export(ctx context.Context, name, kubeconfig string) error {
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

	id, err := p.clusterID(ctx, name)

	if err != nil {
		return err
	}

	config, _, err := p.client.Kubernetes.GetKubeConfig(ctx, id)

	if err != nil {
		return err
	}

	return os.WriteFile(kubeconfig, config.KubeconfigYAML, 0600)
}

func (p *Provider) clusterID(ctx context.Context, name string) (string, error) {
	opt := &godo.ListOptions{}

	for {
		clusters, resp, err := p.client.Kubernetes.List(ctx, opt)

		if err != nil {
			return "", err
		}

		for _, c := range clusters {
			if !strings.EqualFold(c.Name, name) {
				continue
			}

			return c.ID, nil
		}

		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()

		if err != nil {
			return "", err
		}

		opt.Page = page + 1
	}

	return "", errors.New("cluster not found")
}
