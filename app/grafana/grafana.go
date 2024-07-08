package grafana

import (
	"fmt"

	"github.com/adrianliechti/devkube/app"
	"github.com/adrianliechti/devkube/pkg/cli"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:  "grafana",
		Usage: "open Grafana in Browser",

		Action: func(c *cli.Context) error {
			client := app.MustClient(c)

			port := app.MustPortOrRandom(c, 3000)

			ready := make(chan struct{})

			go func() {
				<-ready

				url := fmt.Sprintf("http://127.0.0.1:%d", port)
				cli.OpenURL(url)
			}()

			return client.ServicePortForward(c.Context, "platform", "grafana", "", map[int]int{port: 3000}, ready)
		},
	}
}
