# Ingress Demo

## Demo

### Install Application

```bash
kubectl apply -f manifest.yaml
```

### Access Demo App

```bash
# trust Root CA
devkube trust

devkube ingress
```

```bash
open https://demo.example.org
```
