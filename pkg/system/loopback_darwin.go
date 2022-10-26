//go:build darwin

package system

import (
	"context"
	"errors"
	"os/exec"
)

func AliasIP(ctx context.Context, alias string) error {
	output, err := exec.CommandContext(ctx, "ifconfig", "lo0", "alias", alias).CombinedOutput()

	if err != nil {
		return errors.New(string(output))
	}

	return nil
}

func UnaliasIP(ctx context.Context, alias string) error {
	output, err := exec.CommandContext(ctx, "ifconfig", "lo0", "-alias", alias).CombinedOutput()

	if err != nil {
		return errors.New(string(output))
	}

	return nil
}
