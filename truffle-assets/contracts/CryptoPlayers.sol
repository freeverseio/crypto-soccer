pragma solidity ^0.4.24;

import "openzeppelin-solidity/contracts/token/ERC721/ERC721Full.sol";

contract CryptoPlayers is ERC721Full {
    constructor(string name, string symbol) public 
    ERC721Full(name, symbol)
    {
    }
}
