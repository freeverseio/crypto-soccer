pragma solidity ^0.4.24;

import "../ERC721/CryptoPlayersMetadata.sol";

/**
 * @title CryptoPlayersMetadataMock
 * This mock just provides a public mint and burn functions for testing purposes
 */
contract CryptoPlayersMetadataMock is CryptoPlayersMetadata {
    function setTokensURI(string uri) public {
        super._setTokensURI(uri);
    }
}