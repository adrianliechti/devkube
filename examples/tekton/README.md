# Tekont

## Install Tekton Dashboard, Pipelines & Triggers

```shell
$ devkube install tekton
```

### Open Dashboard

```shell
$ kubectl port-forward -n tekton-pipelines service/tekton-dashboard 9097:9097
```

open [http://localhost:9097](http://localhost:9097) in browser

### Sample Pipeline

```shell
$ kubectl apply -f pipeline/
```

```shell
$ kubectl apply -f pipeline-run/
```

### Sample Pipeline Trigger

```shell
$ kubectl apply -f pipeline-trigger/
```

```shell
$ kubectl port-forward service/el-hello-listener 8080

$ curl -v -H 'content-Type: application/json' -d '{"username": "Tekton"}' http://localhost:8080
```