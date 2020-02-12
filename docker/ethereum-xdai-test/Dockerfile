#### address: 0x83a909262608c650bd9b0ae06e29d90d0f67ac5e 
#### private key: FE058D4CE3446218A7B4E522D9666DF5042CF582A44A9ED64A531A81E7494A85

FROM parity/parity:v2.7.2-stable

COPY ./devnet.chain /devnet.chain
RUN mkdir -p /home/parity/.local/share/io.parity.ethereum/

EXPOSE 8545
EXPOSE 30303
EXPOSE 30303/udp

ENTRYPOINT ["/bin/parity","--chain","/devnet.chain","--no-discovery","--jsonrpc-interface","all","--mode","offline", "--logging","info"]
