pragma solidity ^0.4.24;

import "./CryptoPlayers.sol";

contract CryptoPlayersTeamed is CryptoPlayers {
    // Mapping from player ID to its team
    mapping (uint256 => uint256) private _playerTeam;

    function getTeam(uint256 playerId) public view returns (uint256) {
        require(_exists(playerId), "unexistent player");
        return _playerTeam[playerId];
    }

    function _setTeam(uint256 playerId, uint256 teamId) internal {
        require(_exists(playerId), "unexistent player");
        _playerTeam[playerId] = teamId;
    }
}
