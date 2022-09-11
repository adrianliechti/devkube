package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/adrianliechti/devkube/pkg/cli"
	"github.com/adrianliechti/devkube/pkg/docker"

	"github.com/adrianliechti/devkube/provider"
	"github.com/adrianliechti/devkube/provider/aws"
	"github.com/adrianliechti/devkube/provider/azure"
	"github.com/adrianliechti/devkube/provider/digitalocean"
	"github.com/adrianliechti/devkube/provider/kind"
	"github.com/adrianliechti/devkube/provider/linode"
	"github.com/adrianliechti/devkube/provider/none"
	"github.com/adrianliechti/devkube/provider/vultr"
)

var ProviderFlag = &cli.StringFlag{
	Name:  "provider",
	Usage: "Cluster provider",
}

func ListProviders() []string {
	return []string{
		"local",
		"aws",
		"azure",
		"digitalocean",
		"linode",
		"vultr",
	}
}

func SelectProvider(c *cli.Context) (provider.Provider, error) {
	if name := c.String(ProviderFlag.Name); name != "" {
		return ProviderFromName(c.Context, name)
	}

	items := ListProviders()

	if len(items) == 0 {
		return nil, errors.New("no provider found")
	}

	i, _, err := cli.Select("Select provider", items)

	if err != nil {
		return nil, err
	}

	name := items[i]
	return ProviderFromName(c.Context, name)
}

func MustProvider(c *cli.Context) provider.Provider {
	provider, err := SelectProvider(c)

	if err != nil {
		cli.Fatal(err)
	}

	return provider
}

func ProviderFromName(ctx context.Context, name string) (provider.Provider, error) {
	switch name {
	case "none":
		return none.NewFromEnvironment()

	case "local", "":
		if _, _, err := docker.Info(ctx); err != nil {
			return nil, err
		}

		return kind.New(), nil

	case "aws":
		return aws.NewFromEnvironment()

	case "azure":
		return azure.NewFromEnvironment()

	case "digitalocean":
		return digitalocean.NewFromEnvironment()

	case "linode":
		return linode.NewFromEnvironment()

	case "vultr":
		return vultr.NewFromEnvironment()

	default:
		return nil, fmt.Errorf("unknown provider %q", name)
	}
}
