//go:build linux

package certstore

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
)

func AddRootCA(ctx context.Context, name string) error {
	path, err := certPath(name)

	if err != nil {
		return err
	}

	data, err := os.ReadFile(name)

	if err != nil {
		return err
	}

	tee := exec.CommandContext(ctx, "sudo", "tee", path)
	tee.Stdin = bytes.NewReader(data)
	tee.Stdout = os.Stdout
	tee.Stderr = os.Stderr

	if err := tee.Run(); err != nil {
		return err
	}

	tool, err := certTool()

	if err != nil {
		return err
	}

	update := exec.CommandContext(ctx, "sudo", tool...)
	update.Stdin = os.Stdin
	update.Stdout = os.Stdout
	update.Stderr = os.Stderr

	return update.Run()
}

func RemoveRootCA(ctx context.Context, name string) error {
	path, err := certPath(name)

	if err != nil {
		return err
	}

	rm := exec.CommandContext(ctx, "sudo", "rm", "-f", path)
	rm.Stdin = os.Stdin
	rm.Stdout = os.Stdout
	rm.Stderr = os.Stderr

	if err := rm.Run(); err != nil {
		return err
	}

	tool, err := certTool()

	if err != nil {
		return err
	}

	update := exec.CommandContext(ctx, "sudo", tool...)
	update.Stdin = os.Stdin
	update.Stdout = os.Stdout
	update.Stderr = os.Stderr

	return update.Run()
}

func certPath(name string) (string, error) {
	fingerprint, err := certFingerprint(name)

	if err != nil {
		return "", err
	}

	if _, err := os.Stat("/etc/pki/ca-trust/source/anchors/"); !os.IsNotExist(err) {
		return fmt.Sprintf("/etc/pki/ca-trust/source/anchors/devkube-%s.pem", fingerprint), nil
	}

	if _, err := os.Stat("/usr/local/share/ca-certificates/"); !os.IsNotExist(err) {
		return fmt.Sprintf("/usr/local/share/ca-certificates/devkube-%s.crt", fingerprint), nil
	}

	if _, err := os.Stat("/etc/ca-certificates/trust-source/anchors/"); !os.IsNotExist(err) {
		return fmt.Sprintf("/etc/ca-certificates/trust-source/anchors/devkube-%s.crt", fingerprint), nil
	}

	if _, err := os.Stat("/usr/share/pki/trust/anchors/"); !os.IsNotExist(err) {
		return fmt.Sprintf("/usr/share/pki/trust/anchors/devkube-%s.pem", fingerprint), nil
	}

	return "", errors.New("could not determine certificate path")
}

func certTool() ([]string, error) {
	if _, err := os.Stat("/etc/pki/ca-trust/source/anchors/"); !os.IsNotExist(err) {
		return []string{"update-ca-trust", "extract"}, nil
	}

	if _, err := os.Stat("/usr/local/share/ca-certificates/"); !os.IsNotExist(err) {
		return []string{"update-ca-certificates"}, nil
	}

	if _, err := os.Stat("/etc/ca-certificates/trust-source/anchors/"); !os.IsNotExist(err) {
		return []string{"trust", "extract-compat"}, nil
	}

	if _, err := os.Stat("/usr/share/pki/trust/anchors/"); !os.IsNotExist(err) {
		return []string{"update-ca-certificates"}, nil
	}

	return nil, errors.New("could not determine certificate tool")
}
