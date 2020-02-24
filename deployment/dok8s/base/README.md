- to build and deploy
```bash
kustomaize build . > output.yaml
kubectl apply -f output.yaml
```

- or to deploy directly
```bash
kubectl apply -k .
```
