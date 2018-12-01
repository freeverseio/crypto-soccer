pragma solidity ^0.4.24;

import "../ERC721/CryptoTeamsBase.sol";

/**
 * @title CryptoTeamsBaseMock
 * This mock just provides a public mint and burn functions for testing purposes
 */
contract CryptoTeamsBaseMock is CryptoTeamsBase {
    function mint(address to, uint256 tokenId) public {
        _mint(to, tokenId);
    }

    function burn(uint256 tokenId) public {
        _burn(ownerOf(tokenId), tokenId);
    }
}