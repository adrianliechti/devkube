//go:build darwin || linux

package system

import (
	"os"
	"os/exec"
	"strings"
)

func IsElevated() (bool, error) {
	cmd := exec.Command("id", "-u")
	output, err := cmd.Output()

	if err != nil {
		return false, err
	}

	id := string(output)

	id = strings.TrimRight(id, "\n\r")
	id = strings.TrimSpace(id)

	if id == "0" {
		return true, nil
	}

	return false, nil
}

func RunElevated() error {
	exe, _ := os.Executable()
	cwd, _ := os.Getwd()

	args := []string{
		"-p",
		"[local sudo] Password: ",
		exe,
	}

	for _, arg := range os.Args[1:] {
		args = append(args, arg)
	}

	cmd := exec.Command("sudo", args...)
	cmd.Dir = cwd
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		return err
	}

	os.Exit(0)
	return nil
}
