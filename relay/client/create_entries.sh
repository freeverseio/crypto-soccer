#curl -v -H "Content-Type: application/json" -X POST -d '{"account":"alice","mnemonic":"a b c d"}' http://localhost:8888/relay/v1/createuser
#curl -v -H "Content-Type: application/json" -X POST -d '{"account":"bob","mnemonic":"a b c d"}' http://localhost:8888/relay/v1/createuser

#curl -v -H "Content-Type: application/json" -X POST -d '{}' http://localhost:8888/relay/v1/createuser
curl -v -H "Content-Type: application/json" -X POST http://localhost:8888/relay/v1/createuser
echo
