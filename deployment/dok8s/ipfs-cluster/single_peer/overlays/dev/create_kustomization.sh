#!/bin/sh
set -e

MY_DIR=`dirname "$0"`
MY_DIR=`cd "$MY_DIR" ; pwd`

BASE_DIR=../../base
IPFS_TAG="v0.4.23"
IPFS_CLUSTER_TAG="v0.12.1"
#NAMESPACE=ipfs

# create kustomization.yaml
kustomize create
kustomize edit add base ${BASE_DIR}

#kustomize edit set namespace ${NAMESPACE}

# set image tags
kustomize edit set image ipfs/ipfs-cluster:${IPFS_CLUSTER_TAG}
kustomize edit set image ipfs/go-ipfs:${IPFS_TAG}

# change to n replicas
kustomize edit set replicas ipfs-cluster=3

# patching
kustomize edit add patch env-configmap.yaml
kustomize edit add patch secret.yaml
