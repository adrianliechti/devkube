package app

import (
	"errors"
	"strconv"
	"strings"

	"github.com/adrianliechti/loop/pkg/cli"
	"github.com/adrianliechti/loop/pkg/system"
)

var PortFlag = &cli.IntFlag{
	Name:  "port",
	Usage: "port",
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

var PortsFlag = &cli.StringSliceFlag{
	Name:  "port",
	Usage: "port mappings",
}

func Ports(c *cli.Context) (map[int]int, error) {
	s := c.StringSlice(PortsFlag.Name)

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

func MustPorts(c *cli.Context) map[int]int {
	ports, err := Ports(c)

	if err != nil {
		cli.Fatal(err)
	}

	if len(ports) == 0 {
		cli.Fatal(errors.New("ports missing"))
	}

	return ports
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
