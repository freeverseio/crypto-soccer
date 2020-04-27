#!/bin/sh
set -E
output_filename=cicd.kubeconfig.yaml
server=https://31b1a095-1ef3-4940-aad6-4df9ee9ae947.k8s.ondigitalocean.com
cluster_name=do-ams3-k8s-development-goalrevolution
username=cicd
namespace=freeverse

# the name of the secret containing the service account token goes here
secret_name=cicd-token-2gg8r


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
" > ${output_filename}

cat ${output_filename} | base64 -w 0
# cat ${output_filename} | base64
