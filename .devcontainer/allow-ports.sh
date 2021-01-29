#!/bin/bash

ipp=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' crypto-soccer_devcontainer_dockerhost_1)
ufw allow from $ipp to any port 8545
ufw allow from $ipp to any port 5432
ufw allow from $ipp to any port 4001
ufw allow from $ipp to any port 8080
ufw allow from $ipp to any port 8081
ufw allow from $ipp to any port 5001
ufw allow from $ipp to any port 9094
ufw allow from $ipp to any port 9095
ufw allow from $ipp to any port 9096