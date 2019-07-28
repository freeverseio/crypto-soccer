FROM alpine:3.10
RUN apk update
RUN apk upgrade
RUN apk add --no-cache nodejs nodejs-npm go lsof python3 python git alpine-sdk openssh docker geth
RUN npm install -g ganache-cli
# Build a geth bootnode, regular node from the develop branch and unit test binary
RUN apk add --update git make gcc musl-dev curl jq linux-headers

# Checkout source, pin version to ensure that RPC test suite doesn't break due
# to incompatibility issues in the client packages/generated abi sources.
RUN (git clone -b master --depth 1 --single-branch https://github.com/ethereum/go-ethereum /go/src/github.com/ethereum/go-ethereum)

# Build binaries
RUN (cd /go/src/github.com/ethereum/go-ethereum && GOPATH=/go make all)
RUN (cd /go/src/github.com/ethereum/go-ethereum && GOPATH=/go go get)

# Build and copy binaries to /
RUN (cd /go/src/github.com/ethereum/go-ethereum && GOPATH=/go make all)
RUN cp /go/src/github.com/ethereum/go-ethereum/build/bin/abigen /bin/abigen
RUN rm -rf ./go
