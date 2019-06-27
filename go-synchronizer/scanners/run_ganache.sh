#!/bin/sh

#DATABASE=/tmp/freeverse/ganache
PORT=8545
NETWORKID=5777

#mkdir -p ${DATABASE}
ganache-cli \
--deterministic \
--mnemonic "much repair shock carbon improve miss forget sock include bullet interest solution" \
--port ${PORT} \
--networkId ${NETWORKID} \
--blockTime 1 \
--gasLimit 200000000000 \
--verbose \
--account="0xf1b3f8e0d52caec13491368449ab8d90f3d222a3e485aa7f02591bbceb5efba5,1000000000000000000000"

#--db ${DATABASE} \
