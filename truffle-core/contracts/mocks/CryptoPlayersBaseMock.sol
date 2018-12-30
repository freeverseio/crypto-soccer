pragma solidity ^0.4.24;

import "../ERC721/CryptoPlayersBase.sol";

/**
 * @title CryptoPlayersBaseMock
 * This mock just provides a public mint and burn functions for testing purposes
 */
contract CryptoPlayersBaseMock is CryptoPlayersBase {
    function mint(address to, uint256 tokenId) public {
        _mint(to, tokenId);
    }

    function setState(uint256 playerId, uint256 state) public {
        _setState(playerId, state);
    }

    function setName(uint256 playerId, string name) public {
        _setName(playerId, name);
    }

    function setTeam(uint256 playerId, uint256 teamId) public {
        _setTeam(playerId, teamId);
    }
}