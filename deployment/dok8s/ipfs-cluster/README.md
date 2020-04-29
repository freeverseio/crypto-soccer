# description
This directory contains kustomizations for creating an ipfs cluster (see cluster directory) and an external peer (see single_peer directory)
I have seen a lot of instability when deployed in kukernetes without a service mesh, so for this reason I recommend that you install linkerd in your cluster before deploying (https://linkerd.io/2/faq/).
Also before starting you need to create ids and private keys for all the components so they are know and can be configured via the env-configmap.yaml and secret.yaml (see overlays).


# how to install linkerd
```bash
linkerd install | kubectl apply -f -
```
# cluster secret
The cluster and all its peers need to know/share the same secret. In order to generate it do:
1. generate some random data
```bash
secret=$(od  -vN 32 -An -tx1 /dev/urandom | tr -d ' \n')
```
2. encode it with base64
```bash
encoded_secret=`echo -n $secret | base64`
```
3. copy the encoded secret to cluster-secret field in secret.yam (see overlays)

# IPFS IDs and private keys
In order to generate the keys that ipfs needs you need to use ipfs-key (https://github.com/whyrusleeping/ipfs-key) and encode it with base64.

## bootstrapper ID and private key
1. generate id and key with
```bash
ipfs-key | base64
```
2. copy the output private key and encode it again with base64 and then copy it to field bootstrap-peer-priv-key in secret.yaml. It needs to be re-encoded because that's how k8s secret works.
```bash
<output_private_key> | base64
```
3. copy the ID into bootstrap-peer-id in env-configmap.yaml


## peer ID and private key
Similarly to cluster but select key type Ed25519.
1. generate id and key with
```bash
ipfs-key -type Ed25519 | base64
```
2. copy the output private key and encode it again with base64 and then copy it to field peer-priv-key in secret.yaml. It needs to be re-encoded because that's how k8s secret works.
```bash
<output_private_key> | base64
```
3. copy the ID into peer-id in env-configmap.yaml

# cluster
## description
This directory contains the configuration for deploying the main cluster. Currently deploys 1 bootstrapper and 2 peers in the same namespace via an statefulset. The names of the deployed pods are ipfs-cluster-0 (bootstrapper), ipfs-cluster-1 (peer 1) and ipfs-cluster-2 (peer 2). Both of these peers are considered trusted peers when creating the cluster. Note that a trusted peer is such as it has permission to change pinset. In principle these two peers should never pin anything as pins will be added by an external trusted peer: the Goal Revolution game. The reason for creating these 2 internal peers is due to the need for availability of trusted peers for upcoming stakers followers.

## deploy cluster
1. make sure linkerd was installed first
2. modify overlay accordingly to suit your needs
3. deploy to kubernetes as:
```bash
cd cluster/overlays/<your_customization>
kustomize build -o output.yaml
kubectl apply -f output.yaml -n <namespace>
cat output.yaml | linkerd inject - | kubectl apply -f - -n <namespace>
```
4. watch via normal kubernetes dashboard or use linkerd dashboard
```bash
linkerd dashboard &
```

# single_peer
## description
This directory contains the configuration for deploying a peer external to the above ipfs cluster. It deploys ipfs and ipfs-cluster via a statefulset and needs the peer addresses of the the above trusted peers.
The idea is that the Goal Revolution game will deploy this configuration in its kubernetes cluster.
## deploy peer
1. make sure linkerd was installed first
2. modify overlay accordingly to suit your needs
3. deploy to kubernetes as:
```bash
cd single_peer/overlays/<your_customization>
kustomize build -o output.yaml
kubectl apply -f output.yaml -n <namespace>
cat output.yaml | linkerd inject - | kubectl apply -f - -n <namespace>
```
4. watch via normal kubernetes dashboard or use linkerd dashboard
```bash
linkerd dashboard &
