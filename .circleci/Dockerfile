FROM alpine:3.11
RUN apk update
RUN apk upgrade
RUN apk add --no-cache nodejs nodejs-npm go lsof python3 python git alpine-sdk openssh docker py-pip python-dev libffi-dev openssl-dev gcc libc-dev make geth postgresql-client
# abigen compiling
RUN (git clone -b master --single-branch https://github.com/ethereum/go-ethereum /go-ethereum) 
RUN (cd /go-ethereum && git checkout v1.9.11)
# Build binaries
RUN (cd /go-ethereum/cmd/abigen && go build -o abigen main.go)
# Build and copy binaries to /
RUN cp /go-ethereum/cmd/abigen/abigen /usr/bin/abigen
RUN rm -rf /root/go /go-ethereum

RUN wget https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize%2Fv3.5.4/kustomize_v3.5.4_linux_amd64.tar.gz
RUN tar xvzf ./kustomize_v3.5.4_linux_amd64.tar.gz
RUN rm ./kustomize_v3.5.4_linux_amd64.tar.gz
RUN chmod +x ./kustomize
RUN mv ./kustomize /usr/local/bin