package app

import (
	"context"

	"github.com/adrianliechti/go-cli"
	"github.com/adrianliechti/loop/pkg/system"
)

var PortFlag = &cli.IntFlag{
	Name:  "port",
	Usage: "local port",
}

// PortOrRandom returns the explicitly requested port or, if none was given, a
// free one - preferring the passed port if it is still available.
func PortOrRandom(ctx context.Context, cmd *cli.Command, preference int) (int, error) {
	if port := cmd.Int(PortFlag.Name); port > 0 {
		return port, nil
	}

	return system.FreePort(preference)
}

func MustPortOrRandom(ctx context.Context, cmd *cli.Command, preference int) int {
	port, err := PortOrRandom(ctx, cmd, preference)

	if err != nil {
		cli.Fatal(err)
	}

	return port
}
