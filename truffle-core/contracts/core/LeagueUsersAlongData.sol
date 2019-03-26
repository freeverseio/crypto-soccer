pragma solidity ^0.5.0;

import "./LeaguesScheduler.sol";

contract LeagueUsersAlongData is LeaguesScheduler {
    mapping(uint256 => bytes32) private _usersAlongDataHash;

    function getUsersAlongDataHash(uint256 id) external view returns (bytes32) {
        require(_exists(id), "unexistent league");
        return _usersAlongDataHash[id];
    }

    function updateUsersAlongDataHash(uint256 id, bytes32 usersAlongDataHash) public {
        require(_exists(id), "unexistent league");
        require(!hasFinished(id), "finished league");
        _usersAlongDataHash[id] = usersAlongDataHash;
    }

    function computeUsersAlongDataHash(uint256[] memory teamIds, uint8[3][] memory tactics, uint256[] memory blocks) public pure returns (bytes32) {
        uint256 nTeams = teamIds.length;
        require(tactics.length == nTeams, "teams and tactics mismatch");
        require(blocks.length == nTeams, "teams and blocks mismatch");
        bytes32 usersAlongDataHash;
        for(uint256 i = 0 ; i < teamIds.length ; i++)
            usersAlongDataHash = keccak256(abi.encode(usersAlongDataHash, teamIds[i], tactics[i], blocks[i]));
        return usersAlongDataHash;
    }
}