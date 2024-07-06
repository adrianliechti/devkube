package certmanager

import (
	"context"
	_ "embed"
	"strings"

	"github.com/adrianliechti/devkube/pkg/apply"
	"github.com/adrianliechti/loop/pkg/kubernetes"
)

var (
	//go:embed manifest.yaml
	manifest string
)

const (
	version   = "1.15.1"
	namespace = "cert-manager"
)

func Ensure(ctx context.Context, client kubernetes.Client) error {
	if err := apply.ApplyURL(ctx, client, namespace, "https://github.com/cert-manager/cert-manager/releases/download/v"+version+"/cert-manager.yaml"); err != nil {
		return err
	}

	if err := apply.Apply(ctx, client, namespace, strings.NewReader(manifest)); err != nil {
		return err
	}

	return nil
}
