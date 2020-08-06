- setup: contains initial configuration files
- base: contains application specific yaml files
- dev: contains customization from base


# Cluster creation

```bash
helm install nginx-ingress stable/nginx-ingress --set controller.publishService.enabled=true
```

# Securing the Ingress Using Cert-Manager
follow https://www.digitalocean.com/community/tutorials/how-to-set-up-an-nginx-ingress-on-digitalocean-kubernetes-using-helm
```bash
kubectl create namespace cert-manager
kubectl apply -f https://raw.githubusercontent.com/jetstack/cert-manager/release-0.11/deploy/manifests/00-crds.yaml
helm repo add jetstack https://charts.jetstack.io
helm install cert-manager --version v0.11.0 --namespace cert-manager jetstack/cert-manager
kubectl apply -f cert-issuer.yaml
```

# CircleCI
Create a service account to be used in circleci
```bash
kubectl -n freeverse -f ./cicd/cicd-serviceaccount.yaml apply
```

Get the tocken of cicd
```bash
kubectl -n freeverse get secrets | grep cicd
```

# Secrets
market, staker, relay needs credentials to sign transactions. Let's create a secret where to store them
```bash
kubectl -n freeverse create secret generic blockchain-accounts --from-literal=relay=<private_key> --from-literal=market=<private_key> --from-literal=staker=<private_key>

kubectl -n freeverse create secret generic google-iap-key --from-file=./key.json
```

## Repeling pods from special hardware nodes

### Requirements

The cluster should have a pool of nodes with special hardware, in this case, nodes with more than 2 CPUs.

The cluster admin should taint these nodes so the kube API can differentiate them.

This command will give you the nodes in your cluster with more than 2 cpus: 
```
kubectl get nodes -o go-template --template='{{range .items}}{{if gt .status.capacity.cpu "2"}}{{printf "%v\n" .metadata.name}}{{end}}{{end}}'
```

This command will taint the nodes with more than 2 cpus:
 ```
 kubectl get nodes -o go-template --template='{{range .items}}{{if ge .status.capacity.cpu "2"}}{{printf "%v\n" .metadata.name}}{{end}}{{end}}' | xargs -I '{}' -n 1 $(echo "kubectl taint nodes {} cpu=dedicated:NoSchedule")
 ```

 Once your nodes are tainted, the kube scheduler will only schedule pods in this nodes if the pods have a corresponding toleration that matches the tainted nodes.

 A custom admission controller should be created in order to mutate the pods of the namespace freeverse adding the toleration. So the expected result will be that only the pods of the namespace freeverse will go to the >2 cpu node.