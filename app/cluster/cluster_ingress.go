package cluster

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/adrianliechti/devkube/app"
	"github.com/adrianliechti/devkube/pkg/cli"
	"github.com/adrianliechti/devkube/pkg/hostsfile"
	"github.com/adrianliechti/devkube/pkg/kubectl"
	"github.com/adrianliechti/devkube/pkg/kubernetes"
	"github.com/adrianliechti/devkube/pkg/system"

	"github.com/samber/lo"

	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

func IngressCommand() *cli.Command {
	return &cli.Command{
		Name:  "ingress",
		Usage: "Tunnel Ingress",

		Category: app.ConnectCategory,

		Flags: []cli.Flag{
			app.ProviderFlag,
			app.ClusterFlag,
		},

		Before: func(c *cli.Context) error {
			if _, _, err := kubectl.Info(c.Context); err != nil {
				return err
			}

			return nil
		},

		Action: func(c *cli.Context) error {
			elevated, err := system.IsElevated()

			if err != nil {
				return err
			}

			if !elevated {
				if err := system.RunElevated(); err != nil {
					cli.Fatal("This command must be run as root!")
				}

				os.Exit(0)
			}

			provider, cluster := app.MustCluster(c)

			kubeconfig, closer := app.MustClusterKubeconfig(c, provider, cluster)
			defer closer()

			client, err := kubernetes.NewFromConfig(kubeconfig)

			if err != nil {
				return err
			}

			address := "127.244.0.0"

			cli.Info("Ingress availalbe at:")
			cli.Infof("  http://%s:80", address)
			cli.Infof("  https://%s:443", address)

			return tunnelIngress(c.Context, client, address, 80, 443)
		},
	}
}

func tunnelIngress(ctx context.Context, client kubernetes.Client, address string, httpPort, httpsPort int) error {
	if err := system.AliasIP(ctx, address); err != nil {
		return err
	}

	defer func() {
		system.UnaliasIP(context.Background(), address)
		hostsfile.RemoveByAddress(address)
	}()

	httpTunnel := 5080
	httpsTunnel := 5443

	secret, err := client.CoreV1().Secrets(app.DefaultNamespace).Get(ctx, "platform-ca", metav1.GetOptions{})

	if err != nil {
		return err
	}

	ca, err := tls.X509KeyPair(secret.Data["tls.crt"], secret.Data["tls.key"])

	if err != nil {
		return err
	}

	cacert, err := x509.ParseCertificate(ca.Certificate[0])

	if err != nil {
		return err
	}

	capool := x509.NewCertPool()
	capool.AddCert(cacert)

	timestamp := time.Now()

	httpListener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", address, httpPort))

	if err != nil {
		return err
	}

	httpsListener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", address, httpsPort))

	if err != nil {
		return err
	}

	tlsConfig := &tls.Config{
		RootCAs: capool,
	}

	certificates := map[string]*tls.Certificate{}
	var certificatesLock sync.Mutex

	tlsConfig.GetCertificate = func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
		if cert, ok := certificates[info.ServerName]; ok {
			return cert, nil
		}

		//println("generate certificate", info.ServerName)

		certificatesLock.Lock()
		defer certificatesLock.Unlock()

		template := &x509.Certificate{
			SerialNumber: big.NewInt(timestamp.Unix()),

			// Subject: pkix.Name{
			// 	CommonName: info.ServerName,
			// },

			DNSNames: []string{
				info.ServerName,
			},

			NotBefore: timestamp,
			NotAfter:  timestamp.AddDate(1, 0, 0),

			KeyUsage:    x509.KeyUsageDigitalSignature,
			ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		}

		privateKey, err := rsa.GenerateKey(rand.Reader, 4096)

		if err != nil {
			return nil, err
		}

		certificate, err := x509.CreateCertificate(rand.Reader, template, cacert, &privateKey.PublicKey, ca.PrivateKey)

		if err != nil {
			return nil, err
		}

		pair := &tls.Certificate{
			Certificate: [][]byte{certificate},
			PrivateKey:  privateKey,
		}

		certificates[info.ServerName] = pair
		return pair, nil
	}

	httpsListener = tls.NewListener(httpsListener, tlsConfig)

	go func() {
		updateIngressHosts(ctx, client, address)
	}()

	go func() {
		target, _ := url.Parse(fmt.Sprintf("http://%s:%d", address, httpTunnel))

		proxy := httputil.NewSingleHostReverseProxy(target)

		http.Serve(httpListener, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			proxy.ServeHTTP(w, r)
		}))
	}()

	go func() {
		target, _ := url.Parse(fmt.Sprintf("https://%s:%d", address, httpsTunnel))

		transport := http.DefaultTransport.(*http.Transport).Clone()
		transport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}

		proxy := httputil.NewSingleHostReverseProxy(target)
		proxy.Transport = transport

		http.Serve(httpsListener, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			proxy.ServeHTTP(w, r)
		}))
	}()

	if err := kubectl.Invoke(ctx, []string{"port-forward", "service/ingress-nginx-controller", "--address", address, fmt.Sprintf("%d:80", httpTunnel), fmt.Sprintf("%d:443", httpsTunnel)}, kubectl.WithKubeconfig(client.ConfigPath()), kubectl.WithNamespace(app.DefaultNamespace)); err != nil {
		return err
	}

	return nil
}

func updateIngressHosts(ctx context.Context, client kubernetes.Client, address string) error {
	watcher, err := client.NetworkingV1().Ingresses("").Watch(ctx, metav1.ListOptions{})

	if err != nil {
		return err
	}

	for event := range watcher.ResultChan() {
		ingress, ok := event.Object.(*networkingv1.Ingress)

		if !ok {
			continue
		}

		hosts := make([]string, 0)

		for _, rule := range ingress.Spec.Rules {
			hosts = append(hosts, rule.Host)
		}

		hosts = lo.Uniq(hosts)

		switch event.Type {
		case watch.Added:
			hostsfile.AddAlias(address, hosts...)
			println("added", ingress.Namespace, ingress.Name, strings.Join(hosts, ","))

		case watch.Deleted:
			hostsfile.RemoveByAlias(hosts...)
			println("ingress deleted", ingress.Namespace, ingress.Name, strings.Join(hosts, ","))
		}
	}

	return nil
}
