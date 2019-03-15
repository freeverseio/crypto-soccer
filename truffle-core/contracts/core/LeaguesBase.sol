pragma solidity ^0.5.0;

contract LeaguesBase {
    struct League {
        // teams ids in the league
        uint256[] teamIds;
        // init block of the league
        uint256 initBlock;
        // step blocks of the league
        uint256 step;

        bytes32 usersInitDataHash;
    }

    mapping(uint256 => League) private _leagues;

    function create(uint256 id, uint256 initBlock, uint256 step, uint256[] memory teamIds, uint256[] memory tactics) public {
        require(initBlock > 0, "invalid init block");
        require(step > 0, "invalid block step");
        require(teamIds.length > 1, "minimum 2 teams per league");
        require(teamIds.length % 2 == 0, "odd teams count");
        require(!_exists(id), "league already created");
        bytes32 usersInitDataHash = keccak256(abi.encode(teamIds, tactics));
        _leagues[id] = League(
            teamIds, 
            initBlock, 
            step,
            usersInitDataHash
        );
    }

    function getInitBlock(uint256 id) public view returns (uint256) {
        require(_exists(id), "unexistent league");
        return _leagues[id].initBlock;
    }

    function getStep(uint256 id) public view returns (uint256) {
        require(_exists(id), "unexistent league");
        return _leagues[id].step;
    }

    function getTeamIds(uint256 id) public view returns (uint256[] memory) {
        require(_exists(id), "unexistent league");
        return _leagues[id].teamIds;
    }

    function countTeams(uint256 id) public view returns (uint256) {
        require(_exists(id), "unexistent league");
        return _leagues[id].teamIds.length;
    }

    function _exists(uint256 id) internal view returns (bool) {
        return _leagues[id].initBlock != 0;
    }
}