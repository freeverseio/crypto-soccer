pragma solidity ^0.4.24;

import "../ERC721/CryptoTeamsStorage.sol";

/**
 * @title CryptoTeamsStorageMock
 * This mock just provides a public mint and burn functions for testing purposes
 */
contract CryptoTeamsStorageMock is CryptoTeamsStorage {
    function mint(address to, uint256 tokenId) public {
        _mint(to, tokenId);
    }
    
    function setName(uint256 teamId, string name) public {
        _setName(teamId, name);
    }

    function addPlayer(uint256 teamId, uint256 playerId) public {
        _addPlayer(teamId, playerId);
    }
}