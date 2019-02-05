pragma solidity ^0.4.24;

import "./PlayersTeam.sol";
import "./ERC721MetadataBaseURI.sol";

contract Players is PlayersTeam, ERC721MetadataBaseURI("Players", "P") {
}
