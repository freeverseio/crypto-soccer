# installing nginx ingress controller
- install tiller serviceaccount
    ```bash
    kubectl -n kube-system create serviceaccount tiller
    ```
- Next, bind the tiller serviceaccount to the cluster-admin role:
    ```bash
    kubectl create clusterrolebinding tiller --clusterrole cluster-admin --serviceaccount=kube-system:tiller
    ```
- use helm to install tiller
    ```bash
    helm init --service-account tiller
    ```
    TODO:
    Please note: by default, Tiller is deployed with an insecure 'allow unauthenticated users' policy.
    To prevent this, run `helm init` with the --tiller-tls-verify flag.
    For more information on securing your installation see: https://docs.helm.sh/using_helm/#securing-your-helm-installation

- install ngingx ingress controller
    ```bash
    helm install stable/nginx-ingress --name nginx-ingress --set controller.publishService.enabled=true
    ```

# Securing the Ingress Using Cert-Manager
follow https://www.digitalocean.com/community/tutorials/how-to-set-up-an-nginx-ingress-on-digitalocean-kubernetes-using-helm
```bash
kubectl create namespace cert-manager
kubectl apply -f https://raw.githubusercontent.com/jetstack/cert-manager/release-0.11/deploy/manifests/00-crds.yaml
helm repo add jetstack https://charts.jetstack.io
helm install --name cert-manager --version v0.11.0 --namespace cert-manager jetstack/cert-manager
kubectl apply -f cert-issuer.yaml
```
