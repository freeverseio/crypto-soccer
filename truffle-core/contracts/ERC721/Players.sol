pragma solidity ^0.5.0;

import "./PlayersTeam.sol";
import "./ERC721MetadataBaseURI.sol";

contract Players is PlayersTeam, ERC721MetadataBaseURI("Players", "P") {
}
