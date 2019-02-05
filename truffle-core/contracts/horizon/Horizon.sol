pragma solidity ^0.4.24;

import "../HelperFunctions.sol";
import "../ERC721/Players.sol";
import "../ERC721/Teams.sol";

contract Horizon is HelperFunctions {
    Players private _players;
    Teams private _teams;

    constructor(address teams) public {
        _teams = Teams(teams);
        _players = Players(_teams.getPlayersAddress());
    }

    function createTeam(string name) public {
        _teams.mint(msg.sender, name);
        uint256 teamId = _teams.getTeamId(name);

        for (uint i = 0; i<11; i++) {
            string memory postFix = uint2str(i);
            string memory playerName = strConcat(name, "_", postFix);
            _players.mint(msg.sender, playerName);
            uint256 playerId = _players.getPlayerId(playerName);
            _teams.addPlayer(teamId, playerId);
        }
    }
}