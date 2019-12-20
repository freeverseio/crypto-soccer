
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
- start deployment and (1) pod. For more pods increase replicas number.
```
kubectl apply -f universe.yaml
```

- expose universe deployment so it's accessible from outside the kubernetes virtual network (i.e. from your browser). This step may also be set in the configuration file... just don't know how to do it yet
```
kubectl expose deployment universe  --type=LoadBalancer --port=4000
```

- open service on your browser:
```
minikube service universe
```
this will return an address similar to http://192.168.64.2:30770 . Append graphiql to it, so it looks like http://192.168.64.2:30770/graphiql
