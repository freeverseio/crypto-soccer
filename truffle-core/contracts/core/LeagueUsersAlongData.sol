pragma solidity ^0.5.0;

import "./LeaguesScheduler.sol";

contract LeagueUsersAlongData is LeaguesScheduler {
    mapping(uint256 => bytes32) private _usersAlongDataHash;

    function getUsersAlongDataHash(uint256 id) external view returns (bytes32) {
        require(_exists(id), "unexistent league");
        return _usersAlongDataHash[id];
    }

    function updateUsersAlongDataHash(uint256 id, uint256[] memory teamIds, uint8[3][] memory tactics) public {
        require(_exists(id), "unexistent league");
        require(!hasFinished(id), "finished league");
        require(teamIds.length == tactics.length, "teams and tactics mismatch");
        bytes32 usersAlongDataHash = _usersAlongDataHash[id];
        for(uint256 i = 0 ; i < teamIds.length ; i++)
            usersAlongDataHash = keccak256(abi.encode(usersAlongDataHash, teamIds[i], tactics[i], block.number));
        _usersAlongDataHash[id] = usersAlongDataHash;
    }
}