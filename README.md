# devkube

devkube bootraps a feature-rich Kubernetse cluster locally using Docker, or on a specified cloud provider on top of their managed Kubernetes offering.

Batteries included

- [Kubernetes Dashboard](https://kubernetes.io/docs/tasks/access-application-cluster/web-ui-dashboard/) - web-based user interface
- [Prometheus Operator](https://prometheus-operator.dev) - monitoring system
- [Grafana](https://grafana.com/grafana/) - data observability
- [Loki](https://grafana.com/oss/loki/) - log aggregation system
- [Tempo](https://grafana.com/oss/tempo/) - distributed tracing backend

Optional Add-ons

- [Falco](https://falco.org) - Kubernetes threat detection engine
- [Trivy](https://aquasecurity.github.io/trivy-operator/latest/) - Kubernetse workload vulnerability scanning


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


## Setup Cluster

### Using local Docker Engine

```shell
devkube create
```

### Using [AWS](https://aws.amazon.com/eks/) Cloud Provider

```shell
export AWS_ACCESS_KEY_ID=...
export AWS_SECRET_ACCESS_KEY=...
export AWS_DEFAULT_REGION=...

devkube create --provider aws
```

### Using [Azure](https://azure.microsoft.com/en-us/services/kubernetes-service/) Cloud Provider

```shell
export AZURE_TENANT_ID=...
export AZURE_SUBSCRIPTION_ID=...

devkube create --provider azure
```

### Using [DigitalOcean](https://www.digitalocean.com/products/kubernetes) Cloud Provider

```shell
export DIGITALOCEAN_TOKEN=...

devkube create --provider digitalocean
```

### Using [Linode](https://www.linode.com/) Cloud Provider

```shell
export LINODE_TOKEN=...

devkube create --provider linode
```

### Using [Vultr](https://www.vultr.com/) Cloud Provider

```shell
export VULTR_API_KEY=...

devkube create --provider vultr
```

## Administration Consoles

#### Kubernetes Dashboard

```shell
devkube dashboard
```

> Press "Skip" on the login page to access the dashboard as admin

#### Observability Stack

```shell
devkube grafana
```


## Optional Features

#### Trivy

```shell
devkube enable trivy
```

#### Falco

```shell
devkube enable falco
```