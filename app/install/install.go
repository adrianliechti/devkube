package install

import (
	"errors"
	"strings"

	"github.com/adrianliechti/devkube/app"
	"github.com/adrianliechti/devkube/extension"
	"github.com/adrianliechti/devkube/pkg/cli"

	"github.com/adrianliechti/devkube/extension/argocd"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:  "install",
		Usage: "install extension",

		Hidden: true,

		Action: func(c *cli.Context) error {
			provider, cluster, _ := app.Cluster(c)

			if provider == nil {
				return errors.New("invalid provider specified")
			}

			if cluster == "" {
				return errors.New("invalid cluster specified")
			}

			client := app.MustClient(c)

			fns := []extension.EnsureFunc{}

			for _, extension := range c.Args().Slice() {
				switch strings.ToLower(extension) {
				case "argocd":
					fns = append(fns, argocd.Ensure)
				default:
					return errors.New("unknown extension: " + extension)
				}
			}

			cli.Info("★ installing extensions(s)...")

			var result error

			for _, fn := range fns {
				if err := fn(c.Context, client); err != nil {
					result = errors.Join(result, err)
				}
			}

			return result
		},
	}
}
