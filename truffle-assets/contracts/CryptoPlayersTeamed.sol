pragma solidity ^0.4.24;

import "./CryptoPlayersMintable.sol";
import "./CryptoTeams.sol";

contract CryptoPlayersTeamed is CryptoPlayersMintable {
    CryptoTeams private _criptoTeams;
    // Mapping from player ID to its team ID
    mapping (uint256 => uint256) private _playerTeam;

    constructor(CryptoTeams cryptoTeams)  public {
        _criptoTeams = cryptoTeams;
    }   

    function getTeam(uint256 playerId) external view returns(uint256){
        require(_exists(playerId), "unexistent player");
        return _playerTeam[playerId];
    }
}
