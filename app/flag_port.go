package app

import (
	"errors"

	"github.com/adrianliechti/devkube/pkg/cli"
	"github.com/adrianliechti/devkube/pkg/system"
)

var PortFlag = &cli.IntFlag{
	Name:  "port",
	Usage: "Local Port",
}

func Port(c *cli.Context) int {
	return c.Int(PortFlag.Name)
}

func MustPort(c *cli.Context) int {
	port := Port(c)

	if port <= 0 {
		cli.Fatal(errors.New("port missing"))
	}

	return port
}

func PortOrRandom(c *cli.Context, preference int) (int, error) {
	port := Port(c)

	if port > 0 {
		return port, nil
	}

	return system.FreePort(preference)
}

func MustPortOrRandom(c *cli.Context, preference int) int {
	port, err := PortOrRandom(c, preference)

	if err != nil {
		cli.Fatal(err)
	}

	return port
}

func RandomPort(c *cli.Context, preference int) (int, error) {
	return system.FreePort(preference)
}

func MustRandomPort(c *cli.Context, preference int) int {
	port, err := RandomPort(c, preference)

	if err != nil {
		cli.Fatal(err)
	}

	return port
}
