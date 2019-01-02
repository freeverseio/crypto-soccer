pragma solidity ^0.4.24;

import "../ERC721/CryptoPlayersStorage.sol";

/**
 * @title CryptoPlayersStorageMock
 * This mock just provides a public mint and burn functions for testing purposes
 */
contract CryptoPlayersStorageMock is CryptoPlayersStorage {
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

    function setGenome(
        uint256 playerId,
        uint16 defence,
        uint16 speed,
        uint16 pass,
        uint16 shoot,
        uint16 endurance
    ) public {
        _setGenome(playerId, defence, speed, pass, shoot, endurance);
    }
}