1. Generating self-signed certificate

```sh
../kuberos self-signed-cert \
    -u kuberos.default.svc \
    -n kuberos.default.svc \
    -o manifests/tls/cert.crt \
    -k manifests/tls/cert.key
```

2. Set `caBundle` in `admissionregistration.yaml` with `manifests/tls/cert.crt`

```sh
cat manifests/tls/cert.crt | base64 | tr -d '\n'
```

3. Apply manifests with `kubectl apply -k .`
