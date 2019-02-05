pragma solidity ^0.4.24;

import "./TeamsMintable.sol";
import "./PlayersTeam.sol";

contract TeamsPlayers is TeamsMintable {
    PlayersTeam private _Players;

    constructor(address Players) public {
        _Players = PlayersTeam(Players);
    }

    function getPlayersAddress() external view returns (address) {
        return address(_Players);
    }

    function addPlayer(uint256 teamId, uint256 playerId) public {
        require(_playerExists(playerId), "unexistent player");
        _Players.setTeam(playerId, teamId);
        _addPlayer(teamId, playerId);
    }

    function transferFrom(address from, address to, uint256 teamId) public {
        super.transferFrom(from, to, teamId);
        uint256[] memory players = getPlayers(teamId);
        uint count = players.length;
        for (uint i = 0 ; i < count ; i++){
            _Players.transferFrom(from, to, players[i]);
        }
    }

    function _playerExists(uint256 playerId) private view returns (bool) {
        address owner = _Players.ownerOf(playerId);
        return owner != address(0);
    }
}
