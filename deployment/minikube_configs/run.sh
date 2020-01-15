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
    echo "    start            start minikube"
    echo "    secret           create namespace and docker credentials"
    echo "    ethereum         start ethereum node and deploy contracts"
    echo "    deploycontracts  deploy contracts on ethereum node"
    echo "    freeverse        deploy all freeverse pods"
    echo "    clean            remove everything"
    echo
    echo "The usual workflow would be:"
    echo "1. ./run.sh start"
    echo "2. ./run.sh secret"
    echo "3. ./run.sh ethereum [optional]"
    echo "4. ./run.sh deploycontracts [optional]"
    echo "5. ./run.sh freeverse"
}

start_minikube()
{
    minikube start --vm-driver=hyperkit
    mkdir -p /tmp/minikube-storage # create a directory to share with minikube
    minikube mount /tmp/minikube-storage:/data & # mount into minikube as /data
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
    kubectl apply -f ${MY_DIR}/ethereum.yaml -n ${NAMESPACE}
    echo -- waiting until ethereum POD is running
    kubectl wait --for=condition=available --timeout=600s deployment/ethereum -n ${NAMESPACE}
    POD=$(kubectl get pod -l app=ethereum -n ${NAMESPACE} -o jsonpath="{.items[0].metadata.name}")
    kubectl wait --for=condition=Ready --timeout=600s pod/${POD} -n ${NAMESPACE}
}

deploycontracts()
{
    eval $(minikube docker-env)
    echo In order to deploy contracts we need credentials to freeverse git repo.
    read -p "Enter git token user name: "  username
    read -s -p "Enter git token password: " password
    cd ${MY_DIR}/../../go
    docker build -f Dockerfile.contractdeployment -t contractdeployment:0.0.1 --build-arg GIT_TOKEN_USR=$username --build-arg GIT_TOKEN_PWD=$password .
    cd ${MY_DIR}
    echo -- deploying contracts
    kubectl apply -f ${MY_DIR}/contractdeployment.yaml -n ${NAMESPACE}
    # wait for job to complete
    kubectl wait --for=condition=complete --timeout=600s job/contractdeployment -n ${NAMESPACE}
    # print contract addresses
    echo -- copy the following addresses into configmap.yaml
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
    elif [ $arg == 'deploycontracts' ];  then deploycontracts
    elif [ $arg == 'freeverse' ]; then freeverse
    elif [ $arg == 'clean' ];     then clean
    else help
    fi
    shift 1
done
