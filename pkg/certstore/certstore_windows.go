//go:build windows

package certstore

import (
	"context"
	"os"
	"os/exec"
)

func AddRootCA(ctx context.Context, name string) error {
	cmd := exec.CommandContext(ctx, "certutil", "-user", "-addstore", "Root", name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func RemoveRootCA(ctx context.Context, name string) error {
	fingerprint, err := certFingerprint(name)

	if err != nil {
		return err
	}

	cmd := exec.CommandContext(ctx, "certutil", "-user", "-delstore", "Root", fingerprint)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
