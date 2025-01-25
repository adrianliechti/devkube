package build

import (
	"context"
	"errors"

	"github.com/adrianliechti/devkube/app"
	"github.com/adrianliechti/devkube/pkg/cli"
	"github.com/adrianliechti/loop/pkg/remote/build"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:  "build",
		Usage: "build image into registry",

		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := app.MustClient(ctx, cmd)

			if cmd.Args().Len() != 2 {
				return errors.New("needs two arguments: image and context path")
			}

			image, err := build.ParseImage("registry.platform/" + cmd.Args().Get(0))

			if err != nil {
				return err
			}

			path, err := build.ParsePath(cmd.Args().Get(1))

			if err != nil {
				return err
			}

			image.Insecure = true

			return build.Run(ctx, client, image, path, "", nil)
		},
	}
}
