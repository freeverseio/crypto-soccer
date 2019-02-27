pragma solidity ^0.5.0;

contract LeaguesStorage {
    struct League {
        // teams ids in the league
        uint256[] teamIds;
        // init block of the league
        uint256 initBlock;
        // step blocks of the league
        uint256 step;
        // hash of tactics
        bytes32 tacticsHash;
    }

    mapping(uint256 => League) private _leagues;

    function getEndBlock(uint256 id) external view returns (uint256) {
        require(_exists(id), "unexistent league");
        uint256 nTeams = _leagues[id].teamIds.length;
        uint256 nMatchDays = 2 * (nTeams - 1);
        return _leagues[id].initBlock + (nMatchDays - 1) * _leagues[id].step;
    }

    function create(uint256 id, uint256 initBlock, uint256 step, uint256[] memory teamIds) public {
        require(initBlock > 0, "invalid init block");
        require(step > 0, "invalid block step");
        require(teamIds.length > 1, "minimum 2 teams per league");
        require(teamIds.length % 2 == 0, "odd teams count");
        require(!_exists(id), "league already created");
        bytes32 tacticsHash = 0;
        _leagues[id] = League(
            teamIds, 
            initBlock, 
            step, 
            tacticsHash
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