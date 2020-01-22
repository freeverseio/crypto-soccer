FROM alpine:3.10.3
RUN apk update
RUN apk upgrade
RUN apk add --no-cache nodejs nodejs-npm go lsof python3 python git alpine-sdk openssh docker py-pip python-dev libffi-dev openssl-dev gcc libc-dev make geth
RUN pip install docker-compose
RUN docker-compose --version
# abigen compiling
RUN (git clone -b master --single-branch https://github.com/ethereum/go-ethereum /go-ethereum) 
RUN (cd /go-ethereum && git checkout v1.9.9)
# Build binaries
RUN (cd /go-ethereum/cmd/abigen && go build -o abigen main.go)
# Build and copy binaries to /
RUN cp /go-ethereum/cmd/abigen/abigen /usr/bin/abigen
RUN rm -rf /root/go /go-ethereum