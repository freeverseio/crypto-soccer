#!/bin/sh
export JOB_NAME=fillgamedb
export DEBUG=true
export UNIVERSE_URL=http://universeapi:4000/graphql
export GAME_URL=http://gameapi:4000/graphql

cat job.yml | envsubst | kubectl apply -f - -n freeverse