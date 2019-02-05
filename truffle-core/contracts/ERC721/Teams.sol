pragma solidity ^0.4.24;

import "./TeamsPlayers.sol";
import "./ERC721MetadataBaseURI.sol";

contract Teams is TeamsPlayers, ERC721MetadataBaseURI("Teams", "CT") {
    constructor (address Players) TeamsPlayers(Players) public {}
}

