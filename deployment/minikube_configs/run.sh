#!/bin/sh

MY_DIR=`dirname "$0"`
argerr=0

NAMESPACE=foo
DOCKER_REGISTRY_SERVER=docker.io
DOCKER_USER=freeversedigitalocean
DOCKER_PASSWORD=fdggreTRGDSBw45ergseth4hDGHD

help()
{
    echo usage:
    echo "    ./run.sh <option>"
    echo
    echo "Option list:"
    echo "    start            start minikube"
    echo "    secret           create namespace and docker credentials"
    echo "    ethereum         start ethereum node and deploy contracts"
    echo "    freeverse        deploy all freeverse pods"
    echo "    clean            remove everything"
    echo
    echo "The usual workflow would be:"
    echo "1. ./run.sh start"
    echo "2. ./run.sh secret"
    echo "3. ./run.sh ethereum [optional, but make sure you have built contractdeployment docker image first]"
    echo "4. ./run.sh freeverse"
}

start_minikube()
{
    minikube start --vm-driver=hyperkit
    eval $(minikube docker-env)
}

clean()
{
    minikube stop && minikube delete
}

namespace_and_secret()
{
    kubectl create ns ${NAMESPACE}
    kubectl create secret docker-registry docker-secret --docker-server=${DOCKER_REGISTRY_SERVER} --docker-username=${DOCKER_USER} --docker-password=${DOCKER_PASSWORD} --docker-email=freeversedigitalocean@freeverse.io -n ${NAMESPACE}
}

ethereum()
{
    echo "Running ethereum"
    kubectl apply -f ${MY_DIR}/ethereum.yaml -n ${NAMESPACE}
    echo waiting until ethereum POD is running
    kubectl wait --for=condition=available --timeout=600s deployment/ethereum -n ${NAMESPACE}
    POD=$(kubectl get pod -l app=ethereum -n ${NAMESPACE} -o jsonpath="{.items[0].metadata.name}")
    kubectl wait --for=condition=Ready --timeout=600s pod/${POD} -n ${NAMESPACE}
    # deploy contracts
    echo "Deploying contracts"
    kubectl apply -f ${MY_DIR}/contractdeployment.yaml -n ${NAMESPACE}
    # wait for job to complete
    kubectl wait --for=condition=complete --timeout=600s job/contractdeployment -n ${NAMESPACE}
    # print contract addresses
    echo copy the following addresses into configmap.yaml
    kubectl logs job/contractdeployment -n ${NAMESPACE}
}

freeverse()
{
    # deploy rest of contracts
    kubectl apply -f ${MY_DIR}/configmap.yaml    -n ${NAMESPACE}
    kubectl apply -f ${MY_DIR}/universedb.yaml   -n ${NAMESPACE}
    kubectl apply -f ${MY_DIR}/universeapi.yaml  -n ${NAMESPACE}
    kubectl apply -f ${MY_DIR}/synchronizer.yaml -n ${NAMESPACE}
    kubectl apply -f ${MY_DIR}/ipfsnode.yaml     -n ${NAMESPACE}
    kubectl apply -f ${MY_DIR}/relaydb.yaml      -n ${NAMESPACE}
    #kubectl apply -f ${MY_DIR}/relayapi.yaml -n ${NAMESPACE}
}

until [ $# -eq 0 ]
do
    arg=$1
    if   [ $arg == 'start' ];     then start_minikube
    elif [ $arg == 'secret' ];    then namespace_and_secret
    elif [ $arg == 'ethereum' ];  then ethereum
    elif [ $arg == 'freeverse' ]; then freeverse
    elif [ $arg == 'clean' ];     then clean
    else help
    fi
    shift 1
done
