//go:build linux

package certstore

func AddRootCA(ctx context.Context, name string) error {
	return errors.New("Adding Root CA on Linux is currently not supprted")
}

func RemoveRootCA(ctx context.Context, name string) error {
	return errors.New("Removing Root CA on Linux is currently not supprted")
}