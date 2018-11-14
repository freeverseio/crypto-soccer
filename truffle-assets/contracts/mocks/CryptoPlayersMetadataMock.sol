pragma solidity ^0.4.24;

import "../CryptoPlayersMetadata.sol";

/**
 * @title CryptoPlayersMetadataMock
 * This mock just provides a public mint and burn functions for testing purposes
 */
contract CryptoPlayersMetadataMock is CryptoPlayersMetadata {
    constructor( string CID) public 
    CryptoPlayersMetadata(CID)
    {
    }

    function mint(address to, uint256 tokenId, uint256 state) public {
        _mint(to, tokenId);
        _setState(tokenId, state);
    }

    function burn(uint256 tokenId) public {
        _burn(ownerOf(tokenId), tokenId);
    }
}