# How to deploy to k8s
```bash
./create_customization.sh
kustomize build -o output.yaml
kubectl apply -f output.yaml -n <namespace>
```

# inject linkerd
```bash
cat output.yaml | linkerd inject - | kubectl apply -f - -n <namespace>
linkerd dashboard &
```
