- start minikube
    ```bash
    minikube start --vm-driver=hyperkit
    ```
- create freeverse namespace and docker credentials to download images
    ```bash
    NAMESPACE=freeverse
    DOCKER_REGISTRY_SERVER=docker.io
    DOCKER_USER=freeversedigitalocean
    DOCKER_PASSWORD=fdggreTRGDSBw45ergseth4hDGHD

    kubectl create ns ${NAMESPACE}
    kubectl create secret docker-registry docker-secret --docker-server=${DOCKER_REGISTRY_SERVER} --docker-username=${DOCKER_USER} --docker-password=${DOCKER_PASSWORD} --docker-email=freeversedigitalocean@freeverse.io -n ${NAMESPACE}
    ```

- run ethereum node and wait until POD is running:
    ```bash
    kubectl apply -f ethereum.yaml -n ${NAMESPACE}
    ```

- build go contracts
    ```bash
    cd freeverseio/cryptosoccer && make contracts
    ```

- build dockder image to deploy contracts (use appropriate values for YOUR_GIT_USR and YOUR_GIT_TOKEN)
    ```bash
    eval $(minikube docker-env)
    cd freeverseio/cryptosoccer/go
    docker build -f Dockerfile.contractdeployment -t contractdeployment:0.0.1 --build-arg GIT_TOKEN_USR=<YOUR_GIT_USER> --build-arg GIT_TOKEN_PWD=<YOUR_GIT_TOKEN> .
    ```

- deploy POD to create contracts
    ```bash
    kubectl apply -f contractdeployment.yaml -n ${NAMESPACE}
    ```

- The logs from the previous POD will show the contract addresses that need to be edited in configmap.yaml. After setting the new contract addresses, we can now run the rest of PODs:

    ```bash
    kubectl apply -f configmap.yaml -n ${NAMESPACE}
    kubectl apply -f universedb.yaml -n ${NAMESPACE}
    kubectl apply -f universeapi.yaml -n ${NAMESPACE}
    kubectl apply -f synchronizer.yaml -n ${NAMESPACE}
    kubectl apply -f ipfsnode.yaml -n ${NAMESPACE}
    kubectl apply -f relaydb.yaml -n ${NAMESPACE}
    kubectl apply -f relayapi.yaml -n ${NAMESPACE}
    ```

- open universeapi service on your browser:
    ```bash
    minikube service universeapi -n ${NAMESPACE}
    ```
this will return an address similar to http://192.168.64.2:30770 . Append graphiql to it, so it looks like http://192.168.64.2:30770/graphiql
