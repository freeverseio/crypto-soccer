FROM ethereum/client-go:alltools-v1.8.27

# Set the working directory to /app
WORKDIR /test

RUN apk update
RUN apk add nodejs npm alpine-sdk python

COPY . /test

RUN npm i
RUN ./node_modules/.bin/truffle compile

RUN geth --datadir /test/geth-testnode init /test/geth-testnode/genesis.json
EXPOSE 8545

ENTRYPOINT /bin/sh /test/geth-testnode/run.sh