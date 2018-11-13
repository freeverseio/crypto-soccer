pragma solidity ^0.4.24;

import "openzeppelin-solidity/contracts/token/ERC721/ERC721Full.sol";

contract CryptoPlayers is ERC721Full {
    struct Props {
        uint8 defense;
        uint8 attack;
    }

    // Mapping from token ID to owner
    mapping (uint256 => Props) private _tokenProps;

    constructor(string name, string symbol) public 
    ERC721Full(name, symbol)
    {
    }
}
