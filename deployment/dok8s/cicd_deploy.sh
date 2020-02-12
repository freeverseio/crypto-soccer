#! /bin/bash

set -e

MY_DIR=`dirname "$0"`
MY_DIR=`cd "$MY_DIR" ; pwd`

NAMESPACE=freeverse # TODO: pass as argument so we can use the same script to deploy to testing namespace or production namespace

clean()
{
    echo -- removing deployments
    kubectl delete -f ${MY_DIR}/ingress.yaml      -n ${NAMESPACE}
    kubectl delete -f ${MY_DIR}/configmap.yaml    -n ${NAMESPACE}
    kubectl delete -f ${MY_DIR}/universedb.yaml   -n ${NAMESPACE}
    kubectl delete -f ${MY_DIR}/universeapi.yaml  -n ${NAMESPACE}
    kubectl delete -f ${MY_DIR}/ipfsnode.yaml     -n ${NAMESPACE}
    kubectl delete -f ${MY_DIR}/marketdb.yaml     -n ${NAMESPACE}
    kubectl delete -f ${MY_DIR}/trader.yaml       -n ${NAMESPACE}
    kubectl delete -f ${MY_DIR}/notary.yaml       -n ${NAMESPACE}
}

deploy()
{
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
}

clean
deploy
