pragma solidity ^0.4.24;

import "../ERC721/CryptoPlayersMintable.sol";

/**
 * @title CryptoPlayersMintableMock
 * This mock just provides a public mint and burn functions for testing purposes
 */
contract CryptoPlayersMintableMock is CryptoPlayersMintable {
    function setTeam(uint256 playerId, uint256 teamId) public {
        _setTeam(playerId, teamId);
    }

    function computeId(string name) public pure returns (uint256) {
        return _computeId(name);
    }

    function computeState(uint256 rndSeed, uint8 playerRole) public view returns(uint256) {
        return _computeState(rndSeed, playerRole);
    }
}