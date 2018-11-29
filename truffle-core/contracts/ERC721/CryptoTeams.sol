pragma solidity ^0.4.24;

import "./CryptoTeamsBase.sol";
import "./CryptoTeamsMetadata.sol";

contract CryptoTeams is CryptoTeamsBase {
    function addTeam(string memory name, address owner) public {
        uint256 nextTeamId = totalSupply() + 1;
        _mint(owner, nextTeamId, name);
    }

    function getTeamName(uint idx) public view returns(string) { 
        require(_exists(idx));
        return _teamProps[idx].name;
    }

    function setTeamPlayersIdx(uint256 team, uint256 playersIdx) public {
        require(_exists(team));
        _teamProps[team].playersIdx = playersIdx;
    }

    function getTeamPlayersIdx(uint256 team) public view returns (uint256) {
        require(_exists(team));
        return _teamProps[team].playersIdx;
    }
}

