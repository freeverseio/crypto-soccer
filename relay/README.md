# client
- cd client
- npm install
- npm start

the client should now be running at http://localhost:8888

# server
- cd server
- go run main.go

the server should now be running at http://localhost:8080

# interaction between server and client
First post to the client to create one ore more wallets:
```sh
curl -v -H "Content-Type: application/json" -X POST http://localhost:8888/createwallet
```
or if you want to use a particular mnemonic (with at least 12 words)
```sh
curl -v -H "Content-Type: application/json" -X POST -d '{"mnemonic":"a b c d e f g h i j k l"}' http://localhost:8888/relay/v1/createwallet
```
The application should respon with a json message similar to the following:
```json
{
    "success":"true",
    "entry":
    {
        "id":0,
        "account":"0x82973f0ceed111576c508bcd999c92c9e83e49f0",
        "privatekey":"0x8db7ef46326102762035e5d165f5e1bd69b8b97769b050b25b5d563c6cf2419b",
        "mnemonic":"easy squirrel priority convince green shift random gesture arena body frozen summer"
    }
}
```
Account is what we will use to communicate with the server from now on. That is
```
0x82973f0ceed111576c508bcd999c92c9e83e49f0
```

Once a wallet has been created, register it with the server:
```http
http://localhost:8888/relay/v1/0x82973f0ceed111576c508bcd999c92c9e83e49f0
```
Once registered, user actions can be submitted specifying a type and a value:
```http
http://localhost:8888/relay/v1/0x82973f0ceed111576c508bcd999c92c9e83e49f0/action?type=sell&value=player
```
This will sign the current action and will be verified by the server before submitting the action. If the action was succesfully submitted you should obtain a response like the following
```json
{
    "success":true,
    "useraddr":"0x82973f0ceed111576c508bcd999c92c9e83e49f0",
    "action": {
        "Type":"sell",
        "Value":"player"
        },
    "verified":true
}
```

Wallets created by the client application can be inspected in:
```http
http://localhost:8888/relay/debug
```

Current server storage can be inspected in:
```http
http://localhost:8080/relay/db
```