pragma solidity ^0.4.24;

import "./CryptoTeamsPlayers.sol";
import "./ERC721MetadataBaseURI.sol";

contract CryptoTeams is CryptoTeamsPlayers, ERC721MetadataBaseURI("CryptoTeams", "CT") {
    constructor (address cryptoPlayers) CryptoTeamsPlayers(cryptoPlayers) public {}
}

