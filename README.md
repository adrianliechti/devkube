# devkube

devkube bootstraps feature-rich Kubernetes clusters locally using Docker or on a specified cloud provider on top of their managed Kubernetes offering.

## Batteries included

- [Registry](https://github.com/distribution/distribution) - image distribution
- [Dashboard](https://kubernetes.io/docs/tasks/access-application-cluster/web-ui-dashboard/) - web-based user interface
- [Cert-Manager](https://cert-manager.io)- certificate management
- [Ingress](https://kubernetes.github.io/ingress-nginx/) - NGINX Ingress Controller
- [Grafana](https://grafana.com/grafana/) - data observability
- [Prometheus](https://prometheus-operator.dev) - monitoring system
- [Loki](https://grafana.com/oss/loki/) - log aggregation system
- [Tempo](https://grafana.com/oss/tempo/) - distributed tracing backend

### Optional Add-ons

- [Linkerd](https://linkerd.io) - Service Mesh
- [Kyverno](https://kyverno.io) - Kubernetes Policy Management
- [Falco](https://falco.org) - Kubernetes threat detection engine
- [Trivy](https://aquasecurity.github.io/trivy-operator/latest/) - Kubernetes workload vulnerability scanning

### Cloud providers

- [AWS (Alpha)](https://aws.amazon.com/eks/)
- [Azure (Beta)](https://azure.microsoft.com/en-us/services/kubernetes-service/)
- [DigitalOcean (Alpha)](https://www.digitalocean.com/products/kubernetes)
- [Linode (Alpha)](https://www.linode.com/products/kubernetes/)
- [Vultr (Alpha)](https://www.vultr.com/kubernetes/)

![Overview](docs/assets/overview.svg)

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/) - Container daemon
- [Kind](https://kind.sigs.k8s.io/) - Kubernetes in Docker, for local cluster

## Install

MacOS / Linux with [Homebrew](https://brew.sh)

```shell
brew install adrianliechti/tap/devkube
```

Windows with [Scoop](https://scoop.sh)

```shell
scoop bucket add adrianliechti https://github.com/adrianliechti/scoop-bucket
scoop install kubectl helm adrianliechti/devkube
```

## Create Cluster

```shell
devkube create
```

![Cluster](docs/assets/cluster.png)

## Access Dashboard

```shell
devkube dashboard
```

![Dashboard](docs/assets/dashboard.png)

> Press "Skip" on the login page to access the dashboard as admin

## Access Grafana

```shell
devkube grafana
```

![Grafana](docs/assets/grafana.png)

## Advanced Features

### Ingress Controller

![Ingress](docs/assets/ingress.png)

This CLI can forward traffic to the ingress controller and simulate DNS by adding entries in `/etc/hosts` temporary. It also allows to trust the pre-configured certificate authority (CA) to support TLS rules.

```shell
# Trust Platform CA (use --uninstall to remove)
devkube trust

# Tunnel Traffic (needs sudo)
devkube ingress
```

### OpenTelemetry

```mermaid
flowchart LR
    A[App] -->|OTLP| B(Collector<br>telemetry.loop)
    B --> C{Forward}
    C -->|Logs| D[Loki<br>loki.loop]
    C -->|Traces| E[Tempo<br>tempo.loop]
    C -->|Metrics| F[Prometheus<br>prometheus.loop]
    D <--- G((Grafana))
    E <--- G
    F <--- G
```

![OpenTelemetry](docs/assets/otel.png)

### Trivy

Trivy is a comprehensive security scanner. It is reliable, fast, extremely easy to use, and it works wherever you need it.

```shell
devkube enable trivy
```

![Trivy](docs/assets/trivy.png)

### Kyverno

Kyverno is a policy engine designed for Kubernetes. With Kyverno, policies are managed as Kubernetes resources and no new language is required to write policies.

```shell
devkube enable kyverno
```

![Trivy](docs/assets/kyverno.png)

### Falco

The Falco Project is a cloud native runtime security tool. Falco makes it easy to consume kernel events, and enrich those events with information from Kubernetes and the rest of the cloud native stack.

```shell
devkube enable falco
```

![Falco](docs/assets/falco.png)

### Linkerd

Linkerd is a service mesh for Kubernetes. It makes running services easier and safer by giving you runtime debugging, observability, reliability, and security—all without requiring any changes to your code.

```shell
devkube enable linkerd
```

#### Install CLI

MacOS / Linux with [Homebrew](https://brew.sh)

```shell
brew install linkerd
```

Windows with [Scoop](https://scoop.sh)

```shell
scoop install linkerd
```

Open Dashboard

```shell
linkerd viz dashboard
```

![Linkerd](docs/assets/linkerd.png)

![Linkerd Grafana](docs/assets/linkerd_grafana.png)
