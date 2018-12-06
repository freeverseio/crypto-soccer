pragma solidity ^0.4.24;

import "../ERC721/CryptoTeamsMetadata.sol";

/**
 * @title CryptoTeamsMetadataMock
 * This mock just provides a public mint and burn functions for testing purposes
 */
contract CryptoTeamsMetadataMock is CryptoTeamsMetadata {
    function setPlayersIds(uint256 tokenId, uint256 playersIdx) public {
        _setPlayersIds(tokenId, playersIdx);
    }

    function setTokensURI(string uri) public {
        super._setTokensURI(uri);
    }
}