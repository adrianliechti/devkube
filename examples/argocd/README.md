# Argo CD

## Install Argo CD

```shell
$ devkube install argocd
```

### Open Server UI

```shell
$ kubectl port-forward -n argocd service/argocd-server 8080:80
```

open [http://localhost:8080](http://localhost:8080) in browser

### Sample App

```shell
$ kubectl apply -n argocd -f guestbook/
```