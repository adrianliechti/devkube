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
	"strings"
	"sync"
	"time"

	"github.com/adrianliechti/devkube/app"
	"github.com/adrianliechti/devkube/pkg/cli"
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

			httpPort, err := system.FreePort(c.Int("http-port"))

			if err != nil {
				return err
			}

			httpsPort, err := system.FreePort(c.Int("https-port"))

			if err != nil {
				return err
			}

			client, err := kubernetes.NewFromConfig(kubeconfig)

			if err != nil {
				return err
			}

			cli.Info("Ingress availalbe at:")
			cli.Infof("  http://localhost:%d", httpPort)
			cli.Infof("  https://localhost:%d", httpsPort)

			return tunnelIngress(c.Context, client, httpPort, httpsPort)
		},
	}
}

func tunnelIngress(ctx context.Context, client kubernetes.Client, httpPort, httpsPort int) error {
	httpTunnel, err := system.FreePort(0)

	if err != nil {
		return err
	}

	httpsTunnel, err := system.FreePort(0)

	if err != nil {
		return err
	}

	secret, err := client.CoreV1().Secrets(DefaultNamespace).Get(ctx, "platform-ca", metav1.GetOptions{})

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

	httpListener, err := net.Listen("tcp", fmt.Sprintf(":%d", httpPort))

	if err != nil {
		return err
	}

	httpsListener, err := net.Listen("tcp", fmt.Sprintf(":%d", httpsPort))

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
		updateIngressHosts(ctx, client)
	}()

	go func() {
		target, _ := url.Parse(fmt.Sprintf("http://localhost:%d", httpTunnel))

		proxy := httputil.NewSingleHostReverseProxy(target)

		http.Serve(httpListener, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			proxy.ServeHTTP(w, r)
		}))
	}()

	go func() {
		target, _ := url.Parse(fmt.Sprintf("https://localhost:%d", httpsTunnel))

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

	if err := kubectl.Invoke(ctx, []string{"port-forward", "service/ingress-nginx-controller", fmt.Sprintf("%d:80", httpTunnel), fmt.Sprintf("%d:443", httpsTunnel)}, kubectl.WithKubeconfig(client.ConfigPath()), kubectl.WithNamespace(DefaultNamespace)); err != nil {
		return err
	}

	return nil
}

func updateIngressHosts(ctx context.Context, client kubernetes.Client) error {
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
			println("ingress added", ingress.Namespace, ingress.Name, strings.Join(hosts, ","))

		case watch.Deleted:
			println("ingress deleted", ingress.Namespace, ingress.Name, strings.Join(hosts, ","))
		}
	}

	return nil
}
