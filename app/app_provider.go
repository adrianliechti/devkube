package app

import (
	"fmt"
	"strings"

	"github.com/adrianliechti/devkube/pkg/cli"
	"github.com/adrianliechti/devkube/provider"
	"github.com/adrianliechti/devkube/provider/kind"
)

func MustProvider(c *cli.Context) provider.Provider {
	provider, err := Provider(c)

	if err != nil {
		cli.Fatal(err)
	}

	return provider
}

func Provider(c *cli.Context) (provider.Provider, error) {
	provider := ""

	switch strings.ToLower(provider) {
	case "", "local":
		return kind.New(), nil

	default:
		return nil, fmt.Errorf("unknown provider %q", provider)
	}
}
