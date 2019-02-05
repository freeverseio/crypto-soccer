pragma solidity ^ 0.4.24;

contract Leagues {
    struct League {
        // teams ids in the league
        uint256[] _teamIds;
        // init block of the league
        uint256 _initBlock;
        // step blocks of the league
        uint256 _step;
        // hash of the init status of the league
        bytes32 _initHash;
        // hash of the state of the league
        bytes32 _hash;
    }

    mapping(uint256 => League) private _leagues;

    function getInitBlock(uint256 id) external view returns (uint256) {
        require(_exists(id), "unexistent league");
        return _leagues[id]._initBlock;
    }

    function getStep(uint256 id) external view returns (uint256) {
        require(_exists(id), "unexistent league");
        return _leagues[id]._step;
    }

    function getEndBlock(uint256 id) external view returns (uint256) {
        require(_exists(id), "unexistent league");
        uint256 nTeams = _leagues[id]._teamIds.length;
        uint256 nMatchDays = 2 * (nTeams - 1);
        return _leagues[id]._initBlock + (nMatchDays - 1) * _leagues[id]._step;
    }

    function getTeamIds(uint256 id) external view returns (uint256[] memory) {
        require(_exists(id), "unexistent league");
        return _leagues[id]._teamIds;
    }

    function getInitHash(uint256 id) external view returns (bytes32) {
        require(_exists(id), "unexistent league");
        return _leagues[id]._initHash;
    }

    function getHash(uint256 id) external view returns (bytes32) {
        require(_exists(id), "unexistent league");
        return _leagues[id]._hash;
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
        return _leagues[id]._initBlock != 0;
    }
}