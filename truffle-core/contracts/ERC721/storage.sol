pragma solidity ^ 0.4.24;

import "../CryptoSoccer.sol";
import "./CryptoTeams.sol";
import "./CryptoPlayers.sol";
/*
    Defines all storage structures and mappings
*/

contract Storage is CryptoSoccer, CryptoPlayers {
    CryptoTeams private _cryptoTeams; 

    constructor(address cryptoTeams) public {
        _cryptoTeams = CryptoTeams(cryptoTeams);
    }

    function addTeam(string memory name, address owner) public {
        _cryptoTeams.addTeam(name, owner);
    }

    function getTeamName(uint idx) public view returns(string) { 
        return _cryptoTeams.getTeamName(idx);
    }

    function getNCreatedTeams() public view returns(uint) {
        return _cryptoTeams.getNCreatedTeams();
    }

    function setTeamPlayersIdx(uint256 team, uint256 playersIdx) public {
        _cryptoTeams.setTeamPlayersIdx(team, playersIdx);
    }

    function getTeamPlayersIdx(uint256 team) public returns (uint256) {
        return _cryptoTeams.getTeamPlayersIdx(team);
    }
 
    function teamNameByPlayer(string name) public view returns(string){
        uint256 teamIdx = getTeamIndexByPlayer(name);
        return getTeamName(teamIdx);
    }
}
