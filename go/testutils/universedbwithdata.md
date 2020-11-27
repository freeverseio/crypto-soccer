## Method 1

Launch a postgresql container with a global volume and insert or modify the data in that db. In this case we added a bunch of scripts to /docker-entrypoint-initdb.d/ that executes them on the first run.

```
version: "3"
services:
  universe.db:
    build: universe.db.with.data/
    ports:
      - "5432:5432"
    volumes:
      - universe.db-data:/var/lib/postgresql/data
volumes:
  universe.db-data:
```

Create backup folder in pwd and launch:

```
docker run --rm --volumes-from universe.db -v $(pwd)/backup:/backup ubuntu tar cvf /backup/backup.tar /var/lib/postgresql/data
```

Launch a new postgresql container

```
$ docker run -v /dbdata --name dbstore2 ubuntu /bin/bash
```

Untar the tar file inside the last container

```
$ docker run --rm --volumes-from dbstore2 -v $(pwd):/backup ubuntu bash -c "cd /dbdata && tar xvf /backup/backup.tar --strip 1"
```

## Method 2

Launch the postgresql container and copy the sql insert scripts to /docker-entrypoint-initdb.d/


## Method 3

docker build -t freeverseio/universe.db.with.data:latest .

docker save -o universe.db.with.data.tar.gz freeverseio/universe.db.with.data:latest

docker load --input universe.db.with.data.tar.gz 

docker run freeverseio/universe.db.with.data:latest

It will rerun all the sql scripts the same way for every docker run.
