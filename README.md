# devkube

devkube bootstraps feature-rich Kubernetes clusters locally

## Create Cluster (using Docker)

The default `create` command will spin up a local DevKube cluster in a single container running in [Docker](https://docs.docker.com/get-docker/). [Lima](https://lima-vm.io), [Podman Desktop](https://podman-desktop.io), [Rancher Desktop](https://rancherdesktop.io) should work too.

```shell
$ devkube create
★ installing Kubernetes Cluster...
★ installing Cert-Manager...
★ installing Gatekeeper...
★ installing Crossplane...
★ installing Prometheus...
★ installing Grafana Loki...
★ installing Grafana Tempo...
★ installing Grafana...
★ installing Registry...
★ installing Dashboard...
★ installing Promtail...
★ installing OpenTelemetry...
```

Installed Components
- [Dashboard](https://kubernetes.io/docs/tasks/access-application-cluster/web-ui-dashboard/) - cluster management web ui
- [Registry](https://distribution.github.io/distribution/) - container images distribution
- [Cert-Manager](https://cert-manager.io) -X.509 certificate management
- [Gatekeeper](https://open-policy-agent.github.io/gatekeeper/website/) - customizable policy controller
- [Crossplane](https://www.crossplane.io) - universal control plane
- [Grafana](https://grafana.com/oss/grafana/) - observability and data visualization
- [Grafana Loki](https://grafana.com/oss/loki/) - log aggregation system
- [Grafana Tempo](https://grafana.com/oss/tempo/) - distributed tracing
- [Prometheus](https://prometheus.io) - monitoring and alerting system
- [Open Telemetry](https://opentelemetry.io/docs/collector/) - receive, process and export telemetry data


## Management Tools

### Kubernetes Web UI

open Dashboard in browser

```shell
$ devkube dashboard
```

### Grafana Web UI

open Grafana in browser

```shell
$ devkube grafana
```

### Access cluster workload

To access workload services within your cluster, `connect` allows you to forward these adresses and ports locally and allow easy access.

```shell
sudo devkube connect
...
5:58PM INF adding tunnel namespace=platform hosts="[dashboard.platform dashboard.platform.svc.cluster.local]" ports=[80]
5:58PM INF adding tunnel namespace=platform hosts="[grafana.platform grafana.platform.svc.cluster.local]" ports=[80]
...
```

### Import a local image

```shell
# pull an image (or build one)
docker pull alpine:3

# import image into cluster registry
devkube load alpine:3
```

### Build an image within cluster

```shell
cd /path/to/your/project

# cat Dockerfile
# FROM alpine:3
# RUN apk add --no-cache bash

# build image using buildkitd
devkube build demo .
★ creating container (default/loop-buildkit-704c031)...
★ copying build context...
...
★ removing container (default/loop-buildkit-704c031)..

# run impage in kubernetes
kubectl run -it --rm demo --image registry.platform/demo /bin/bash
demo:/# exit
```


## Installation

MacOS / Linux with [Homebrew](https://brew.sh)

```shell
brew install adrianliechti/tap/devkube
```

Windows with [Scoop](https://scoop.sh)

```shell
scoop bucket add adrianliechti https://github.com/adrianliechti/scoop-bucket
scoop install kubectl helm adrianliechti/devkube
```