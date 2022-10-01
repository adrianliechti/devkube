package tekton

import (
	"context"
	_ "embed"
	"strings"

	"github.com/adrianliechti/devkube/pkg/kubectl"
)

const (
	pipelineVersion   = "v0.40.1"
	pipelineNamespace = "tekton-pipelines"

	dashboardVersion   = "v0.29.2"
	dashboardNamespace = "tekton-pipelines"
)

var (
	//go:embed pipelineManifest.yaml
	pipelineManifest string

	//go:embed dashboardManifest.yaml
	dashboardManifest string
)

func Install(ctx context.Context, kubeconfig string) error {

	if err := installPipeline(ctx, kubeconfig); err != nil {
		return err
	}

	if err := installDashboard(ctx, kubeconfig); err != nil {
		return err
	}

	return nil
}

func Uninstall(ctx context.Context, kubeconfig string) error {
	if err := uninstallPipeline(ctx, kubeconfig); err != nil {
		return err
	}

	if err := uninstallDashboard(ctx, kubeconfig); err != nil {
		return err
	}

	return nil
}

func installPipeline(ctx context.Context, kubeconfig string) error {
	// manifest := pipelineManifest()

	if err := kubectl.Invoke(ctx, []string{"apply", "-f", "-"}, kubectl.WithKubeconfig(kubeconfig), kubectl.WithNamespace(pipelineNamespace), kubectl.WithInput(strings.NewReader(pipelineManifest)), kubectl.WithDefaultOutput()); err != nil {
		return err
	}

	return nil
}

func installDashboard(ctx context.Context, kubeconfig string) error {
	// manifest := dashboardManifest()

	if err := kubectl.Invoke(ctx, []string{"apply", "-f", "-"}, kubectl.WithKubeconfig(kubeconfig), kubectl.WithNamespace(dashboardNamespace), kubectl.WithInput(strings.NewReader(dashboardManifest)), kubectl.WithDefaultOutput()); err != nil {
		return err
	}

	return nil
}

func uninstallPipeline(ctx context.Context, kubeconfig string) error {
	// manifest := pipelineManifest()

	if err := kubectl.Invoke(ctx, []string{"delete", "-f", "-"}, kubectl.WithKubeconfig(kubeconfig), kubectl.WithNamespace(pipelineNamespace), kubectl.WithInput(strings.NewReader(pipelineManifest)), kubectl.WithDefaultOutput()); err != nil {
		// return err
	}

	return nil
}

func uninstallDashboard(ctx context.Context, kubeconfig string) error {
	// manifest := dashboardManifest()

	if err := kubectl.Invoke(ctx, []string{"delete", "-f", "-"}, kubectl.WithKubeconfig(kubeconfig), kubectl.WithNamespace(dashboardNamespace), kubectl.WithInput(strings.NewReader(dashboardManifest)), kubectl.WithDefaultOutput()); err != nil {
		// return err
	}

	return nil
}

// func pipelineManifest() string {
// 	return "https://storage.googleapis.com/tekton-releases/pipeline/previous/" + pipelineVersion + "/release.yaml"
// }

// func dashboardManifest() string {
// 	return "https://storage.googleapis.com/tekton-releases/dashboard/previous/" + dashboardVersion + "/tekton-dashboard-release.yaml"
// }
