pragma solidity ^0.4.24;

import "../ERC721/CryptoTeams.sol";

/**
 * @title CryptoTeamsMock
 * This mock just provides a public mint and burn functions for testing purposes
 */
contract CryptoTeamsMock is CryptoTeams {
    constructor (address cryptoPlayers) CryptoTeams(cryptoPlayers) public {}

    function mint(address to, uint256 tokenId) public {
        _mint(to, tokenId);
    }

    function burn(uint256 tokenId) public {
        _burn(ownerOf(tokenId), tokenId);
    }
}