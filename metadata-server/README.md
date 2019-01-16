# crypto-soccer-metadata-server
Return a JSON description generated from the token (player/team) id.

# UML Sequence diagram

![](https://www.draw.io/?lightbox=1&highlight=0000ff&edit=_blank&layers=1&nav=1&title=Untitled%20Diagram.xml#R7ZrbdqM2FIafxpfJAoRtfBl7krRdM23WeKaHSxlkrEYgKuQ47tNXAnGSwIfEeFIv5yZooxPan%2F69ER6AWfT6yGCy%2BkIDRAaOFbwOwKeB49iTkSf%2BSctWWRygLCHDgbJVhjn%2BFymjpaxrHKC0UZFTSjhOmkafxjHyecMGGaObZrUlJc1RExgiwzD3ITGtf%2BCAr5TVtqzqxk8Ihys1tDdUNxbQfw4ZXcdqvIEDltlffjuCRV%2BqfrqCAd3UTOB%2BAGaMUp5fRa8zROTiFsuWt3vouFvOm6GYH9JgvIALaPvuwnfGS3dh3Th5Dy%2BQrNVazNg24fSJwC1iqZo13xYrJR4gkZfriHzGS0RwLErTBDEcIY6YuEOU%2BamyTTcrzNE8gb5suhEQCduKR0SUbHEp3MqhaMLKMiEwSfEiG9USFob8NUvxC%2FqK0pweaaVrLkealVRkVaU3UKC6KhfcyvqNsK%2BuCVwgMi3dN6OEyuFjmj1Qyhl9RoVReNXK%2Fso7BSVyiCUmpFbzIfuTdvFUDzDCRO6K3xELYAyVWe0A21HltoEgwWEsbL5wbbaIpq%2BV%2B18Q4%2Bi1ZlK%2Bf0RUOIBtRRV1t9xwaqN6qripQV8wv6rxbo%2BUEaqNFpZdV7CJC8XbgewVk6nBZ%2FBWoyShOObZ%2BMPpYPhJw44yvqIhjSGpg1fBYF06DJ1b%2B2A63GEDjpKEBh0mHL2wAQw0RJ8wgBxKPhATT3XVpo%2BM49HaBJy3aZMz7kWbTACv2nQ6bQLH0qFrUysdJhyTPthwDTRwJBPMqzD9H1g8Wphc8DZhAr0Ik2fA91uC4jmCV%2BouiroSqSOpc60%2BqJtco2GP0dA7Fo43ZuolRO9hg%2FBf507qfGERDJznX%2F6JhktwY7sGDygI0VwVK5ffV1bN5TV%2B%2Fkacb9ViwzWnUj%2FKHj5Tmqh6%2BZhyoMYma9l3KV0zH3XleopuDlmIunw0aXcIQwRyIXaN0dqWN2t6xxjc1iqojVL1%2FCQNtfd1tykCrqsd7xxXX1zkM6gcXT7KG9PkqzD80NRk3BSCMuXYIwSn0AGTBdtggYuljb9%2F%2FXkXE%2FJpsQ%2FJnVqfBeWcRuIGioM7ebYrbYT6z5lJTPJP5fGs8NdB7j%2FQdXXJQmRBN3W1ygziRjHfXR7cKTeHqE3hNE3b3qs%2FpmBoL922%2FsKUP4lqpcHxTu0AQ4OXAKcJ5P7qVLychgrT8ftjm8iBX8rMtsTWbmB763ijwlBTvBzuV8yzFrcW8FQ5bwRE8M%2FLeptOHHeiBtzzoDacNLXK8ZyDUNvfEdCYzR%2FX6OjY4DsaWq0T7pqXXh84pw2%2BrYmXGW7PkHcFMF2VcJ8qCbM%2BahJWaqLGWxcHe%2Br3w4F5JtVyXj4iYmGnC3kV8tJ1NXZE2sGbSDD5Eg%2Brt3szgzPyGl2gIxwEGXcdalx9My1mZCRA5ZdfNZNBGbfqmHVvkM4Eyrp1R3ZTTm7ccV5%2Bb2gFWmrW7IAulyl6b0xtfeLhj5CEE6nAQbkR6CU%2BAVeX%2B37Civ6dozij7pKTPfX7kZORISeP999OIxcCInUuaNfUg6Alb9EOLgGbiqTQx3H4LaNNSN0OnTCA6%2F4Wrh2iOC3vTm7bIUpfZyjexYRycK5QbsbekZaD9ZQb6u9NYE9OsKd%2BP5vYPJvRv1P1kxB8mB3uaKcjtrnDvXPu8GL8yw3NH%2BKUFOyJqHvq97MZW07HLi6kAu3nZW3HkWcNqbZ5xJRJ4G0Shxevfl5T%2FVx7ZDhjfBpniGL1E9Z8w1Q%2FFAb3%2FwE%3D)

Endpoints:
* /players/\<ID>
* /teams/\<ID>
  
# install & run
``` 
$ npm install
$ npm start
```

# config.json
```
{
    "provider": "ws://127.0.0.1:8545", // ex. Ganache
    "crypto_player_address": "0x73E02FefD1d31607b8d5c9Ee9F0c465033C1ebf3", 
    "players_image_base_URL": "https://www.freenode.io/players/image/"
}
```

# ERC721 Metadata JSON Schema 
```
{
    "title": "Asset Metadata",
    "type": "object",
    "properties": {
        "name": {
            "type": "string",
            "description": "Identifies the asset to which this NFT represents"
        },
        "description": {
            "type": "string",
            "description": "Describes the asset to which this NFT represents"
        },
        "image": {
            "type": "string",
            "description": "A URI pointing to a resource with mime type image/* representing the asset to which this NFT represents. Consider making any images at a width between 320 and 1080 pixels and aspect ratio between 1.91:1 and 4:5 inclusive."
        }
    }
}
```

# [OpenSea extended Metadata JSON Schema](https://docs.opensea.io/docs/2-adding-metadata)
```
{
  "attributes": [
    {
      "trait_type": "base", 
      "value": "narwhal"
    }, 
    {
      "trait_type": "eyes", 
      "value": "sleepy"
    }, 
    {
      "trait_type": "mouth", 
      "value": "cute"
    }, 
    {
      "trait_type": "level", 
      "value": 4
    }, 
    {
      "trait_type": "stamina", 
      "value": 90.2
    }, 
    {
      "trait_type": "personality", 
      "value": "boring"
    }, 
    {
      "display_type": "boost_number", 
      "trait_type": "aqua_power", 
      "value": 10
    }, 
    {
      "display_type": "boost_percentage", 
      "trait_type": "stamina_increase", 
      "value": 5
    }, 
    {
      "display_type": "number", 
      "trait_type": "generation", 
      "value": 1
    }
  ], 
  "description": "Friendly OpenSea Creature that enjoys long swims in the ocean.", 
  "external_url": "https://openseacreatures.io/3", 
  "image": "https://storage.googleapis.com/opensea-prod.appspot.com/creature/3.png", 
  "name": "Dave Starbelly"
}
```

