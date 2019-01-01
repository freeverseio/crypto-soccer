pragma solidity ^0.4.24;

import "../ERC721/CryptoTeamsBase.sol";

/**
 * @title CryptoTeamsBaseMock
 * This mock just provides a public mint and burn functions for testing purposes
 */
contract CryptoTeamsBaseMock is CryptoTeamsBase {
    function addPlayer(uint256 teamId, uint256 playerId) public {
        _addPlayer(teamId, playerId);
    }
}