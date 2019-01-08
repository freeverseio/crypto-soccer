pragma solidity ^0.4.24;

import "./CryptoTeamsMintable.sol";
import "./CryptoPlayersTeam.sol";

contract CryptoTeamsPlayers is CryptoTeamsMintable {
    CryptoPlayersTeam private _cryptoPlayers;

    constructor(address cryptoPlayers) public {
        _cryptoPlayers = CryptoPlayersTeam(cryptoPlayers);
    }

    function getCryptoPlayers() external view returns (address) {
        return address(_cryptoPlayers);
    }

    function addPlayer(uint256 teamId, uint256 playerId) public {
        require(_playerExists(playerId), "unexistent player");
        _cryptoPlayers.setTeam(playerId, teamId);
        _addPlayer(teamId, playerId);
    }

    function transferFrom(address from, address to, uint256 teamId) public {
        super.transferFrom(from, to, teamId);
        uint256[] memory players = getPlayers(teamId);
        uint count = players.length;
        for (uint i = 0 ; i < count ; i++){
            _cryptoPlayers.transferFrom(from, to, players[i]);
        }
    }

    function _playerExists(uint256 playerId) private view returns (bool) {
        address owner = _cryptoPlayers.ownerOf(playerId);
        return owner != address(0);
    }
}
