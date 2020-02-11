#! /bin/bash
usage()
{
    echo this script runs a kubectl command on behalf of cicd service account
    echo usage:
    echo "     $0 <server> <service_account_namespace> <kubectl_command>"
    echo
    echo example:
    echo "     ./cicd_kubectl.sh https://12345-12345-abcdefg.k8s.abc.com my_namespace get pods -n my_namespace"
}

set -e

if [ $# -lt 3 ]; then usage
else
    SERVER=$1
    shift 1
    NAMESPACE=$1
    shift 1
    TOKEN=$(kubectl get secret -n ${NAMESPACE} $(kubectl get secret -n ${NAMESPACE} | grep cicd-token | awk '{print $1}') -o jsonpath='{.data.token}' | base64 --decode)
    kubectl --insecure-skip-tls-verify --kubeconfig="/dev/null" --server=https://3ba944e4-aede-4db3-af77-41484e6468ff.k8s.ondigitalocean.com --token=$TOKEN $@
fi
