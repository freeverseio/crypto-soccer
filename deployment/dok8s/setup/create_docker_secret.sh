#!/bin/sh

NAMESPACE=$1
DOCKER_REGISTRY_SERVER=docker.io
DOCKER_USER=freeverseiodev
DOCKER_PASSWORD=MWx4/CgC5/3hY+/8sINqJpMSvalbG1xwUbeZBpPU0qc=

# kubectl create ns ${NAMESPACE}
kubectl create secret docker-registry docker-secret --docker-server=${DOCKER_REGISTRY_SERVER} --docker-username=${DOCKER_USER} --docker-password=${DOCKER_PASSWORD} --docker-email=asiniscalchi@freeverse.io -n ${NAMESPACE}
