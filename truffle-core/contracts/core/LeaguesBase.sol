pragma solidity ^0.5.0;

contract LeaguesBase {
    struct League {
        // teams ids in the league
        uint256[] teamIds;
        // init block of the league
        uint256 initBlock;
        // step blocks of the league
        uint256 step;
        // hash of the init status of the league 
        bytes32 initStateHash;
        // hash of the day hashes of the league
        bytes32[] dayStateHashes;
        // hash of tactics
        bytes32 tacticsHash;
        // scores of the league 
        uint16[] scores;
    }

    mapping(uint256 => League) private _leagues;

    function create(uint256 id, uint256 initBlock, uint256 step, uint256[] memory teamIds) public {
        require(initBlock > 0, "invalid init block");
        require(step > 0, "invalid block step");
        require(teamIds.length > 1, "minimum 2 teams per league");
        require(teamIds.length % 2 == 0, "odd teams count");
        require(!_exists(id), "league already created");
        bytes32 initStateHash;
        bytes32[] memory dayStateHashes;
        bytes32 tacticsHash;
        uint16[] memory scores;
        _leagues[id] = League(
            teamIds, 
            initBlock, 
            step,
            initStateHash,
            dayStateHashes,
            tacticsHash,
            scores
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
    
    function getInitStateHash(uint256 id) external view returns (bytes32) {
        require(_exists(id), "unexistent league");
        return _leagues[id].initStateHash;
    }

    function getDayStateHashes(uint256 id) public view returns (bytes32[] memory) {
        require(_exists(id), "unexistent league");
        return _leagues[id].dayStateHashes;
    }

    function _setInitStateHash(uint256 id, bytes32 stateHash) internal {
        require(_exists(id), "unexistent league");
        _leagues[id].initStateHash = stateHash;
    }

    function _setDayStateHashes(uint256 id, bytes32[] memory hashes) internal {
        require(_exists(id), "unexistent league");
        _leagues[id].dayStateHashes = hashes;
    }

    function _setScores(uint256 id, uint16[] memory leagues) internal {
        require(_exists(id), "unexistent league");
        _leagues[id].scores = leagues;
    }

    function getScores(uint256 id) external view returns (uint16[] memory) {
        require(_exists(id), "unexistent league");
        return _leagues[id].scores;
    }
}