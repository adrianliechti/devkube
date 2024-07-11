package extension

import (
	"context"

	"github.com/adrianliechti/loop/pkg/kubernetes"
)

type EnsureFunc = func(ctx context.Context, client kubernetes.Client) error
