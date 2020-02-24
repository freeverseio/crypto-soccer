#! /bin/bash

set -e

MY_DIR=`dirname "$0"`
MY_DIR=`cd "$MY_DIR" ; pwd`
DOCKER_REGISTRY_SERVER=docker.io
NAMESPACE=freeverse # TODO: pass as argument so we can use the same script to deploy to testing namespace or production namespace

namespace_and_secret()
{
    NS=`kubectl get pods --all-namespaces | awk '{print $1;}' | uniq | grep ${NAMESPACE}`

    if [ $NS != $NAMESPACE ]; then kubectl create ns ${NAMESPACE}
    else echo "namespace $NAMESPACE already exist. Omitting creation."
    fi
    SECRET=`kubectl get secret -n ${NAMESPACE} | grep docker-secret | awk '{print $1}'`
    if [ $SECRET != 'docker-secret' ]; then kubectl create secret docker-registry docker-secret --docker-server=${DOCKER_REGISTRY_SERVER} --docker-username=${DROPLET_DOCKER_ID} --docker-password=${DROPLET_DOCKER_PASSWD} --docker-email=freeversedigitalocean@freeverse.io -n ${NAMESPACE}
    else echo "docker secret already exist. Omitting creation."
    fi
}

clean()
{
    echo -- removing deployments
    kubectl delete -k ${MY_DIR}/dev/ingress      -n ${NAMESPACE}
    kubectl delete -k ${MY_DIR}/dev/configmap    -n ${NAMESPACE}
    # kubectl delete -k ${MY_DIR}/universedb   -n ${NAMESPACE}
    kubectl delete -k ${MY_DIR}/dev/universeapi  -n ${NAMESPACE}
    kubectl delete -k ${MY_DIR}/dev/marketdb     -n ${NAMESPACE}
    kubectl delete -k ${MY_DIR}/dev/trader       -n ${NAMESPACE}
    kubectl delete -k ${MY_DIR}/dev/notary       -n ${NAMESPACE}
    kubectl delete -k ${MY_DIR}/dev/synchronizer -n ${NAMESPACE}
    kubectl delete -k ${MY_DIR}/dev/relayactions -n ${NAMESPACE}
    kubectl delete -k ${MY_DIR}/dev/authproxy    -n ${NAMESPACE}
}

deploy()
{
    echo -- deploying
    kubectl apply -k ${MY_DIR}/dev            -n ${NAMESPACE}
    kubectl apply -f ${MY_DIR}/authproxy.yaml -n ${NAMESPACE}

    #kubectl apply -f ${MY_DIR}/ingress.yaml      -n ${NAMESPACE}
    #kubectl apply -f ${MY_DIR}/configmap.yaml    -n ${NAMESPACE}
    ## kubectl apply -f ${MY_DIR}/universedb.yaml   -n ${NAMESPACE}
    #kubectl apply -f ${MY_DIR}/universeapi.yaml  -n ${NAMESPACE}
    #kubectl apply -f ${MY_DIR}/ipfsnode.yaml     -n ${NAMESPACE}
    #kubectl apply -f ${MY_DIR}/marketdb.yaml     -n ${NAMESPACE}
    #kubectl apply -f ${MY_DIR}/trader.yaml       -n ${NAMESPACE}
    #kubectl apply -f ${MY_DIR}/notary.yaml       -n ${NAMESPACE}
    #kubectl apply -f ${MY_DIR}/relayactions.yaml -n ${NAMESPACE}
    #kubectl apply -f ${MY_DIR}/synchronizer.yaml -n ${NAMESPACE}
    #kubectl apply -f ${MY_DIR}/horizon.yaml      -n ${NAMESPACE}
}

clean
#namespace_and_secret # circle-ci does not have access to do this
deploy
kubectl get pods -n ${NAMESPACE}
