package load

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/adrianliechti/devkube/app"
	"github.com/adrianliechti/go-cli"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:  "load",
		Usage: "load image into registry",

		Flags: []cli.Flag{
			app.PortFlag,
		},

		Action: func(ctx context.Context, cmd *cli.Command) error {
			if cmd.Args().Len() != 1 {
				return errors.New("needs one argument: image")
			}

			image := cmd.Args().Get(0)

			client := app.MustClient(ctx, cmd)
			port := app.MustPortOrRandom(ctx, cmd, 5555)

			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			ready := make(chan struct{})
			done := make(chan error, 1)

			go func() {
				done <- client.ServicePortForward(ctx, "platform", "registry", "", map[int]int{port: 80}, ready)
			}()

			select {
			case <-ready:
			case err := <-done:
				if err == nil {
					err = errors.New("port forward closed unexpectedly")
				}

				return err
			}

			return LoadImage(ctx, image, fmt.Sprintf("localhost:%d", port))
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

	archive := filepath.Join(dir, "image.tar")

	if out, err := exec.CommandContext(ctx, "docker", "image", "save", source, "-o", archive).CombinedOutput(); err != nil {
		if message := strings.TrimSpace(string(out)); message != "" {
			return fmt.Errorf("failed to save image: %s", message)
		}

		return fmt.Errorf("failed to save image: %w", err)
	}

	image, err := tarball.ImageFromPath(archive, &src)

	if err != nil {
		return err
	}

	return remote.Write(dst, image, remote.WithContext(ctx))
}
