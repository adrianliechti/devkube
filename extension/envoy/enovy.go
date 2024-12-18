package envoy

import (
	"context"
	_ "embed"
	"strings"

	"github.com/adrianliechti/devkube/pkg/helm"
	"github.com/adrianliechti/loop/pkg/kubernetes"
)

const (
	name      = "envoy-gateway"
	namespace = "platform"

	// https://github.com/envoyproxy/gateway/releases
	chartName    = "oci://docker.io/envoyproxy/gateway-helm"
	chartVersion = "v1.2.4"
)

var (
	//go:embed manifest.yaml
	manifest string
)

func Ensure(ctx context.Context, client kubernetes.Client) error {
	values := map[string]any{}

	if err := helm.Ensure(ctx, client, namespace, name, "", chartName, chartVersion, values); err != nil {
		return err
	}

	if err := client.Apply(ctx, namespace, strings.NewReader(manifest)); err != nil {
		return err
	}

	return nil
}
