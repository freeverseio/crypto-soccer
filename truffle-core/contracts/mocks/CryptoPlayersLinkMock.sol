pragma solidity ^0.4.24;

import "../ERC721/CryptoPlayersLink.sol";

/**
 * @title CryptoPlayersLinkMock
 * This mock just provides a public mint and burn functions for testing purposes
 */
contract CryptoPlayersLinkMock is CryptoPlayersLink {
    function setTeam(uint256 playerId, uint256 teamId) public onlyMinter {
        _setTeam(playerId, teamId);
    }
}