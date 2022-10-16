package cluster

import (
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
	"sync"
	"time"

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

			httptunnel, err := system.FreePort(0)

			if err != nil {
				return err
			}

			httpstunnel, err := system.FreePort(0)

			if err != nil {
				return err
			}

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

			httpListener, err := net.Listen("tcp", fmt.Sprintf(":%d", httpport))

			if err != nil {
				return err
			}

			httpsListener, err := net.Listen("tcp", fmt.Sprintf(":%d", httpsport))

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
				target, _ := url.Parse(fmt.Sprintf("http://localhost:%d", httptunnel))

				proxy := httputil.NewSingleHostReverseProxy(target)

				http.Serve(httpListener, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					proxy.ServeHTTP(w, r)
				}))
			}()

			go func() {
				target, _ := url.Parse(fmt.Sprintf("https://localhost:%d", httpstunnel))

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

			cli.Info("Ingress availalbe at:")
			cli.Infof("  http://localhost:%d", httpport)
			cli.Infof("  https://localhost:%d", httpsport)

			if err := kubectl.Invoke(c.Context, []string{"port-forward", "service/ingress-nginx-controller", fmt.Sprintf("%d:80", httptunnel), fmt.Sprintf("%d:443", httpstunnel)}, kubectl.WithKubeconfig(kubeconfig), kubectl.WithNamespace(DefaultNamespace)); err != nil {
				return err
			}

			return nil
		},
	}
}
