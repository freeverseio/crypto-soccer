pragma solidity ^0.5.0;

import "./LeaguesScheduler.sol";

contract LeagueUsersAlongData is LeaguesScheduler {
    mapping(uint256 => bytes32) private _usersAlongDataHash;

    function getUsersAlongDataHash(uint256 id) external view returns (bytes32) {
        require(_exists(id), "unexistent league");
        return _usersAlongDataHash[id];
    }

    function updateUsersAlongDataHash(uint256 id, uint256 teamId, uint8[3] memory tactic) public {
        require(_exists(id), "unexistent league");
        require(!hasFinished(id), "finished league");
        bytes32 dataHash = keccak256(abi.encode(_usersAlongDataHash[id], teamId, tactic, block.number));
        _usersAlongDataHash[id] = dataHash;
    }
}