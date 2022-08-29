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

### Using Kubernetes-in-Docker

```shell
devkube create
```

## Open Kubernetes Dashboard

```shell
devkube dashboard
```

> Press "Skip" on the login page to access the dashboard as admin

## Open Observability Stack

```shell
devkube grafana
```


## Optional Features

### Trivy

```shell
devkube enable trivy
```