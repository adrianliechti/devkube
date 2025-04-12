package otel

import (
	"context"

	"github.com/adrianliechti/devkube/app"
	"github.com/adrianliechti/go-cli"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:  "otel",
		Usage: "forward otel collector",

		Action: func(ctx context.Context, cmd *cli.Command) error {
			client := app.MustClient(ctx, cmd)

			port := app.MustPortOrRandom(ctx, cmd, 4318)

			ready := make(chan struct{})

			go func() {
				<-ready

				cli.Infof("export OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:%d", port)
				cli.Infof("export OTEL_EXPORTER_OTLP_PROTOCOL=http/protobuf")
			}()

			return client.ServicePortForward(ctx, "platform", "otel", "", map[int]int{port: 4318}, ready)
		},
	}
}
