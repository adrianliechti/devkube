

```shell
$ kubectl apply -f https://storage.googleapis.com/tekton-releases/pipeline/latest/release.yaml
$ kubectl apply -f https://storage.googleapis.com/tekton-releases/dashboard/latest/release.yaml

$ kubectl apply -f https://storage.googleapis.com/tekton-releases/triggers/latest/release.yaml
$ kubectl apply -f https://storage.googleapis.com/tekton-releases/triggers/latest/interceptors.yaml
```

```shell
$ kubectl port-forward -n tekton-pipelines service/tekton-dashboard 9097:9097
```

```shell
$ kubectl apply -f pipeline/
```

```shell
$ kubectl apply -f pipeline-run/
```

```shell
$ kubectl apply -f pipeline-trigger/
```

```shell
$ kubectl port-forward service/el-hello-listener 8080

$ curl -v -H 'content-Type: application/json' -d '{"username": "Tekton"}' http://localhost:8080
```