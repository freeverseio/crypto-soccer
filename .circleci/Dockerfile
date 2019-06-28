FROM circleci/node:10.16.0
RUN sudo apt-get update
RUN sudo apt-get dist-upgrade
RUN sudo wget https://dl.google.com/go/go1.12.2.linux-amd64.tar.gz
RUN sudo tar -xvf go1.12.2.linux-amd64.tar.gz
RUN sudo mv go /usr/local
RUN sudo npm install -g ganache-cli

ENV GOROOT=/usr/local/go
#ENV GOPATH=$HOME/Projects/Proj1
ENV PATH=$GOPATH/bin:$GOROOT/bin:$PATH


