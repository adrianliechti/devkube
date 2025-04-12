package app

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/adrianliechti/go-cli"
	"github.com/adrianliechti/loop/pkg/system"
)

var PortFlag = &cli.IntFlag{
	Name:  "port",
	Usage: "port",
}

func Port(ctx context.Context, cmd *cli.Command) int {
	return int(cmd.Int(PortFlag.Name))
}

func MustPort(ctx context.Context, cmd *cli.Command) int {
	port := Port(ctx, cmd)

	if port <= 0 {
		cli.Fatal(errors.New("port missing"))
	}

	return port
}

var PortsFlag = &cli.StringSliceFlag{
	Name:  "port",
	Usage: "port mappings",
}

func Ports(ctx context.Context, cmd *cli.Command) (map[int]int, error) {
	s := cmd.StringSlice(PortsFlag.Name)

	result := map[int]int{}

	for _, p := range s {
		pair := strings.Split(p, ":")

		if len(pair) > 2 {
			return nil, errors.New("invalid port mapping")
		}

		if len(pair) == 1 {
			pair = []string{pair[0], pair[0]}
		}

		source, err := strconv.Atoi(pair[0])

		if err != nil {
			return nil, err
		}

		target, err := strconv.Atoi(pair[1])

		if err != nil {
			return nil, err
		}

		result[source] = target
	}

	return result, nil
}

func MustPorts(ctx context.Context, cmd *cli.Command) map[int]int {
	ports, err := Ports(ctx, cmd)

	if err != nil {
		cli.Fatal(err)
	}

	if len(ports) == 0 {
		cli.Fatal(errors.New("ports missing"))
	}

	return ports
}

func PortOrRandom(ctx context.Context, cmd *cli.Command, preference int) (int, error) {
	port := Port(ctx, cmd)

	if port > 0 {
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

func RandomPort(ctx context.Context, cmd *cli.Command, preference int) (int, error) {
	return system.FreePort(preference)
}

func MustRandomPort(ctx context.Context, cmd *cli.Command, preference int) int {
	port, err := RandomPort(ctx, cmd, preference)

	if err != nil {
		cli.Fatal(err)
	}

	return port
}
