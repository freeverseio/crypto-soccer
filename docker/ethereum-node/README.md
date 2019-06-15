# Setting a local ethereum node via docker

docker build -t testnode .

docker run -p 8545:8545 testnode

truffle migrate --network dockertest --reset