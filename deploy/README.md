1. Custom the ConfigMap `pod-policy` at `kustomization.yaml` as below:

```yaml
configMapGenerator:
  - name: pod-policy
    literals:
      - image.registry=docker.io
      - image.org=qqbuby
```

2. Apply manifests with `kubectl apply -k .`

3. Use `kubectl label ns [namespace] pod-policy.kuberos.io=true`.
