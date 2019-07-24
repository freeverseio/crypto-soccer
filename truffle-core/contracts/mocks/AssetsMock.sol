pragma solidity >=0.4.21 <0.6.0;

import "../assets/Assets.sol";

contract AssetsMock is Assets {

    function computeSkills(uint256 rnd) public pure returns (uint16[5] memory) {
        return _computeSkills(rnd);
    }

    function intHash(string memory arg) public pure returns (uint256) {
        return _intHash(arg);
    }

    function computeBirth(uint256 rnd, uint256 currentTime) public pure returns (uint16) {
        // uint256 currentTime = 1557495456;
        return _computeBirth(rnd, currentTime);
    }

    function computeSeed(bytes32 teamNameHash, uint256 posInTeam) public pure returns (uint256) {
        return _computeSeed(teamNameHash, posInTeam);
    }

    function playerExists(uint256 playerId) public view returns (bool) {
        return _playerExists(playerId);
    }

    function isPlayerVirtual(uint256 playerId) public view returns (bool) {
        return _isPlayerVirtual(playerId);
    }

    function setPlayerState(uint256 state) public {
        _setPlayerState(state);
    }

    function teamExists(uint256 teamId) public view returns (bool){
        return _teamExists(teamId);
    }
    
    function signTeamInLeagueMock(
        uint256 leagueId, 
        uint256 teamId, 
        uint8[PLAYERS_PER_TEAM] memory teamOrder, 
        uint8 teamTactics
    ) public {
        _signTeamInLeague(leagueId, teamId, teamOrder, teamTactics);
    }

}