pragma solidity ^0.4.24;

import "./CryptoTeamsBase.sol";
import "./CryptoPlayersBase.sol";

contract CryptoTeamsLink is CryptoTeamsBase {
    CryptoPlayersBase private _cryptoPlayers;

    function _playerExists(uint256 playerId) internal view returns (bool) {
        address owner = _cryptoPlayers.ownerOf(playerId);
        return owner != address(0);
    }

    function addPlayer(uint256 teamId, uint256 playerId) public {
        require(_playerExists(playerId), "unexistent player");
        _cryptoPlayers.setTeam(playerId, teamId);
        super.addPlayer(teamId, playerId);
    }

    function transferFrom(address from, address to, uint256 teamId) public {
        super.transferFrom(from, to, teamId);
        uint256[] memory players = getPlayers(teamId);
        uint count = players.length;
        for (uint i = 0 ; i < count ; i++){
            _cryptoPlayers.transferFrom(from, to, players[i]);
        }
    }

    function setPlayersContract(address cryptoPlayers) public {
        _cryptoPlayers = CryptoPlayersBase(cryptoPlayers);
    }
}

