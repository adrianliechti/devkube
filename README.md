# devkube

devkube bootstraps feature-rich Kubernetes clusters locally using Docker or on a specified cloud provider on top of their managed Kubernetes offering.

## Batteries included

- [Registry](https://github.com/distribution/distribution) - image distribution
- [Dashboard](https://kubernetes.io/docs/tasks/access-application-cluster/web-ui-dashboard/) - web-based user interface
- [Cert-Manager](https://cert-manager.io)- certificate management
- [Grafana](https://grafana.com/grafana/) - data observability
- [Prometheus](https://prometheus-operator.dev) - monitoring system
- [Loki](https://grafana.com/oss/loki/) - log aggregation system
- [Tempo](https://grafana.com/oss/tempo/) - distributed tracing backend

### Optional Add-ons

- [Falco](https://falco.org) - Kubernetes threat detection engine
- [Trivy](https://aquasecurity.github.io/trivy-operator/latest/) - Kubernetse workload vulnerability scanning

### Cloud providers

- [AWS (Alpha)](https://aws.amazon.com/eks/)
- [Azure (Beta)](https://azure.microsoft.com/en-us/services/kubernetes-service/)
- [DigitalOcean (Alpha)](https://www.digitalocean.com/products/kubernetes)
- [Linode (Alpha)](https://www.linode.com/products/kubernetes/)
- [Vultr (Alpha)](https://www.vultr.com/kubernetes/)

![Overview](docs/assets/overview.svg)

## Install

#### MacOS / Linux

[Homebrew](https://brew.sh)

```
brew install adrianliechti/tap/devkube
```

#### Windows

[Scoop](https://scoop.sh)

```shell
scoop bucket add adrianliechti https://github.com/adrianliechti/scoop-bucket
scoop install kubectl helm adrianliechti/devkube
```


## Create Cluster

![Cluster](docs/assets/cluster.png)

```shell
devkube create
```

## Access Dashboard

![Dashboard](docs/assets/dashboard.png)

```shell
devkube dashboard
```

> Press "Skip" on the login page to access the dashboard as admin

## Access Grafana

![Grafana](docs/assets/grafana.png)

```shell
devkube grafana
```

## Advanced Features

### OpenTelememetry

![OpenTelemetry](docs/assets/otel.png)


### Trivy

Trivy is a comprehensive security scanner. It is reliable, fast, extremely easy to use, and it works wherever you need it.

![Trivy](docs/assets/trivy.png)

```shell
devkube enable trivy
```

#### Falco

The Falco Project is a cloud native runtime security tool. Falco makes it easy to consume kernel events, and enrich those events with information from Kubernetes and the rest of the cloud native stack.

```shell
devkube enable falco
```