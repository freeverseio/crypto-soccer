#!/bin/sh
set -e

MY_DIR=`dirname "$0"`
MY_DIR=`cd "$MY_DIR" ; pwd`

BASE_DIR=../../base
NAMESPACE=ipfs
TAG="latest"

# create kustomization.yaml
kustomize create
kustomize edit add base ${BASE_DIR}

#kustomize edit set namespace ${NAMESPACE}

# set image tags
kustomize edit set image ipfs/ipfs-cluster:${TAG}
kustomize edit set image ipfs/go-ipfs:${TAG}

# change to n replicas
kustomize edit set replicas ipfs-cluster=3

# patching
#kustomize edit add patch env-configmap.yaml
