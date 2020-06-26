FROM openethereum/openethereum:v3.0.1
COPY ./data /data

RUN ["./openethereum","account","import","--chain","/data/devnet.chain","/data/eb3ce112d8610382a994646872c4361a96c82cf8"]

RUN mkdir -p $HOME/.local/share/io.parity.ethereum/
RUN mv $HOME/.local/share/io.parity.ethereum/keys $HOME/keys

EXPOSE 8545
EXPOSE 30303
EXPOSE 30303/udp

ENTRYPOINT ["./openethereum","--chain","/data/devnet.chain","--no-discovery","--jsonrpc-interface","all","--mode","offline", "--unlock", "0xeb3ce112d8610382a994646872c4361a96c82cf8", "--password", "/data/passwd", "--logging","info", "--keys-path","$HOME/keys"]
