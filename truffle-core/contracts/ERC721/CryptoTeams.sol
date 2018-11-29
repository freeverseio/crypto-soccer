pragma solidity ^0.4.24;

import "./CryptoTeamsBase.sol";
import "./CryptoTeamsMetadata.sol";

contract CryptoTeams is CryptoTeamsBase, CryptoTeamsMetadata {
    function addTeam(string memory name, address owner) public {
        uint256 nextTeamId = totalSupply() + 1;
        _mint(owner, nextTeamId, name);
    }

    function getTeamName(uint idx) public view returns(string) { 
        return _getTeamName(idx);
    }

    function setTeamPlayersIdx(uint256 team, uint256 playersIdx) public {
        _setTeamPlayersIdx(team, playersIdx);
    }

    function getTeamPlayersIdx(uint256 team) public view returns (uint256) {
        return _getTeamPlayersIdx(team);
    }
}

