#!/bin/sh

# kubectl port-forward svc/synchronizer-with-universe-db 5432:5432 -n freeverse
# kubectl port-forward svc/universeapi 4000:4000 -n freeverse  

export DEBUG=true
export UNIVERSE_URL=http://localhost:4000/graphql
export UNIVERSE_DB=postgres://freeverse:freeverse@localhost:5432/cryptosoccer?sslmode=disable
export JOB_NAME=regenerateplayernames
./jobs