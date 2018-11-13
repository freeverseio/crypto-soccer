pragma solidity ^0.4.24;

import "../CryptoPlayers.sol";

/**
 * @title CryptoPlayersMock
 * This mock just provides a public mint and burn functions for testing purposes
 */
contract CryptoPlayersMock is CryptoPlayers {
    constructor(string name, string symbol, string CID) public 
    CryptoPlayers(name, symbol, CID)
    {
    }

    function mint(address to, uint256 tokenId, uint256 state) public {
        _mint(to, tokenId);
    }

    function burn(uint256 tokenId) public {
        _burn(ownerOf(tokenId), tokenId);
    }
}