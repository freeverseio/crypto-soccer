# How to deploy to k8s
```bash
./create_customization.sh
kustomize build -o output.yaml
kubectl apply -f output.yaml -n <namespace>
```
