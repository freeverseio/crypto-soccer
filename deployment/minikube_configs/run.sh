UNIVERSE_NAMESPACE=universe

kubectl create ns ${UNIVERSE_NAMESPACE}
kubectl apply -f configmap.yaml -n ${UNIVERSE_NAMESPACE}
kubectl apply -f universedb.yaml -n ${UNIVERSE_NAMESPACE}
kubectl apply -f universeapi.yaml -n ${UNIVERSE_NAMESPACE}
