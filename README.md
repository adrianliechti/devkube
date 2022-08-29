# DevKube


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
scoop install kind kubectl helm adrianliechti/devkube
```


## Setup Cluster

### Using local Docker Engine

```shell
devkube create
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