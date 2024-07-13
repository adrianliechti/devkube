package load

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/adrianliechti/devkube/app"
	"github.com/adrianliechti/devkube/pkg/cli"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:  "load",
		Usage: "load image into registry",

		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := app.MustClient(ctx, cmd)

			if cmd.Args().Len() != 1 {
				return errors.New("needs one arguments: image")
			}

			image := cmd.Args().Get(0)

			port := app.MustPortOrRandom(ctx, cmd, 5555)
			ready := make(chan struct{})

			go func() {
				if err := client.ServicePortForward(ctx, "platform", "registry", "", map[int]int{port: 80}, ready); err != nil {
					log.Fatal(err)
				}
			}()

			<-ready

			if err := LoadImage(ctx, image, fmt.Sprintf("localhost:%d", port)); err != nil {
				log.Fatal(err)
			}

			return nil
		},
	}
}

func LoadImage(ctx context.Context, source, registry string) error {
	src, err := name.NewTag(source)

	if err != nil {
		return err
	}

	dst, err := name.ParseReference(path.Join(registry, source))

	if err != nil {
		return err
	}

	dir, err := os.MkdirTemp("", "container")

	if err != nil {
		return err
	}

	defer os.RemoveAll(dir)

	path := path.Join(dir, "image.tar")

	if err := exec.CommandContext(ctx, "docker", "image", "save", source, "-o", path).Run(); err != nil {
		return errors.New("failed to save image")
	}

	image, err := tarball.ImageFromPath(path, &src)

	if err != nil {
		return err
	}

	if err := remote.Write(dst, image); err != nil {
		return err
	}

	return nil
}
