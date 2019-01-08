pragma solidity ^0.4.24;

import "../ERC721/ERC721MetadataBaseURI.sol";

/**
 * @title ERC721MetadataBaseURIMock
 * This mock just provides a public mint and burn functions for testing purposes
 */
contract ERC721MetadataBaseURIMock is ERC721MetadataBaseURI {
    constructor (string name, string symbol) ERC721MetadataBaseURI(name, symbol) public {}

    function mint(address to, uint256 tokenId) public {
        _mint(to, tokenId);
    }

    function burn(uint256 tokenId) public {
        _burn(ownerOf(tokenId), tokenId);
    }
}