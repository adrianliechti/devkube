package build

import (
	"errors"

	"github.com/adrianliechti/devkube/app"
	"github.com/adrianliechti/devkube/pkg/cli"
	"github.com/adrianliechti/loop/app/build"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:  "build",
		Usage: "build image into registry",

		Action: func(c *cli.Context) error {
			client := app.MustClient(c)

			if c.Args().Len() != 2 {
				return errors.New("needs two arguments: image and context path")
			}

			image, err := build.ParseImage("registry.platform/" + c.Args().Get(0))

			if err != nil {
				return err
			}

			path, err := build.ParsePath(c.Args().Get(1))

			if err != nil {
				return err
			}

			image.Insecure = true

			return build.Run(c.Context, client, "", image, path, "")
		},
	}
}
