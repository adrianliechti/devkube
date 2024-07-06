package gatekeeper

import (
	"context"

	"github.com/adrianliechti/devkube/pkg/apply"
	"github.com/adrianliechti/loop/pkg/kubernetes"
)

const (
	namespace = "gatekeeper-system"

	version = "v3.16.3"
)

func Ensure(ctx context.Context, client kubernetes.Client) error {
	if err := apply.ApplyURL(ctx, client, namespace, "https://raw.githubusercontent.com/open-policy-agent/gatekeeper/"+version+"/deploy/gatekeeper.yaml"); err != nil {
		return err
	}

	return nil
}
