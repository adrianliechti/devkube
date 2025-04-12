package grafana

import (
	"context"
	"fmt"

	"github.com/adrianliechti/devkube/app"
	"github.com/adrianliechti/go-cli"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:  "grafana",
		Usage: "open Grafana in Browser",

		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := app.MustClient(ctx, cmd)

			port := app.MustPortOrRandom(ctx, cmd, 3000)

			ready := make(chan struct{})

			go func() {
				<-ready

				url := fmt.Sprintf("http://127.0.0.1:%d", port)
				cli.OpenURL(url)
			}()

			return client.ServicePortForward(ctx, "platform", "grafana", "", map[int]int{port: 3000}, ready)
		},
	}
}
