pragma solidity ^ 0.4.24;

contract Leagues {
    // teams ids in the league
    uint256[] private _teamIds;
    // init block of the league
    uint256 private _initBlock;
    // step blocks of the league
    uint256 private _step;
    // hash of the init status of the league
    bytes32 private _initHash;
    // hash of the state of the league
    bytes32 private _hash;

    function getInitBlock(uint256 id) external view returns (uint256) {
        return _initBlock;
    }

    function getStep(uint256 id) external view returns (uint256) {
        return _step;
    }

    function getEndBlock(uint256 id) external view returns (uint256) {
        uint256 nTeams = _teamIds.length;
        uint256 nMatchDays = 2 * (nTeams - 1);
        return _initBlock + (nMatchDays - 1) * _step;
    }

    function getTeamIds(uint256 id) external view returns (uint256[] memory) {
        return _teamIds;
    }

    function getInitHash(uint256 id) external view returns (bytes32) {
        return _initHash;
    }

    function getHash(uint256 id) external view returns (bytes32) {
        return _hash;
    }

    function create(uint256 blocksToInit, uint256 step, uint256[] memory teamIds) public {
        require(step > 0, "invalid block step");
        require(teamIds.length > 1, "minimum 2 teams per league");
        _teamIds = teamIds;
        _initBlock = block.number + blocksToInit;
        _step = step;
    }
}