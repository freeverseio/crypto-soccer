# Goal Revolution Dev Node

Goal Revolution Dev client that will expose a GraphQL endpoint to access the Goal Revolution dev universe
```
http://localhost:4000/graphiql
```

Goal Revolution dev universe without graphical interface
```
http://localhost:4000/graphql
```

#### Download the last version
As this is a development version of node you'll need to manually sync with the latest releases
```
docker-compose pull
```

#### Start
It will start all the services needed to create the local universe. A fresh sync could last more the 4 hours.

```
docker-compose up
```

if you prefear to start it as a daemon:
```
docker-compose up -d
```



#### Stop
```
docker-compose down
```

#### When you want to restart only one of the dbs (e.g. ipfs/ethereum do not need resstart)
```
docker volume ls
docker volume rm name_of_volume
```

