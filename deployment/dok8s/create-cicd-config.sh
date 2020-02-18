#!/bin/sh
set -E
server=https://365763e3-86ce-4576-8628-d6eebca12dab.k8s.ondigitalocean.com
cluster_name=do-ams3-k8s-1-16-6-do-0-ams3-cryptosoccer-test2
username=cicd
namespace=freeverse

# the name of the secret containing the service account token goes here
secret_name=cicd-token-sbflk


ca=$(kubectl get secret/$secret_name -n ${namespace} -o jsonpath='{.data.ca\.crt}')
token=$(kubectl get secret/$secret_name -n ${namespace} -o jsonpath='{.data.token}' | base64 --decode)

echo "apiVersion: v1
kind: Config
clusters:
- name: ${cluster_name}
  cluster:
    certificate-authority-data: ${ca}
    server: ${server}
contexts:
- name: cicd-context
  context:
    cluster: ${cluster_name}
    namespace: ${namespace}
    user: ${username}
current-context: cicd-context
users:
- name: ${username}
  user:
    token: ${token}
" > cicd.kubeconfig.yaml
