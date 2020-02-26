FROM parity/parity:v2.7.2-stable
COPY ./data /data

RUN ["/bin/parity","account","import","--chain","/data/devnet.chain","/data/eb3ce112d8610382a994646872c4361a96c82cf8"]

RUN mkdir -p /home/parity/.local/share/io.parity.ethereum/

EXPOSE 8545
EXPOSE 30303
EXPOSE 30303/udp

ENTRYPOINT ["/bin/parity","--chain","/data/devnet.chain","--no-discovery","--jsonrpc-interface","all","--mode","offline", "--unlock", "0xeb3ce112d8610382a994646872c4361a96c82cf8", "--password", "/data/passwd", "--logging","info"]
