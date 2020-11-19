#!/bin/bash

docker-compose -f ./docker-compose-with-uni-data.yml down
docker-compose -f ./docker-compose-with-uni-data.yml build
docker-compose -f ./docker-compose-with-uni-data.yml up
