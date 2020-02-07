#!/bin/sh

MY_DIR=`dirname "$0"`
MY_DIR=`cd "$MY_DIR" ; pwd`

NAMESPACE=freeverse
DOCKER_REGISTRY_SERVER=docker.io
DOCKER_USER=freeversedigitalocean
DOCKER_PASSWORD=fdggreTRGDSBw45ergseth4hDGHD

help()
{
    echo usage:
    echo "    ./run.sh <option>"
    echo
    echo "Option list:"
    echo "    secret           create namespace and docker credentials"
    echo "    ethereum         start ethereum node"
    echo "    freeverse        deploy all freeverse pods"
    echo "    clean            remove everything"
    echo
    echo "The usual workflow would be:"
    echo "1. ./run.sh secret"
    echo "2. ./run.sh ethereum"
    echo "3. ./run.sh freeverse"
}

namespace_and_secret()
{
    kubectl create ns ${NAMESPACE}
    kubectl create secret docker-registry docker-secret --docker-server=${DOCKER_REGISTRY_SERVER} --docker-username=${DOCKER_USER} --docker-password=${DOCKER_PASSWORD} --docker-email=freeversedigitalocean@freeverse.io -n ${NAMESPACE}
    kubectl apply -f ${MY_DIR}/helloworld.yaml -n ${NAMESPACE}
}

ethereum()
{
    kubectl apply -f ${MY_DIR}/ethereum.yaml -n ${NAMESPACE}
    echo -- waiting until ethereum POD is running
    kubectl wait --for=condition=available --timeout=600s deployment/ethereum -n ${NAMESPACE}
    POD=$(kubectl get pod -l app=ethereum -n ${NAMESPACE} -o jsonpath="{.items[0].metadata.name}")
    kubectl wait --for=condition=Ready --timeout=600s pod/${POD} -n ${NAMESPACE}
}

freeverse()
{
    # deploy contracts that have no dependencies
    echo deploying ingress,configmap,universedb,universeapi,ipsnode,marketdb,trader,notary
    kubectl apply -f ${MY_DIR}/ingress.yaml      -n ${NAMESPACE}
    kubectl apply -f ${MY_DIR}/configmap.yaml    -n ${NAMESPACE}
    kubectl apply -f ${MY_DIR}/universedb.yaml   -n ${NAMESPACE}
    kubectl apply -f ${MY_DIR}/universeapi.yaml  -n ${NAMESPACE}
    kubectl apply -f ${MY_DIR}/ipfsnode.yaml     -n ${NAMESPACE}
    kubectl apply -f ${MY_DIR}/marketdb.yaml     -n ${NAMESPACE}
    kubectl apply -f ${MY_DIR}/trader.yaml       -n ${NAMESPACE}
    kubectl apply -f ${MY_DIR}/notary.yaml       -n ${NAMESPACE}

    echo waiting for pods to be ready...
    kubectl wait --for=condition=available --timeout=600s deployment/universedb -n ${NAMESPACE}
    UNIVERSEDB_POD=$(kubectl get pod -l app=universedb -n ${NAMESPACE} -o jsonpath="{.items[0].metadata.name}")

    kubectl wait --for=condition=available --timeout=600s deployment/universeapi -n ${NAMESPACE}
    UNIVERSEAPI_POD=$(kubectl get pod -l app=universeapi -n ${NAMESPACE} -o jsonpath="{.items[0].metadata.name}")

    kubectl wait --for=condition=available --timeout=600s deployment/ipfsnode -n ${NAMESPACE}
    IPFSNODE_POD=$(kubectl get pod -l app=ipfsnode -n ${NAMESPACE} -o jsonpath="{.items[0].metadata.name}")

    kubectl wait --for=condition=available --timeout=600s deployment/trader -n ${NAMESPACE}
    TRADER_POD=$(kubectl get pod -l app=trader -n ${NAMESPACE} -o jsonpath="{.items[0].metadata.name}")

    kubectl wait --for=condition=Ready --timeout=600s pod/${UNIVERSEDB_POD} -n ${NAMESPACE}
    kubectl wait --for=condition=Ready --timeout=600s pod/${UNIVERSEAPI_POD} -n ${NAMESPACE}
    kubectl wait --for=condition=Ready --timeout=600s pod/${IPFSNODE_POD} -n ${NAMESPACE}
    kubectl wait --for=condition=Ready --timeout=600s pod/${TRADER_POD} -n ${NAMESPACE}

    echo deploying relayactions, synchronizer and horizon 
    kubectl apply -f ${MY_DIR}/relayactions.yaml -n ${NAMESPACE}
    kubectl apply -f ${MY_DIR}/synchronizer.yaml -n ${NAMESPACE}
    kubectl apply -f ${MY_DIR}/horizon.yaml -n ${NAMESPACE}
}

until [ $# -eq 0 ]
do
    arg=$1
    if   [ $arg == 'secret' ];    then namespace_and_secret
    elif [ $arg == 'ethereum' ];  then ethereum
    elif [ $arg == 'freeverse' ]; then freeverse
    else help
    fi
    shift 1
done
