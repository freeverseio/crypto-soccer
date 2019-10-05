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
--gasPrice ${GASPRICE} \
--verbose \
--account="0xFE058D4CE3446218A7B4E522D9666DF5042CF582A44A9ED64A531A81E7494A85,10000000000000000000000" 

#--db ${DATABASE} \
