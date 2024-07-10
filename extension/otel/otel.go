package otel

import (
	"context"
	_ "embed"
	"strings"

	"github.com/adrianliechti/devkube/pkg/apply"
	"github.com/adrianliechti/devkube/pkg/helm"
	"github.com/adrianliechti/loop/pkg/kubernetes"
)

const (
	name      = "otel-collector"
	namespace = "platform"

	// https://artifacthub.io/packages/helm/opentelemetry-helm/opentelemetry-collector
	repoURL      = "https://open-telemetry.github.io/opentelemetry-helm-charts"
	chartName    = "opentelemetry-collector"
	chartVersion = "0.97.1"
)

var (
	//go:embed manifest.yaml
	manifest string
)

func Ensure(ctx context.Context, client kubernetes.Client) error {
	if err := ensureDaemonSet(ctx, client); err != nil {
		return err
	}

	if err := ensureDeployment(ctx, client); err != nil {
		return err
	}

	if err := apply.Apply(ctx, client, namespace, strings.NewReader(manifest)); err != nil {
		return err
	}

	return nil
}

func ensureDaemonSet(ctx context.Context, client kubernetes.Client) error {
	values := map[string]any{
		"nameOverride":     name,
		"fullnameOverride": name,

		"image": map[string]any{
			"repository": "otel/opentelemetry-collector-contrib",
		},

		"mode": "daemonset",

		"presets": map[string]any{
			"logsCollection": map[string]any{
				"enabled": true,

				"storeCheckpoints": true,
			},

			"hostMetrics": map[string]any{
				"enabled": true,
			},

			"kubernetesAttributes": map[string]any{
				"enabled": true,

				"extractAllPodLabels":      true,
				"extractAllPodAnnotations": true,
			},

			"kubeletMetrics": map[string]any{
				"enabled": true,
			},
		},

		"config": map[string]any{
			"receivers": map[string]any{
				"kubeletstats": map[string]any{
					"insecure_skip_verify": true,
				},
			},

			"exporters": collectorExporters(),

			"service": map[string]any{
				"pipelines": collectorPipelines(),
			},
		},
	}

	if err := helm.Ensure(ctx, client, namespace, name, repoURL, chartName, chartVersion, values); err != nil {
		return err
	}

	return nil
}

func ensureDeployment(ctx context.Context, client kubernetes.Client) error {
	values := map[string]any{
		"nameOverride":     name + "-cluster",
		"fullnameOverride": name + "-cluster",

		"image": map[string]any{
			"repository": "otel/opentelemetry-collector-contrib",
		},

		"mode": "deployment",

		"replicaCount": 1,

		"presets": map[string]any{
			"kubernetesEvents": map[string]any{
				"enabled": true,
			},

			"clusterMetrics": map[string]any{
				"enabled": true,
			},
		},

		"config": map[string]any{
			"exporters": collectorExporters(),

			"service": map[string]any{
				"pipelines": collectorPipelines(),
			},
		},
	}

	if err := helm.Ensure(ctx, client, namespace, name+"-cluster", repoURL, chartName, chartVersion, values); err != nil {
		return err
	}

	return nil
}

func collectorExporters() map[string]any {
	return map[string]any{
		"otlp": map[string]any{
			"endpoint": "tempo:4317",

			"tls": map[string]any{
				"insecure": true,
			},
		},

		"otlphttp": map[string]any{
			"endpoint": "http://loki:3100/otlp",
		},

		"prometheusremotewrite": map[string]any{
			"endpoint": "http://prometheus:9090/api/v1/write",
		},
	}
}

func collectorPipelines() map[string]any {
	return map[string]any{
		"traces": map[string]any{
			"receivers": []string{
				"otlp",
			},

			"exporters": []string{
				"otlp",
			},
		},

		"metrics": map[string]any{
			"receivers": []string{
				"otlp",
			},

			"exporters": []string{
				"prometheusremotewrite",
			},
		},

		"logs": map[string]any{
			"receivers": []string{
				"otlp",
			},

			"exporters": []string{
				"otlphttp",
			},
		},
	}
}
