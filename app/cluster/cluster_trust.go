package cluster

import (
	"context"
	"crypto/sha1"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/adrianliechti/devkube/app"
	"github.com/adrianliechti/devkube/pkg/cli"
	"github.com/adrianliechti/devkube/pkg/kubernetes"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TrustCommand() *cli.Command {
	return &cli.Command{
		Name:  "trust",
		Usage: "Trust Cluster Root CA",

		Flags: []cli.Flag{
			app.ProviderFlag,
			app.ClusterFlag,
		},

		Action: func(c *cli.Context) error {
			provider, cluster := app.MustCluster(c)

			kubeconfig, closer := app.MustClusterKubeconfig(c, provider, cluster)
			defer closer()

			client, err := kubernetes.NewFromConfig(kubeconfig)

			if err != nil {
				return err
			}

			secret, err := client.CoreV1().Secrets("cert-manager").Get(c.Context, "platform-ca", metav1.GetOptions{})

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

			if err := installCertificate(c.Context, file); err != nil {
				return err
			}

			return nil
		},
	}
}

func installCertificate(ctx context.Context, name string) error {
	keychain, err := userKeychain()

	if err != nil {
		return err
	}

	if err := exec.CommandContext(ctx, "security", "add-trusted-cert", "-r", "trustRoot", "-k", keychain, name).Run(); err != nil {
		return err
	}

	return nil
}

func uninstallCertificate(ctx context.Context, name string) error {
	keychain, err := userKeychain()

	if err != nil {
		return err
	}

	data, err := os.ReadFile(name)

	if err != nil {
		return err
	}

	block, _ := pem.Decode(data)
	cert, err := x509.ParseCertificate(block.Bytes)

	if err != nil {
		return err
	}

	hash := sha1.Sum(cert.Raw)
	fingerprint := strings.ToUpper(hex.EncodeToString(hash[:]))

	if err := exec.CommandContext(ctx, "security", "delete-certificate", "-t", "-Z", fingerprint, keychain).Run(); err != nil {
		return err
	}

	return nil
}

func userKeychain() (string, error) {
	home, err := os.UserHomeDir()

	if err != nil {
		return "", err
	}

	keychain := filepath.Join(home, "/Library/Keychains/login.keychain")

	return keychain, nil
}
