FROM ethereum/client-go:alltools-v1.8.27

# Set the working directory to /app
WORKDIR /test

COPY . /test

# RUN apk update
# RUN apk add nodejs npm alpine-sdk python
# RUN npm i
# RUN ./node_modules/.bin/truffle compile

RUN geth --datadir /data init /test/geth-testnode/genesis.json

EXPOSE 8545

CMD ["geth", "--datadir", "/data", "--mine", "--unlock", "0x83a909262608c650bd9b0ae06e29d90d0f67ac5e", "--password", "/test/geth-testnode/passwd", "--keystore", "/test/geth-testnode/keystore", "--rpc", "--rpcapi", "personal,db,eth,net,web3,txpool,miner,debug", "--rpcvhosts", "*", "--rpccorsdomain", "*", "--rpcaddr", "0.0.0.0", "--networkid", "63819", "--nodiscover"]
# CMD ["geth",  "--networkid", "63819" ,"--datadir", "/", "--mine", "--unlock", "0x83a909262608c650bd9b0ae06e29d90d0f67ac5e", "--password", "/test/geth-testnode/passwd", "--keystore", "/test/geth-testnode/keystore"]
# ENTRYPOINT /bin/sh /test/geth-testnode/run.sh