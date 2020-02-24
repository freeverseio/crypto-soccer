#!/bin/sh
set -e

MY_DIR=`dirname "$0"`
MY_DIR=`cd "$MY_DIR" ; pwd`

BASE_DIR=../base
NAMESPACE=freeverse
APP_NAME='cryptosoccer'
APP_VERSION='1.0.0'
FREEVERSE_TAG="dev"
IPFS_TAG="v0.4.23"

# create kustomization.yaml
kustomize create
kustomize edit add label 'app.kubernetes.io/part-of':${APP_NAME},'app.kubernetes.io/version':${APP_VERSION}

kustomize edit add base ${BASE_DIR}
kustomize edit set namespace ${NAMESPACE}

# set image tags 
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
