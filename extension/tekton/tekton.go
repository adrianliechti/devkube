package tekton

import (
	"context"
	_ "embed"

	"github.com/adrianliechti/loop/pkg/kubernetes"
)

func Ensure(ctx context.Context, client kubernetes.Client) error {
	if err := client.ApplyURL(ctx, "", "https://storage.googleapis.com/tekton-releases/pipeline/latest/release.yaml"); err != nil {
		return err
	}

	if err := client.ApplyURL(ctx, "", "https://storage.googleapis.com/tekton-releases/triggers/latest/release.yaml"); err != nil {
		return err
	}

	if err := client.ApplyURL(ctx, "", "https://storage.googleapis.com/tekton-releases/triggers/latest/interceptors.yaml"); err != nil {
		return err
	}

	if err := client.ApplyURL(ctx, "", "https://storage.googleapis.com/tekton-releases/dashboard/latest/release.yaml"); err != nil {
		return err
	}

	return nil
}
