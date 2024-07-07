package app

import (
	"fmt"
	"strings"

	"github.com/adrianliechti/devkube/pkg/cli"
	"github.com/adrianliechti/devkube/provider"
	"github.com/adrianliechti/devkube/provider/kind"
	"github.com/adrianliechti/devkube/provider/none"
)

var ProviderFlag = &cli.StringFlag{
	Name:  "provider",
	Usage: "provider name",
}

func MustProvider(c *cli.Context) provider.Provider {
	provider, err := Provider(c)

	if err != nil {
		cli.Fatal(err)
	}

	return provider
}

func Provider(c *cli.Context) (provider.Provider, error) {
	provider := c.String(ProviderFlag.Name)

	switch strings.ToLower(provider) {
	case "", "local":
		return kind.New()

	case "none":
		return none.NewFromEnvironment()

	default:
		return nil, fmt.Errorf("unknown provider %q", provider)
	}
}
