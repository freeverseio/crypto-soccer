#!/bin/bash

SCRIPT=$(readlink -f "$0")
# Absolute path this script is in
SCRIPTPATH=$(dirname "$SCRIPT")

psql postgres -U postgres  -h localhost -p 5432 -a -f $SCRIPTPATH/../../universe.db/*.sql
