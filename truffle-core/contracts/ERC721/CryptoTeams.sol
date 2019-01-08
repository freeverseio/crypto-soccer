pragma solidity ^0.4.24;

import "./CryptoTeamsPlayers.sol";
import "./CryptoMetadata.sol";

contract CryptoTeams is CryptoTeamsPlayers, CryptoMetadata("CryptoSoccerTeams", "CST") {
}

