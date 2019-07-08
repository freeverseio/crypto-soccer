FROM alpine:3.10
RUN apk update
RUN apk upgrade
RUN apk add --no-cache nodejs nodejs-npm go lsof python3 python git alpine-sdk openssh
RUN npm install -g ganache-cli
# RUN sudo apt-get -y install lsof
# RUN sudo apt-get -y install python3.6 
# RUN sudo wget https://dl.google.com/go/go1.12.2.linux-amd64.tar.gz
# RUN sudo tar -xvf go1.12.2.linux-amd64.tar.gz
# RUN sudo mv go /usr/local
# RUN sudo npm install -g ganache-cli

# ENV GOROOT=/usr/local/go
# #ENV GOPATH=$HOME/Projects/Proj1
# ENV PATH=$GOPATH/bin:$GOROOT/bin:$PATH


