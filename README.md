# devkube

devkube bootstraps feature-rich Kubernetes clusters locally

## Create Cluster (using Docker)

The default `create` command will spin up a local DevKube cluster in a single container running in [Docker](https://docs.docker.com/get-docker/). [Lima](https://lima-vm.io), [Podman Desktop](https://podman-desktop.io), [Rancher Desktop](https://rancherdesktop.io) should work too.

```shell
$ devkube create
★ installing Kubernetes Cluster...
★ installing Cert-Manager...
★ installing Registry...
★ installing Gatekeeper...
★ installing OpenTelemetry...
```

Installed Components
- [Cert-Manager](https://cert-manager.io) - X.509 Certificate Management
- [Registry](https://distribution.github.io/distribution/) - Container Images Distribution
- [Gatekeeper](https://open-policy-agent.github.io/gatekeeper/website/) - Customizable Policy controller
- [Grafana LGTM](ithub.com/grafana/docker-otel-lgtm) - Observability and Data Visualization

Optional Components
- [Crossplane](https://www.crossplane.io) - Universal Control Plane
- [Envoy Gateway](https://gateway.envoyproxy.io) - Manage Application and API traffic
- [Argo CD](https://argo-cd.readthedocs.io) - Declarative GitOps CD
- [Tekton](https://tekton.dev) - Cloud Native CI/CD

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
5:58PM INF adding tunnel address=127.244.179.12 hosts="[grafana.platform grafana.platform.svc.cluster.local]" ports=[3000]
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