pragma solidity ^0.4.25;

contract Leagues {
    struct League {
        // teams ids in the league
        uint256[] teamIds;
        // init block of the league
        uint256 initBlock;
        // step blocks of the league
        uint256 step;
        // hash of the init status of the league 
        bytes32 initStateHash;
        // hash of the final hashes of the league
        bytes32[] finalTeamStateHashes;
        // hash of tactics
        bytes32 tacticsHash;
        // scores of the league
        uint256[2][] scores;
    }

    mapping(uint256 => League) private _leagues;

    function getEndBlock(uint256 id) external view returns (uint256) {
        require(_exists(id), "unexistent league");
        uint256 nTeams = _leagues[id].teamIds.length;
        uint256 nMatchDays = 2 * (nTeams - 1);
        return _leagues[id].initBlock + (nMatchDays - 1) * _leagues[id].step;
    }
 
    function getScores(uint256 id) external view returns (uint256[2][] memory) {
        require(_exists(id), "unexistent league");
        return _leagues[id].scores;
    }

    function getInitStateHash(uint256 id) external view returns (bytes32) {
        require(_exists(id), "unexistent league");
        return _leagues[id].initStateHash;
    }

    /// TODO: blockToInit -> initBlock: utilize an absolute reference 
    function create(uint256 id, uint256 blocksToInit, uint256 step, uint256[] memory teamIds) public {
        require(step > 0, "invalid block step");
        require(teamIds.length > 1, "minimum 2 teams per league");
        require(!_exists(id), "league already created");
        uint256 initBlock = block.number + blocksToInit;
        bytes32 initStateHash = 0;
        bytes32[] memory finalTeamStateHashes;
        uint256[2][] memory scores;
        bytes32 tacticsHash = 0;
        _leagues[id] = League(
            teamIds, 
            initBlock, 
            step, 
            initStateHash, 
            finalTeamStateHashes, 
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

    function getInitHash(uint256 id) public view returns (bytes32) {
        require(_exists(id), "unexistent league");
        return _leagues[id].initStateHash;
    }

    function getFinalTeamStateHashes(uint256 id) public view returns (bytes32[] memory) {
        require(_exists(id), "unexistent league");
        return _leagues[id].finalTeamStateHashes;
    }
    function countTeams(uint256 id) public view returns (uint256) {
        require(_exists(id), "unexistent league");
        return _leagues[id].teamIds.length;
    }

    function _setInitStateHash(uint256 id, bytes32 stateHash) internal {
        require(_exists(id), "unexistent league");
        _leagues[id].initStateHash = stateHash;
    }

    function _setFinalTeamStateHashes(uint256 id, bytes32[] memory hashes) internal {
        require(_exists(id), "unexistent league");
        _leagues[id].finalTeamStateHashes = hashes;
    }

    function _setScores(uint256 id, uint256[2][] memory scores) internal {
        require(_exists(id), "unexistent league");
        _leagues[id].scores = scores;
    }

    function _exists(uint256 id) private view returns (bool) {
        return _leagues[id].initBlock != 0;
    }
}