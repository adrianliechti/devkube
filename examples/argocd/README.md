# Argo CD

## Install Argo CD

```shell
kubectl create namespace argocd

kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/examples/k8s-rbac/argocd-server-applications/argocd-server-rbac-clusterrole.yaml
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/examples/k8s-rbac/argocd-server-applications/argocd-server-rbac-clusterrolebinding.yaml

kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/examples/k8s-rbac/argocd-server-applications/argocd-notifications-controller-rbac-clusterrole.yaml
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/examples/k8s-rbac/argocd-server-applications/argocd-notifications-controller-rbac-clusterrolebinding.yaml

kubectl apply -n argocd -f argocd-cmd-params.yaml
kubectl apply -n argocd -f argocd-project-default.yaml

kubectl rollout restart -n argocd deployment argocd-server
kubectl rollout restart -n argocd statefulset argocd-application-controller
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