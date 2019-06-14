Descrizione:- ERC721 team + ERC721 player ?
# [Webpage](freeverse.io)

* migrate to google hosting
 
TODO web secciones:
* what is cryptosoccer
* crear RoadMap
* presentar team
* Whitepaper

[example webpage chintai.io](https://www.chintai.io/)

# Presentation of the web infrastructure design.
ERC721Metadata::tokenURI(uint256 tokenId) returns a JSON schema:

```
{
   “title”: “Asset Metadata”,
   “type”: “object”,
   “properties”: {
       “name”: {
           “type”: “string”,
           “description”: “Identifies the asset to which this NFT represents”,
       },
       “description”: {
           “type”: “string”,
           “description”: “Describes the asset to which this NFT represents”,
       },
       “image”: {
           “type”: “string”,
           “description”: “A URI pointing to a resource with mime type image/* representing the asset to which this NFT represents. Consider making any images at a width between 320 and 1080 pixels and aspect ratio between 1.91:1 and 4:5 inclusive.“,
       }
   }
}
```

## Players web:
URI: https://freeverse.io/players/?state=<player_state>
* The player image is calculated on the fly in front of his state.
* Player_state contains team_id => it can be used to add a second layer of customization: team t-shirt.

## Teams web:
URI: https://freeverse.io/teams/<team_is>

# Action point:
* crear mockup de la web del player 
* crear mockup de la web del team 
* mapping del status del player

# Decoupling engine from tokens ?
* how much cost a remote contract call vs local call?


# Setting a local ethereum node via docker

docker build -t testnode .
docker run -p 8545:8545 testnode
truffle migrate --network dockertest --reset
