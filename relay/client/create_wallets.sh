#curl -v -H "Content-Type: application/json" -X POST -d '{"mnemonic":"a b c d e f g h i j k l"}' http://localhost:8888/relay/v1/createwallet
curl -v -H "Content-Type: application/json" -X POST http://localhost:8888/createwallet
echo
