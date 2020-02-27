#! /bin/bash

set -e

MY_DIR=`dirname "$0"`
MY_DIR=`cd "$MY_DIR" ; pwd`

NAMESPACE=freeverse # TODO: pass as argument so we can use the same script to deploy to testing namespace or production namespace

kubectl delete -f ${MY_DIR}/dev/app.yaml -n ${NAMESPACE}
kubectl apply -f ${MY_DIR}/dev/app.yaml -n ${NAMESPACE}
