package app

import (
	"fmt"

	"github.com/adrianliechti/devkube/pkg/cli"
	"github.com/adrianliechti/devkube/provider"
	"github.com/adrianliechti/devkube/provider/kind"
	"github.com/adrianliechti/devkube/provider/none"

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
