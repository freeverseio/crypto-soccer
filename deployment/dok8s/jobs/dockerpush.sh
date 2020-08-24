#!/bin/sh
export DOCKER_PASSWD=foo
export DOCKER_ID=bar
cd ../../../go && docker build -f Dockerfile.jobs -t freeverseio/jobs:latest . --no-cache

# push images
echo $DOCKER_PASSWD | docker login -u $DOCKER_ID --password-stdin
docker push freeverseio/jobs
