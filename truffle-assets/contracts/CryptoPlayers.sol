pragma solidity ^0.4.24;

import "openzeppelin-solidity/contracts/token/ERC721/ERC721MetadataMintable.sol";

contract CryptoPlayers is ERC721MetadataMintable {
    constructor(string name, string symbol) public 
    ERC721Metadata(name, symbol)
    {
    // register the supported interfaces to conform to ERC721 via ERC165
    }
}
