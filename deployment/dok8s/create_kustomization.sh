#!/bin/sh


# create kustomization.yaml
kustomize create

# add resources
kustomize edit add resource cert-issuer.yaml
kustomize edit add resource configmap.yaml
kustomize edit add resource helloworld.yaml
kustomize edit add resource horizon.yaml
kustomize edit add resource ingress.yaml
kustomize edit add resource ipfsnode.yaml
kustomize edit add resource marketdb.yaml
kustomize edit add resource notary.yaml
kustomize edit add resource relayactions.yaml
kustomize edit add resource synchronizer.yaml
kustomize edit add resource trader.yaml
kustomize edit add resource universeapi.yaml
kustomize edit add resource universedb.yaml
kustomize edit add resource xdai.yaml
