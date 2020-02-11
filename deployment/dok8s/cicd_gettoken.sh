#! /bin/bash
NAMESPACE=$1
TOKEN=$(kubectl get secret -n ${NAMESPACE} $(kubectl get secret -n ${NAMESPACE} | grep cicd-token | awk '{print $1}') -o jsonpath='{.data.token}' | base64 --decode)
echo $TOKEN
