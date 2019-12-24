NAMESPACE=freeverse
DOCKER_REGISTRY_SERVER=docker.io
DOCKER_USER=freeversedigitalocean
DOCKER_PASSWORD=fdggreTRGDSBw45ergseth4hDGHD

kubectl create ns ${NAMESPACE}
kubectl create secret docker-registry docker-secret --docker-server=${DOCKER_REGISTRY_SERVER} --docker-username=${DOCKER_USER} --docker-password=${DOCKER_PASSWORD} --docker-email=freeversedigitalocean@freeverse.io -n ${NAMESPACE}
kubectl apply -f configmap.yaml -n ${NAMESPACE}
kubectl apply -f universedb.yaml -n ${NAMESPACE}
kubectl apply -f universeapi.yaml -n ${NAMESPACE}
kubectl apply -f ethereum.yaml -n ${NAMESPACE}
