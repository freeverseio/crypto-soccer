pragma solidity ^0.4.24;

import "../HelperFunctions.sol";
import "../ERC721/CryptoPlayers.sol";
import "../ERC721/CryptoTeams.sol";

contract Horizon is HelperFunctions {
    CryptoPlayers private _cryptoPlayers;
    CryptoTeams private _cryptoTeams;

    constructor(address cryptoPlayers, address cryptoTeams) public {
        _cryptoPlayers = CryptoPlayers(cryptoPlayers);
        _cryptoTeams = CryptoTeams(cryptoTeams);
    }

    function createTeam(string name) public {
        _cryptoTeams.mint(msg.sender, name);
        uint256 teamId = _cryptoTeams.getTeamId(name);

        for (uint i = 0; i<11; i++) {
            string memory postFix = uint2str(i);
            string memory playerName = strConcat(name, "_", postFix);
            _cryptoPlayers.mint(msg.sender, playerName);
            uint256 playerId = _cryptoPlayers.getPlayerId(playerName);
            _cryptoTeams.addPlayer(teamId, playerId);
        }
    }
}