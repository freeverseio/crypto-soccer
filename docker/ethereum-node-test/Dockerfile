FROM ethereum/client-go:alltools-stable

# Set the working directory to /app
WORKDIR /test

COPY ./geth-testnode /test

#### address: 0x83a909262608c650bd9b0ae06e29d90d0f67ac5e 
#### private key: FE058D4CE3446218A7B4E522D9666DF5042CF582A44A9ED64A531A81E7494A85

RUN geth --datadir /test/data init /test/genesis.json

# RUN apk add nodejs npm alpine-sdk python
# RUN npm i
# RUN ./node_modules/.bin/truffle compile
# RUN geth --datadir /data --mine --unlock 0x83a909262608c650bd9b0ae06e29d90d0f67ac5e --password /test/geth-testnode/passwd --keystore /test/geth-testnode/keystore --rpc --rpcapi personal,db,eth,net,web3,txpool,miner,debug --rpcvhosts '*' --rpccorsdomain '*' --rpcaddr '0.0.0.0' --networkid 63819 --nodiscover & sleep 5 && ./node_modules/.bin/truffle migrate --network dockertest --reset

VOLUME /test 

EXPOSE 8545
# EXPOSE 30303
# EXPOSE 30303/udp

ENTRYPOINT ["geth", "--datadir", "/test/data", "--unlock", "0x83a909262608c650bd9b0ae06e29d90d0f67ac5e", "--password", "/test/passwd", "--keystore", "/test/keystore", "--rpc", "--rpcapi", "personal,db,eth,net,web3,txpool,miner,debug", "--rpcvhosts", "*", "--rpccorsdomain", "*", "--rpcaddr", "0.0.0.0", "--networkid", "63819", "--allow-insecure-unlock", "--mine", "--nodiscover"]