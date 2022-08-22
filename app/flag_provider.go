package app

import (
	"fmt"
	"os"

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

func Provider(c *cli.Context) (provider.Provider, error) {
	provider := c.String(ProviderFlag.Name)

	switch provider {
	case "none":
		return none.New(), nil

	case "kind", "":
		if _, _, err := dockercli.Info(c.Context); err != nil {
			return nil, err
		}

		return kind.New(), nil

	case "linode":
		token := os.Getenv("LINODE_TOKEN")

		if token == "" {
			return nil, fmt.Errorf("LINODE_TOKEN is not set")
		}

		return linode.New(token), nil

	case "vultr":
		token := os.Getenv("VULTR_API_KEY")

		if token == "" {
			return nil, fmt.Errorf("VULTR_API_KEY is not set")
		}

		return vultr.New(token), nil

	default:
		return nil, fmt.Errorf("unknown provider %q", provider)
	}
}

func MustProvider(c *cli.Context) provider.Provider {
	provider, err := Provider(c)

	if err != nil {
		cli.Fatal(err)
	}

	return provider
}
