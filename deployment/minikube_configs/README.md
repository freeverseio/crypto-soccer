
- install minikube
- start minikube
```
minikube start --vm-driver=hyperkit
```
- connect docker and minikube
```
eval $(minikube docker-env)
```


# UNIVERSE DB and API

- build universe.db image:
```
cd ../../universe.db && docker build -t universedb:0.0.1 .
```
- build universe.api image:
```
cd ../../universe.api && docker build -t universeapi:0.0.1 .
```
- check containers are available
```
docker ps
```
- create a namespace called universe for this
```
kubectl create ns universe
```

- start deployment and (1) pod. For more pods increase replicas number.
```
kubectl apply -f universe.yaml
```
- open service on your browser:
```
minikube service universe -n universe
```
this will return an address similar to http://192.168.64.2:30770 . Append graphiql to it, so it looks like http://192.168.64.2:30770/graphiql

