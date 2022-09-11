package eksctl

import (
	"context"
)

func Create(ctx context.Context, name string, kubeconfig string, opt ...Option) error {
	e := New(opt...)

	args := []string{
		"create", "cluster",

		"--region", e.region,
		"--name", name,

		"--nodes", "2",
		"--node-type", "m5.large",
		"--node-private-networking",

		"--kubeconfig", kubeconfig,
	}

	return e.Invoke(ctx, args...)
}
