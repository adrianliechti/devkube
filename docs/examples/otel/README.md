# OpenTelemetry Demo 

## Demo

### Install Application

```shell
# Add OpenTelemetry Helm Charts repo
helm repo add open-telemetry https://open-telemetry.github.io/opentelemetry-helm-charts
helm repo update

# Install Demo
helm upgrade --install opentelemetry-demo open-telemetry/opentelemetry-demo --version 0.6.0 --create-namespace --namespace opentelemetry-demo --values values.yaml

# Wait until all pods ready
kubectl get pods --namespace opentelemetry-demo --watch
```

### Access Demo App

```shell
kubectl port-forward service/opentelemetry-demo-frontend 8080:8080 --namespace opentelemetry-demo
```

### Watch Service Graph

```shell
devkube grafana
```

- navigate to "Explore" on the navigation bar left
- choose "Tempo" as data source on the top
- select "Service Graph" view
- click the "Run query" button