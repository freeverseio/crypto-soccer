pragma solidity ^0.4.24;

import "openzeppelin-solidity/contracts/token/ERC721/ERC721Mintable.sol";
import "openzeppelin-solidity/contracts/token/ERC721/ERC721Enumerable.sol";
import "openzeppelin-solidity/contracts/token/ERC721/ERC721Metadata.sol";

contract CryptoTeams is ERC721Mintable, ERC721Enumerable, ERC721Metadata {
    constructor(string name, string symbol) ERC721Metadata(name, symbol)
        public
    {
    }
}

