pragma solidity ^0.5.0;

import "./LeaguesScheduler.sol";

contract LeagueUsersAlongData is LeaguesScheduler {
    mapping(uint256 => bytes32) private _usersAlongDataHash;

    function getUsersAlongDataHash(uint256 id) public view returns (bytes32) {
        require(_exists(id), "unexistent league");
        return _usersAlongDataHash[id];
    }

    function updateUsersAlongDataHash(uint256 id, uint256[] memory teamIds, uint8[3][] memory tactics, uint256[] memory blocks) public {
        require(_exists(id), "unexistent league");
        require(!hasFinished(id), "finished league");
        // TODO: do this well with lionel4
        // bytes32 usersAlongDataHash = _usersAlongDataHash[id];
        // usersAlongDataHash = _computeUsersAlongDataHash(usersAlongDataHash, teamIds, tactics, blocks);
        _usersAlongDataHash[id] = keccak256(abi.encode("TODO"));
    }

    function computeUsersAlongDataHash(uint256[] memory teamIds, uint8[3][] memory tactics, uint256[] memory blocks) public pure returns (bytes32) {
        bytes32 base;
        return _computeUsersAlongDataHash(base, teamIds, tactics, blocks);
    }

    function _computeUsersAlongDataHash(bytes32 base, uint256[] memory teamIds, uint8[3][] memory tactics, uint256[] memory blocks) private pure returns (bytes32 usersAlongDataHash) {
        uint256 nTeams = teamIds.length;
        require(tactics.length == nTeams, "teams and tactics mismatch");
        require(blocks.length == nTeams, "teams and blocks mismatch");
        usersAlongDataHash = base;
        for(uint256 i = 0 ; i < nTeams ; i++)
            usersAlongDataHash = keccak256(abi.encode(usersAlongDataHash, teamIds[i], tactics[i], blocks[i]));
        return usersAlongDataHash;
    }
}