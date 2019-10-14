#!/bin/sh

docker pull freeverseio/ethereum:test
docker run -p 8545:8545 freeverseio/ethereum:test
