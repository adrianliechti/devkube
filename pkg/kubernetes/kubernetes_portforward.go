package kubernetes

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
)

func (c *client) ServicePortForward(ctx context.Context, namespace, name, address string, ports map[int]int, readyChan chan struct{}) error {
	pod, err := c.ServicePod(ctx, namespace, name)

	if err != nil {
		return err
	}

	return c.PodPortForward(ctx, pod.Namespace, pod.Name, address, ports, readyChan)
}

func (c *client) PodPortForward(ctx context.Context, namespace, name, address string, ports map[int]int, readyChan chan struct{}) error {
	if address == "" {
		address = "localhost"
	}

	path := fmt.Sprintf("/api/v1/namespaces/%s/pods/%s/portforward", namespace, name)

	host := c.Config().Host
	host = strings.TrimPrefix(host, "http://")
	host = strings.TrimPrefix(host, "https://")

	transport, upgrader, err := spdy.RoundTripperFor(c.Config())

	if err != nil {
		return err
	}

	mappings := make([]string, 0)

	for s, t := range ports {
		mappings = append(mappings, fmt.Sprintf("%d:%d", s, t))
	}

	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, http.MethodPost, &url.URL{Scheme: "https", Path: path, Host: host})
	forwarder, err := portforward.NewOnAddresses(dialer, []string{address}, mappings, ctx.Done(), readyChan, ioutil.Discard, ioutil.Discard)

	if err != nil {
		return err
	}

	return forwarder.ForwardPorts()
}
