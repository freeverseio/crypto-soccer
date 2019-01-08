pragma solidity ^0.4.24;

import "./CryptoPlayersTeam.sol";
import "./CryptoMetadata.sol";

contract CryptoPlayers is CryptoPlayersTeam, CryptoMetadata("CryptoSoccerPlayers", "CSP") {
}
