#!/bin/sh

MY_DIR=`dirname "$0"`
SOLC=solcjs
ABIGEN=~/go/bin/abigen
BUILD_DIR=${MY_DIR}/build
CONTRACT=${MY_DIR}/../data/test.sol

(cd ${MY_DIR} && ${SOLC} --abi ../data/test.sol -o ${BUILD_DIR})
(cd ${MY_DIR} && ${SOLC} --bin ../data/test.sol -o ${BUILD_DIR})
mkdir -p ${MY_DIR}/lionel
mkdir -p ${MY_DIR}/stakers
${ABIGEN} --bin=${BUILD_DIR}/___data_test_sol_SoccerSim.bin --abi=${BUILD_DIR}/___data_test_sol_SoccerSim.abi --pkg=lionel --out=${MY_DIR}/lionel/lionel.go
${ABIGEN} --bin=${BUILD_DIR}/___data_test_sol_Stakers.bin --abi=${BUILD_DIR}/___data_test_sol_Stakers.abi --pkg=stakers --out=${MY_DIR}/stakers/stakers.go
