FROM ethereum/client-go:alltools-v1.8.27

# Set the working directory to /app
WORKDIR /test

COPY ./genesis.json /test

RUN geth --datadir /test/data init /test/genesis.json

VOLUME /test 

EXPOSE 8545
EXPOSE 30303

CMD ["geth", "--datadir", "/test/data", "--bootnodes", "enode://f1e8d73c81c328938d5c8459c560a3755a9e3c3f72055db6006bb03ba9b80dc3768c33b6826b290223a933babf2d979470a05e3a3574f5960184ead903c7e6ce@165.22.66.118:30303","--rpc", "--rpcapi", "personal,db,eth,net,web3,txpool,miner,debug", "--rpcvhosts", "*", "--rpccorsdomain", "*", "--rpcaddr", "0.0.0.0", "--networkid", "63819"]