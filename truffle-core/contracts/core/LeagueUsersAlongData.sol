pragma solidity ^0.5.0;

import "./LeaguesScheduler.sol";

contract LeagueUsersAlongData is LeaguesScheduler {
    mapping(uint256 => bytes32) private _usersAlongDataHash;

    function getUsersAlongDataHash(uint256 id) public view returns (bytes32) {
        require(_leagueExists(id), "unexistent league");
        return _usersAlongDataHash[id];
    }

    function updateUsersAlongDataHash(uint256 id, uint256[] memory teamIds, uint8[] memory tacticsIds, uint256[] memory blocks) public {
        require(_leagueExists(id), "unexistent league");
        require(!hasFinished(id), "finished league");
        // TODO: do this well with lionel4
        bytes32 usersAlongDataHash = _usersAlongDataHash[id];
        usersAlongDataHash = _computeUsersAlongDataHash(usersAlongDataHash, teamIds, tacticsIds, blocks);
        _usersAlongDataHash[id] = keccak256(abi.encode("TODO"));
    }

    function computeUsersAlongDataHash(uint256[] memory teamIds, uint8[] memory tacticsIds, uint256[] memory blocks) public pure returns (bytes32) {
        bytes32 base;
        return _computeUsersAlongDataHash(base, teamIds, tacticsIds, blocks);
    }

    function _computeUsersAlongDataHash(bytes32 base, uint256[] memory teamIds, uint8[] memory tacticsIds, uint256[] memory blocks) private pure returns (bytes32 usersAlongDataHash) {
        uint256 nTeams = teamIds.length;
        require(tacticsIds.length == nTeams, "teams and tacticsIds mismatch");
        require(blocks.length == nTeams, "teams and blocks mismatch");
        usersAlongDataHash = base;
        for(uint256 i = 0 ; i < nTeams ; i++)
            usersAlongDataHash = keccak256(abi.encode(usersAlongDataHash, teamIds[i], tacticsIds[i], blocks[i]));
        return usersAlongDataHash;
    }
}