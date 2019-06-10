pragma solidity >=0.4.21 <0.6.0;

import "../assets/Assets.sol";

contract AssetsMock is Assets {
    constructor(address playerState) public Assets(playerState) {
    }

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

    function computeSeed(string memory teamName, uint256 posInTeam) public pure returns (uint256) {
        return _computeSeed(teamName, posInTeam);
    }

    function playerExists(uint256 playerId) public view returns (bool) {
        return _playerExists(playerId);
    }

    function isVirtual(uint256 playerId) public view returns (bool) {
        return _isVirtual(playerId);
    }

    function setPlayerState(uint256 state) public {
        _setPlayerState(state);
    }

    function teamExists(uint256 teamId) public view returns (bool){
        return _teamExists(teamId);
    }
}