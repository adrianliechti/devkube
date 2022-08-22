package dashboard

import (
	"context"

	"github.com/adrianliechti/devkube/pkg/helm"
	"github.com/adrianliechti/devkube/pkg/kubectl"
	"github.com/adrianliechti/devkube/pkg/kubernetes"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	dashboardRepo = "https://kubernetes.github.io/dashboard"

	dashboard        = "dashboard"
	dashboardChart   = "kubernetes-dashboard"
	dashboardVersion = "5.8.0"
)

func Install(ctx context.Context, kubeconfig, namespace string) error {
	if namespace == "" {
		namespace = "default"
	}

	client, err := kubernetes.NewFromConfig(kubeconfig)

	if err != nil {
		return err
	}

	values := map[string]any{
		"nameOverride":     dashboard,
		"fullnameOverride": dashboard,

		"extraArgs": []string{
			"--enable-skip-login",
			"--enable-insecure-login",
			"--disable-settings-authorizer",
		},

		"protocolHttp": true,

		"service": map[string]any{
			"externalPort": 80,
		},

		"serviceMonitor": map[string]any{
			"enabled": true,
		},

		"metricsScraper": map[string]any{
			"enabled": true,
		},

		"resources": nil,
	}

	if err := helm.Install(ctx, dashboard, dashboardRepo, dashboardChart, dashboardVersion, values, helm.WithKubeconfig(kubeconfig), helm.WithNamespace(namespace), helm.WithDefaultOutput()); err != nil {
		return err
	}

	client.RbacV1().ClusterRoleBindings().Delete(ctx, dashboard, metav1.DeleteOptions{})

	if _, err := client.RbacV1().ClusterRoleBindings().Create(ctx, &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: dashboard,
		},

		Subjects: []rbacv1.Subject{
			{
				Kind: rbacv1.ServiceAccountKind,

				Name:      dashboard,
				Namespace: namespace,
			},
		},

		RoleRef: rbacv1.RoleRef{
			Kind: "ClusterRole",
			Name: "cluster-admin",
		},
	}, metav1.CreateOptions{}); err != nil {
		return err
	}

	return nil
}

func Uninstall(ctx context.Context, kubeconfig, namespace string) error {
	if namespace == "" {
		namespace = "default"
	}

	if err := kubectl.Invoke(ctx, []string{"delete", "clusterrolebinding", dashboard}, kubectl.WithKubeconfig(kubeconfig)); err != nil {
		//return err
	}

	if err := helm.Uninstall(ctx, dashboard, helm.WithKubeconfig(kubeconfig), helm.WithNamespace(namespace)); err != nil {
		//return err
	}

	return nil
}
