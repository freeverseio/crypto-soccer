#!/bin/sh

DATABASE=/tmp/freeverse/ganache
PORT=8545
NETWORKID=5777

mkdir -p ${DATABASE}
ganache-cli \
--deterministic \
--mnemonic "much repair shock carbon improve miss forget sock include bullet interest solution" \
--port ${PORT} \
--networkId ${NETWORKID} \
--db ${DATABASE} \
--blockTime 1 \
--verbose
#--noVMErrorsOnRPCResponse \
