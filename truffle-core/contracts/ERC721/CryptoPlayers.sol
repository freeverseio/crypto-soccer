pragma solidity ^0.4.24;

import "./CryptoPlayersTeam.sol";
import "./ERC721MetadataBaseURI.sol";

contract CryptoPlayers is CryptoPlayersTeam, ERC721MetadataBaseURI("CryptoPlayers", "CP") {
}
