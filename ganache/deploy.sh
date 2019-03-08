SOLC=solcjs
ABIGEN=~/go/bin/abigen
BUILD_DIR=build

${SOLC} --abi ../data/test.sol -o ${BUILD_DIR}
#${ABIGEN} --abi=${BUILD_DIR}/___data_test_sol_SoccerSim.abi --pkg=lionel --out=lionel.go
#${ABIGEN} --abi=${BUILD_DIR}/___data_test_sol_Stakers.abi --pkg=stakers --out=stakers.go
${SOLC} --bin ../data/test.sol -o ${BUILD_DIR}
mkdir lionel
mkdir stakers
${ABIGEN} --bin=${BUILD_DIR}/___data_test_sol_SoccerSim.bin --abi=${BUILD_DIR}/___data_test_sol_SoccerSim.abi --pkg=lionel --out=lionel/lionel.go
${ABIGEN} --bin=${BUILD_DIR}/___data_test_sol_Stakers.bin --abi=${BUILD_DIR}/___data_test_sol_Stakers.abi --pkg=stakers --out=stakers/stakers.go
