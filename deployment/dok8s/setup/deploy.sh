#!/bin/sh

MY_DIR=`dirname "$0"`
MY_DIR=`cd "$MY_DIR" ; pwd`
NAMESPACE=freeverse

${MY_DIR}/create_docker_secret.sh ${NAMESPACE}
kubectl apply -f ${MY_DIR}/ipfsnode     -n ${NAMESPACE}
kubectl apply -f ${MY_DIR}/xdai         -n ${NAMESPACE}
kubectl apply -f ${MY_DIR}/ingress.yaml -n ${NAMESPACE}
