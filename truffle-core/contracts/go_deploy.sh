#!/bin/sh

MY_DIR=`dirname "$0"`
SOLC=solcjs
ABIGEN=~/go/bin/abigen
BUILD_DIR=${MY_DIR}/build

cd ${MY_DIR}

contracts=()
for dir in assets core game_controller stakers state; do
  for file in ${dir}/*.sol
  do
      contracts+=(${file})
  done
done

echo compiling ${contracts[@]}

# compile solidity
(cd ${MY_DIR} && ${SOLC} --abi ${contracts[@]} -o ${BUILD_DIR})
(cd ${MY_DIR} && ${SOLC} --bin ${contracts[@]} -o ${BUILD_DIR})

# generate go bindings
GO_DESTDIR=${BUILD_DIR}/go
mkdir -p ${BUILD_DIR}/go
ABI_DESTDIR=${BUILD_DIR}/abi
mkdir -p ${ABI_DESTDIR}
BIN_DESTDIR=${BUILD_DIR}/bin
mkdir -p ${BIN_DESTDIR}

find ${BUILD_DIR} -iname "*.abi" | while read name; do
  inputName=`basename $name .abi`
  outputName=`basename \`echo $name|rev | cut -d'_' -f1 | rev\` .abi`
  ${ABIGEN} --bin=${BUILD_DIR}/${inputName}.bin --abi=${BUILD_DIR}/${inputName}.abi --pkg=${outputName} --out=${GO_DESTDIR}/${outputName}.go
  mv ${BUILD_DIR}/${inputName}.abi ${ABI_DESTDIR}/${outputName}.abi
  mv ${BUILD_DIR}/${inputName}.bin ${BIN_DESTDIR}/${outputName}.bin
done
