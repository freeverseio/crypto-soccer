pragma solidity ^0.4.24;

import "../ERC721/Players.sol";

/**
 * @title PlayersMock
 * This mock just provides a public mint and burn functions for testing purposes
 */
contract PlayersMock is Players {
    function mint(address to, uint256 tokenId) public {
        _mint(to, tokenId);
    }

    function burn(uint256 tokenId) public {
        _burn(ownerOf(tokenId), tokenId);
    }
}