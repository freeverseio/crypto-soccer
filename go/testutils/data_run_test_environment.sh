#!/bin/bash

docker load --input universe.db.with.data.tar.gz 
docker-compose -f ./docker-compose-with-uni-data.yml down
docker-compose -f ./docker-compose-with-uni-data.yml build
docker-compose -f ./docker-compose-with-uni-data.yml up
