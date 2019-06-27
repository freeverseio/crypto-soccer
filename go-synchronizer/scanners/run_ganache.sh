#!/bin/sh

PORT=8545
NETWORKID=5777
GASPRICE=200000 # very low gas price otherwise we run out of money quickly

ganache-cli \
--deterministic \
--mnemonic "much repair shock carbon improve miss forget sock include bullet interest solution" \
--port ${PORT} \
--networkId ${NETWORKID} \
--blockTime 0 \
--gasLimit 200000000000 \
--gasPrice ${GASPRICE}
--verbose \
--account="0xf1b3f8e0d52caec13491368449ab8d90f3d222a3e485aa7f02591bbceb5efba5,10000000000000000000000"

#--db ${DATABASE} \
