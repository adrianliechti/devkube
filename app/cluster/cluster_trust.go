package cluster

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/adrianliechti/devkube/app"
	"github.com/adrianliechti/devkube/pkg/certstore"
	"github.com/adrianliechti/devkube/pkg/cli"
	"github.com/adrianliechti/devkube/pkg/kubernetes"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TrustCommand() *cli.Command {
	return &cli.Command{
		Name:  "trust",
		Usage: "Trust Root CA",

		Category: app.ConnectCategory,

		Flags: []cli.Flag{
			app.ProviderFlag,
			app.ClusterFlag,
			&cli.BoolFlag{
				Name:  "uninstall",
				Usage: "Uninstall Cluster Root CA",
			},
		},

		Action: func(c *cli.Context) error {
			provider, cluster := app.MustCluster(c)

			kubeconfig, closer := app.MustClusterKubeconfig(c, provider, cluster)
			defer closer()

			client, err := kubernetes.NewFromConfig(kubeconfig)

			if err != nil {
				return err
			}

			secret, err := client.CoreV1().Secrets(app.DefaultNamespace).Get(c.Context, "platform-ca", metav1.GetOptions{})

			if err != nil {
				return err
			}

			data := secret.Data["ca.crt"]

			if len(data) == 0 {
				return errors.New("invalid certificate data")
			}

			dir, err := os.MkdirTemp("", "devkube")

			if err != nil {
				return err
			}

			defer os.RemoveAll(dir)

			file := filepath.Join(dir, "ca.crt")

			if err := os.WriteFile(file, data, 0644); err != nil {
				return err
			}

			if c.Bool("uninstall") {
				return certstore.RemoveRootCA(c.Context, file)
			}

			return certstore.AddRootCA(c.Context, file)
		},
	}
}
