package eksctl

import (
	"context"
)

func Export(ctx context.Context, name, path string, opt ...Option) error {
	e := New(opt...)

	args := []string{
		"utils", "write-kubeconfig",

		"--region", e.region,
		"--cluster", name,

		"--kubeconfig", path,
	}

	return e.Invoke(ctx, args...)
}
