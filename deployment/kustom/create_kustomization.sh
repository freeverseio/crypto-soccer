#!/bin/sh
set -e

MY_DIR=`dirname "$0"`
MY_DIR=`cd "$MY_DIR" ; pwd`

BASE_DIR=../dok8s
FREEVERSE_TAG="1.0"
IPFS_TAG="v0.4.23"

# if kustomization.yaml does not exist in base_dir then create it
if [ ! -f "${BASE_DIR}/kustomization.yaml" ]; then (cd ${BASE_DIR} && ./create_kustomization.sh) fi

# create kustomization.yaml
kustomize create

# set new images
kustomize edit add base ${BASE_DIR}
kustomize edit set image ipfs/go-ipfs:${IPFS_TAG}
kustomize edit set image freeverseio/horizon:${FREEVERSE_TAG}
kustomize edit set image freeverseio/market.db:${FREEVERSE_TAG}
kustomize edit set image freeverseio/market.notary:${FREEVERSE_TAG}
kustomize edit set image freeverseio/relay.actions:${FREEVERSE_TAG}
kustomize edit set image freeverseio/synchronizer:${FREEVERSE_TAG}
kustomize edit set image freeverseio/market.trader:${FREEVERSE_TAG}
kustomize edit set image freeverseio/universe.api:${FREEVERSE_TAG}
kustomize edit set image freeverseio/universe.db:${FREEVERSE_TAG}
kustomize edit set image freeverseio/xdai:${FREEVERSE_TAG}

# build application to be deployed
kustomize build .

# or alternative apply directly to cluster
# kubectl apply -k .
