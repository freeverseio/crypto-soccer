pragma solidity ^0.5.0;

import "../ERC721/PlayersMintable.sol";

/**
 * @title PlayersMintableMock
 * This mock just provides a public mint and burn functions for testing purposes
 */
contract PlayersMintableMock is PlayersMintable {
    function computeSkills(uint256 seed) public view returns(uint16[5] memory) {
        return _computeSkills(seed);
    }
}