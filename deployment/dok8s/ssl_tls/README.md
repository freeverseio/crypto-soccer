# Securing the Ingress Using Cert-Manager
follow https://www.digitalocean.com/community/tutorials/how-to-set-up-an-nginx-ingress-on-digitalocean-kubernetes-using-helm
```bash
kubectl create namespace cert-manager
kubectl apply -f https://raw.githubusercontent.com/jetstack/cert-manager/release-0.11/deploy/manifests/00-crds.yaml
helm repo add jetstack https://charts.jetstack.io
helm install --name cert-manager --version v0.11.0 --namespace cert-manager jetstack/cert-manager
kubectl apply -f cert-issuer.yaml
```
