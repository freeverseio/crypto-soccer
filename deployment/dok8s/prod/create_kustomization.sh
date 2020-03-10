#!/bin/sh
set -e

MY_DIR=`dirname "$0"`
MY_DIR=`cd "$MY_DIR" ; pwd`

BASE_DIR=../base
NAMESPACE=freeverse
APP_NAME='cryptosoccer'
APP_VERSION='1.0.0'
TAG="1.0.0-alpha"

# create kustomization.yaml
kustomize create
kustomize edit add label 'app.kubernetes.io/part-of':${APP_NAME},'app.kubernetes.io/version':${APP_VERSION}

kustomize edit add base ${BASE_DIR}
#kustomize edit set namespace ${NAMESPACE}

# set image tags 
kustomize edit set image freeverseio/horizon:${TAG}
kustomize edit set image freeverseio/market.db:${TAG}
kustomize edit set image freeverseio/market.notary:${TAG}
kustomize edit set image freeverseio/relay.actions:${TAG}
kustomize edit set image freeverseio/synchronizer:${TAG}
kustomize edit set image freeverseio/market.trader:${TAG}
kustomize edit set image freeverseio/universe.api:${TAG}
kustomize edit set image freeverseio/universe.db:${TAG}
kustomize edit set image freeverseio/authproxy:${TAG}

# change to n replicas
# kustomize edit set replicas horizon=1
# kustomize edit set replicas universeapi=1
# kustomize edit set replicas trader=1

# patching
kustomize edit add patch configmap.yaml

# build application to be deployed
kustomize build ${MY_DIR} -o ${MY_DIR}/app.yaml

# or alternative apply directly to cluster
# kubectl apply -k .
