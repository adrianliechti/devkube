package helm

import (
	"context"
	"fmt"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/kube"

	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func (h *Helm) Uninstall(ctx context.Context, release string) error {
	namespace := h.namespace

	if namespace == "" {
		namespace = "default"
	}

	log := func(format string, v ...interface{}) {
		//log.Printf(format, v)
	}

	config := new(action.Configuration)

	if err := config.Init(kube.GetConfig(h.kubeconfig, h.context, namespace), namespace, "", log); err != nil {
		return err
	}

	client := action.NewUninstall(config)

	result, err := client.Run(release)

	if err != nil {
		return err
	}

	fmt.Println("Successfully installed release", result.Info)

	return nil
}
