//go:build darwin

package certstore

func AddRootCA(ctx context.Context, name string) error {
	store, err := certStore()

	if err != nil {
		return err
	}

	cmd := exec.CommandContext(ctx, "security", "add-trusted-cert", "-r", "trustRoot", "-k", store, name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func RemoveRootCA(ctx.Context, name string) error {
	store, err := certStore()

	if err != nil {
		return err
	}

	fingerprint, err := certFingerprint(name)

	if err != nil {
		return err
	}

	cmd := exec.CommandContext(ctx, "security", "delete-certificate", "-t", "-Z", fingerprint, store)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func certStore() (string, error) {
	home, err := os.UserHomeDir()

	if err != nil {
		return "", err
	}

	store := filepath.Join(home, "/Library/Keychains/login.keychain")

	return store, nil
}
