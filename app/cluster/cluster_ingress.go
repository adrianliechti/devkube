package cluster

import (
	"fmt"

	"github.com/adrianliechti/devkube/app"
	"github.com/adrianliechti/devkube/pkg/cli"
	"github.com/adrianliechti/devkube/pkg/kubectl"
	"github.com/adrianliechti/devkube/pkg/kubernetes"
	"github.com/adrianliechti/devkube/pkg/system"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func IngressCommand() *cli.Command {
	return &cli.Command{
		Name:  "ingress",
		Usage: "Tunnel Ingress",

		Flags: []cli.Flag{
			app.ProviderFlag,
			app.ClusterFlag,
			&cli.IntFlag{
				Name:  "http-port",
				Usage: "Local HTTP Port",
				Value: 8080,
			},
			&cli.IntFlag{
				Name:  "https-port",
				Usage: "Local HTTPS Port",
				Value: 8443,
			},
		},

		Before: func(c *cli.Context) error {
			if _, _, err := kubectl.Info(c.Context); err != nil {
				return err
			}

			return nil
		},

		Action: func(c *cli.Context) error {
			provider, cluster := app.MustCluster(c)

			kubeconfig, closer := app.MustClusterKubeconfig(c, provider, cluster)
			defer closer()

			httpport, err := system.FreePort(c.Int("http-port"))

			if err != nil {
				return err
			}

			httpsport, err := system.FreePort(c.Int("https-port"))

			if err != nil {
				return err
			}

			client, err := kubernetes.NewFromConfig(kubeconfig)

			if err != nil {
				return err
			}

			secret, err := client.CoreV1().Secrets(DefaultNamespace).Get(c.Context, "platform-ca", metav1.GetOptions{})

			if err != nil {
				return err
			}

			key := secret.Data["tls.key"]
			cert := secret.Data["tls.crt"]

			_ = key
			_ = cert

			_ = httpport
			_ = httpsport

			// println(httpport)
			// println(httpsport)

			// println(string(key))
			// println(string(cert))

			if err := kubectl.Invoke(c.Context, []string{"port-forward", "service/ingress-nginx-controller", fmt.Sprintf("%d:80", httpport), fmt.Sprintf("%d:443", httpsport)}, kubectl.WithKubeconfig(kubeconfig), kubectl.WithNamespace(DefaultNamespace), kubectl.WithDefaultOutput()); err != nil {
				return err
			}

			return nil
		},
	}
}
