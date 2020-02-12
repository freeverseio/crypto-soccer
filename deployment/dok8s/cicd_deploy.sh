#! /bin/bash
# you must set the following environmental variables (take values from your kube/config):
#KUBERNETES_CLUSTER_CERTIFICATE
#KUBERNETES_SERVER
#KUBERNETES_TOKEN

set -e

MY_DIR=`dirname "$0"`
MY_DIR=`cd "$MY_DIR" ; pwd`

NAMESPACE=freeverse # TODO: pass as argument so we can use the same script to deploy to testing namespace or production namespace
#KUBERNETES_TOKEN=$(kubectl get secret -n ${NAMESPACE} $(kubectl get secret -n ${NAMESPACE} | grep cicd-token | awk '{print $1}') -o jsonpath='{.data.token}' | base64 --decode)

#echo "${KUBERNETES_CLUSTER_CERTIFICATE}" | base64 --decode > cert.crt
#echo "${KUBERNETES_CLUSTER_CERTIFICATE}" | base64 -d > cert.crt
#KUBECTL="kubectl --kubeconfig=/dev/null --server=${KUBERNETES_SERVER} --certificate-authority=cert.crt --token=${KUBERNETES_TOKEN}"

# example from https://www.digitalocean.com/community/tutorials/how-to-automate-deployments-to-digitalocean-kubernetes-with-circleci
#envsubst <./kube/do-sample-deployment.yml >./kube/do-sample-deployment.yml.out
#kubectl \
#  --kubeconfig=/dev/null \
#  --server=${KUBERNETES_SERVER} \
#  --certificate-authority=cert.crt \
#  --token=${KUBERNETES_TOKEN} \
#  get pods -n ${NAMESPACE}

echo -- deploying ingress,configmap,universedb,universeapi,ipsnode,marketdb,trader,notary
kubectl apply -f ${MY_DIR}/ingress.yaml      -n ${NAMESPACE}
kubectl apply -f ${MY_DIR}/configmap.yaml    -n ${NAMESPACE}
kubectl apply -f ${MY_DIR}/universedb.yaml   -n ${NAMESPACE}
kubectl apply -f ${MY_DIR}/universeapi.yaml  -n ${NAMESPACE}
kubectl apply -f ${MY_DIR}/ipfsnode.yaml     -n ${NAMESPACE}
kubectl apply -f ${MY_DIR}/marketdb.yaml     -n ${NAMESPACE}
kubectl apply -f ${MY_DIR}/trader.yaml       -n ${NAMESPACE}
kubectl apply -f ${MY_DIR}/notary.yaml       -n ${NAMESPACE}

echo -- waiting for pods to be ready...
kubectl wait --for=condition=available --timeout=600s deployment/universedb -n ${NAMESPACE}
UNIVERSEDB_POD=$(kubectl get pod -l app=universedb -n ${NAMESPACE} -o jsonpath="{.items[0].metadata.name}")

kubectl wait --for=condition=available --timeout=600s deployment/universeapi -n ${NAMESPACE}
UNIVERSEAPI_POD=$(kubectl get pod -l app=universeapi -n ${NAMESPACE} -o jsonpath="{.items[0].metadata.name}")

kubectl wait --for=condition=available --timeout=600s deployment/ipfsnode -n ${NAMESPACE}
IPFSNODE_POD=$(kubectl get pod -l app=ipfsnode -n ${NAMESPACE} -o jsonpath="{.items[0].metadata.name}")

kubectl wait --for=condition=available --timeout=600s deployment/trader -n ${NAMESPACE}
TRADER_POD=$(kubectl get pod -l app=trader -n ${NAMESPACE} -o jsonpath="{.items[0].metadata.name}")

kubectl wait --for=condition=Ready --timeout=600s pod/${UNIVERSEDB_POD}  -n ${NAMESPACE}
kubectl wait --for=condition=Ready --timeout=600s pod/${UNIVERSEAPI_POD} -n ${NAMESPACE}
kubectl wait --for=condition=Ready --timeout=600s pod/${IPFSNODE_POD}    -n ${NAMESPACE}
kubectl wait --for=condition=Ready --timeout=600s pod/${TRADER_POD}      -n ${NAMESPACE}


echo -- deploying relayactions, synchronizer and horizon
kubectl apply -f ${MY_DIR}/relayactions.yaml -n ${NAMESPACE}
kubectl apply -f ${MY_DIR}/synchronizer.yaml -n ${NAMESPACE}
kubectl apply -f ${MY_DIR}/horizon.yaml      -n ${NAMESPACE}

kubectl get pods -n ${NAMESPACE}

#rm cert.crt
