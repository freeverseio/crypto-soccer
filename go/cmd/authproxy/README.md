# AuthProxy

This is an authenticated proxy using a secp256k1 signed token without backend.

## Overview

First, the client has to sign a token, to do it, it must

- `A`) get the ASCII representation of the decimal value of the current unix timestamp in seconds 
- `B`) hash `A` with keccak256 and sign it with a secp256k1 key
- the token is the concatenation of the two fields infixed with a colon `A:B`

must look like `1579870008:FFiOsOr1OEmxrOlv1DQW7Nlhcy5qL4QH3xdooWwLQkZzPVJhLggkWvDSRDT0u+DafEi+7A7m2mBHvpQQidjK8QE=`

The client has to send the normal GraphQL POST call with the token in the HTTP header `Authorization: Bearer <Token>`.

One recieved, the AuthProxy will check that:

- The timestamp is correct. Because there can be some small diffs between server and client clocks, the server is configured to check in a relaxed way the diff between the local time and the provided timestap. Note that client should take advantage of this by reusing the tickets in this time lapse, so the server is going to have the tickets cached and the validation will be faster.
- The signature of the above timestamp is also ok.

Note that server has a ratelimit configuration to add backpressure, so clients should check if the response is a [409](https://httpstatuses.com/429) to retry the call later in time. `X-Rate-Limit-Limit` and `X-Rate-Limit-Duration` headers will be available.

After this, the AuthProxy will send the recieved body to the GraphQL backend, removing all existing headers, and setting the following ones:

- `Content-Type` with `application/json`
- `User-Agent` with `AuthProxy`
- `X-AuthProxy-Address` with the ethereum address of the signer of the token, e.g. `0xd30f74aca0259d0136249fb3ce6b2a0f970a90e3`

## Server configuration

The server has the following configutation options:

- `-gqlurl` the URL of the GraphQL backend
- `-serviceurl` the URL where this service is going to be published
- `-debug` to activate verboose 
- `-timeout` maximum timeout in seconds for a request to be processed
- `-ratelimit` maximum amount of requests per second
- `-backdoor` to activate an especial token `joshua` that bypasses token security checks
- `-gracetime` to define the grace time in seconds between the ticket claimed time and the local proxy time

the `-serviceurl` parameter accepts

- `http;//` on the port `8080` or
- `https://` on the port `443` (the default). Full DNS (no IP) should be specified since the TLS certificate is automatically created by using letsencrypt.org. The certificate is stored internally, and it will be created each time that the server starts and it is not found in the `cache-path` folder.

## Telemetry

The server exports metrics for telemetry in the `0.0.0.0:4000` port, additionally from the exported golang metrics, also exports:

- `authproxy_ops_success` : the total number of processed operations
- `authproxy_ops_failed` : the total number of failed operations
- `authproxy_ops_dropped` : the the total number of droped events
- `authproxy_cache_hits` : the total number of ticket cache hits
- `authproxy_cache_misses` : the total number of ticket cache misses

## Testing

### With curl

`curl -X POST -H "Content-Type: application/json" -H "Authorization: Bearer joshua" --data '{"query":"{allTeams (condition: {owner: \"0x83A909262608c650BD9b0ae06E29D90D0F67aC5e\"}){totalCount}}"} <proxy_url>' 
`

### With apache benchmark

Also test the performance using `ab` (needs `apache2-utils` package), first create the file `data.json` with the content

```
{"query":"{allTeams (condition: {owner: \"0x83A909262608c650BD9b0ae06E29D90D0F67aC5e\"}){totalCount}}"}
```

and run it with

`ab -n 1000 -c 100 -p data.json -H "Authorization: Bearer joshua" <proxy_url>` 
