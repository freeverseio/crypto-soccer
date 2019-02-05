pragma solidity ^0.4.24;

import "../ERC721/Teams.sol";

/**
 * @title TeamsMock
 * This mock just provides a public mint and burn functions for testing purposes
 */
contract TeamsMock is Teams {
    constructor (address Players) Teams(Players) public {}

    function mint(address to, uint256 tokenId) public {
        _mint(to, tokenId);
    }

    function burn(uint256 tokenId) public {
        _burn(ownerOf(tokenId), tokenId);
    }
}