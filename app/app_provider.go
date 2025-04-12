package app

import (
	"context"
	"fmt"
	"strings"

	"github.com/adrianliechti/devkube/provider"
	"github.com/adrianliechti/devkube/provider/kind"
	"github.com/adrianliechti/devkube/provider/none"
	"github.com/adrianliechti/go-cli"
)

var ProviderFlag = &cli.StringFlag{
	Name:  "provider",
	Usage: "cluster provider",
}

func MustProvider(ctx context.Context, cmd *cli.Command) provider.Provider {
	provider, err := Provider(ctx, cmd)

	if err != nil {
		cli.Fatal(err)
	}

	return provider
}

func Provider(ctx context.Context, cmd *cli.Command) (provider.Provider, error) {
	provider := cmd.String(ProviderFlag.Name)

	switch strings.ToLower(provider) {
	case "", "local":
		return kind.New()

	case "none":
		return none.NewFromEnvironment()

	default:
		return nil, fmt.Errorf("unknown provider %q", provider)
	}
}
