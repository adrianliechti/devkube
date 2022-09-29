package eksctl

import (
	"context"
)

func Delete(ctx context.Context, name string, opt ...Option) error {
	e := New(opt...)

	args := []string{
		"delete", "cluster",

		"--region", e.region,
		"--name", name,
	}

	return e.Invoke(ctx, args...)
}
