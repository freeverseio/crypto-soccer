#!/bin/bash -xe

export HORIZON_URL=http://localhost:4000/graphql
export PG_CONNECTION_STRING="postgresql://freeverse:freeverse@localhost:5432/game"

node scripts/fillGameDb.js