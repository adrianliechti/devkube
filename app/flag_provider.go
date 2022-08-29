package app

import (
	"context"
	"fmt"

	"github.com/adrianliechti/devkube/pkg/cli"
	"github.com/adrianliechti/devkube/provider"
	"github.com/adrianliechti/devkube/provider/kind"
	"github.com/adrianliechti/devkube/provider/linode"
	"github.com/adrianliechti/devkube/provider/none"
	"github.com/adrianliechti/devkube/provider/vultr"

	dockercli "github.com/adrianliechti/devkube/pkg/docker"
)

var ProviderFlag = &cli.StringFlag{
	Name:  "provider",
	Usage: "Cluster provider",
}

func ListProviders() []string {
	return []string{
		"kind",
		"linode",
		"vultr",
	}
}

func ProviderFromName(ctx context.Context, name string) (provider.Provider, error) {
	switch name {
	case "none":
		return none.NewFromEnvironment()

	case "kind", "":
		if _, _, err := dockercli.Info(ctx); err != nil {
			return nil, err
		}

		return kind.New(), nil

	case "linode":
		return linode.NewFromEnvironment()

	case "vultr":
		return vultr.NewFromEnvironment()

	default:
		return nil, fmt.Errorf("unknown provider %q", name)
	}
}

func Provider(c *cli.Context) (provider.Provider, error) {
	name := c.String(ProviderFlag.Name)
	return ProviderFromName(c.Context, name)
}

func MustProvider(c *cli.Context) provider.Provider {
	provider, err := Provider(c)

	if err != nil {
		cli.Fatal(err)
	}

	return provider
}
