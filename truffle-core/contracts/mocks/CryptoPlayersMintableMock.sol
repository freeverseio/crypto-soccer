pragma solidity ^0.4.24;

import "../ERC721/CryptoPlayersMintable.sol";

/**
 * @title CryptoPlayersMintableMock
 * This mock just provides a public mint and burn functions for testing purposes
 */
contract CryptoPlayersMintableMock is CryptoPlayersMintable {
    function computeId(string name) public pure returns (uint256) {
        return _computeId(name);
    }

    function computeSkills(uint256 seed) public view returns(uint16[5]) {
        return _computeSkills(seed);
    }
}