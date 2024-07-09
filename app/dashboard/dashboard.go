package dashboard

import (
	"fmt"

	"github.com/adrianliechti/devkube/app"
	"github.com/adrianliechti/devkube/pkg/cli"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:  "dashboard",
		Usage: "open Dashboard in Browser",

		Action: func(c *cli.Context) error {
			client := app.MustClient(c)

			port := app.MustPortOrRandom(c, 9090)

			ready := make(chan struct{})

			go func() {
				<-ready

				url := fmt.Sprintf("http://127.0.0.1:%d", port)
				cli.OpenURL(url)
			}()

			return client.ServicePortForward(c.Context, "platform", "dashboard", "", map[int]int{port: 8080}, ready)
		},
	}
}