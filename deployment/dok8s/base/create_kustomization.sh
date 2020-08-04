#!/bin/sh
kustomize create
kustomize edit add resource configmap.yaml
kustomize edit add resource ingress.yaml
kustomize edit add base authproxy
#kustomize edit add base helloworld
kustomize edit add base horizon
kustomize edit add base marketdb
kustomize edit add base notary
kustomize edit add base relayactions
kustomize edit add base synchronizer
kustomize edit add base trader
kustomize edit add base universeapi
kustomize edit add base universedb
kustomize edit add base dashboard
kustomize edit add base webphoenix
kustomize edit add base oauth2-proxy
# kustomize edit add base gamedb
# kustomize edit add base gameapi
# kustomize edit add base gamelayer

