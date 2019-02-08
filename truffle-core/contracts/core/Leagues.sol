pragma solidity ^ 0.4.24;

contract Leagues {
    struct League {
        // teams ids in the league
        uint256[] teamIds;
        // init block of the league
        uint256 initBlock;
        // step blocks of the league
        uint256 step;
        // hash of the init status of the league
        bytes32 initHash;
        // hash of the state of the league
        bytes32 hash;
    }

    mapping(uint256 => League) private _leagues;

    function getInitBlock(uint256 id) public view returns (uint256) {
        require(_exists(id), "unexistent league");
        return _leagues[id].initBlock;
    }

    function getStep(uint256 id) public view returns (uint256) {
        require(_exists(id), "unexistent league");
        return _leagues[id].step;
    }

    function getEndBlock(uint256 id) external view returns (uint256) {
        require(_exists(id), "unexistent league");
        uint256 nTeams = _leagues[id].teamIds.length;
        uint256 nMatchDays = 2 * (nTeams - 1);
        return _leagues[id].initBlock + (nMatchDays - 1) * _leagues[id].step;
    }

    function getTeamIds(uint256 id) public view returns (uint256[] memory) {
        require(_exists(id), "unexistent league");
        return _leagues[id].teamIds;
    }

    function getInitHash(uint256 id) external view returns (bytes32) {
        require(_exists(id), "unexistent league");
        return _leagues[id].initHash;
    }

    function getHash(uint256 id) external view returns (bytes32) {
        require(_exists(id), "unexistent league");
        return _leagues[id].hash;
    }

    function _setHash(uint256 id, bytes32 hash) internal {
        _leagues[id].hash = hash;
    }

    function countTeams(uint256 id) public view returns (uint256) {
        require(_exists(id), "unexistent league");
        return _leagues[id].teamIds.length;
    }

    /// TODO: blockToInit -> initBlock: utilize an absolute reference 
    function create(uint256 id, uint256 blocksToInit, uint256 step, uint256[] memory teamIds) public {
        require(step > 0, "invalid block step");
        require(teamIds.length > 1, "minimum 2 teams per league");
        require(!_exists(id), "league already created");
        uint256 initBlock = block.number + blocksToInit;
        bytes32 initHash = 0;
        bytes32 hash = 0;
        _leagues[id] = League(teamIds, initBlock, step, initHash, hash);
    }

    function _exists(uint256 id) private view returns (bool) {
        return _leagues[id].initBlock != 0;
    }
}