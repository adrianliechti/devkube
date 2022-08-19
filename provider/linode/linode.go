package linode

import (
	"context"
	"encoding/base64"
	"errors"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/linode/linodego"
	"golang.org/x/oauth2"

	"github.com/adrianliechti/devkube/provider"
)

type Provider struct {
	client *linodego.Client
}

func New(token string) provider.Provider {
	client := &http.Client{
		Transport: &oauth2.Transport{
			Source: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}),
		},
	}

	c := linodego.NewClient(client)

	return &Provider{
		client: &c,
	}
}

func (p *Provider) List(ctx context.Context) ([]string, error) {
	var list []string

	cluster, err := p.client.ListLKEClusters(ctx, nil)

	if err != nil {
		return list, err
	}

	for _, c := range cluster {
		list = append(list, c.Label)
	}

	return list, nil
}

func (p *Provider) Create(ctx context.Context, name string, kubeconfig string) error {
	versions, err := p.client.ListLKEVersions(ctx, nil)

	if err != nil {
		return err
	}

	if len(versions) == 0 {
		return errors.New("no versions found")
	}

	version := versions[0].ID

	opts := linodego.LKEClusterCreateOptions{
		Label: name,

		Region:     "eu-central",
		K8sVersion: version,

		ControlPlane: &linodego.LKEClusterControlPlane{
			HighAvailability: false,
		},

		NodePools: []linodego.LKENodePoolCreateOptions{
			{
				Count: 3,
				Type:  "g6-standard-4", // Linode 8GB
			},
		},
	}

	cluster, err := p.client.CreateLKECluster(ctx, opts)

	// _ = opts
	// clusterID, err := p.clusterID(ctx, name)

	// if err != nil {
	// 	return err
	// }

	// cluster, err := p.client.GetLKECluster(ctx, clusterID)

	// if err != nil {
	// 	return err
	// }

	config, err := p.client.GetLKEClusterKubeconfig(ctx, cluster.ID)

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

	return p.client.DeleteLKECluster(ctx, id)
}

func (p *Provider) ExportConfig(ctx context.Context, name, kubeconfig string) error {
	id, err := p.clusterID(ctx, name)

	if err != nil {
		return err
	}

	config, err := p.client.GetLKEClusterKubeconfig(ctx, id)

	if err != nil {
		return err
	}

	return writeKubeconfig(kubeconfig, config)
}

func (p *Provider) clusterID(ctx context.Context, name string) (int, error) {
	cluster, err := p.client.ListLKEClusters(ctx, nil)

	if err != nil {
		return 0, err
	}

	for _, c := range cluster {
		if !strings.EqualFold(c.Label, name) {
			continue
		}

		return c.ID, nil
	}

	return 0, errors.New("cluster not found")
}

func writeKubeconfig(kubeconfig string, config *linodego.LKEClusterKubeconfig) error {
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
