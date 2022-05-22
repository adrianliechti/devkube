package app

import (
	"github.com/adrianliechti/devkube/pkg/cli"
)

var NameFlag = &cli.StringFlag{
	Name:  "name",
	Usage: "Cluster name",
}
