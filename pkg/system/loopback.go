//go:build !darwin

package system

import (
	"context"
)

func AliasIP(ctx context.Context, alias string) error {
	return nil
}

func UnaliasIP(ctx context.Context, alias string) error {
	return nil
}
