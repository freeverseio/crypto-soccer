pragma solidity ^ 0.4.24;

contract Leagues {
    // teams ids in the league
    uint256[] private _teamIds;
    // init block of the league
    uint256 private _blockInit;
    // step blocks of the league
    uint256 private _blockStep;
    // hash of the init status of the league
    bytes32 private _initStateHash;
    // hash of the state of the league
    bytes32 private _stateHash;


    function getBlockInit() external view returns (uint256) {
        return _blockInit;
    }

    function getBlockStep() external view returns (uint256) {
        return _blockStep;
    }

    function getBlockEnd(uint256 id) external view returns (uint256) {
        uint256 nTeams = _teamIds.length;
        uint256 nMatchDays = 2 * (nTeams - 1);
        return _blockInit * (nMatchDays - 1) * _blockStep;
    }

    function getTeamIds() external view returns (uint256[] memory) {
        return _teamIds;
    }

    function getInitStateHash(uint256 id) external view returns (bytes32) {
        return _initStateHash;
    }

    function getStateHash(uint256 id) external view returns (bytes32) {
        return _stateHash;
    }

    function create(uint256 blockInitDelta, uint256 blockStep, uint256[] memory teamIds) public {
        require(blockInitDelta > 0, "invalid init block");
        require(blockStep > 0, "invalid block step");
        require(teamIds.length > 1, "minimum 2 teams per league");
        _teamIds = teamIds;
        _blockInit = block.number + blockInitDelta;
        _blockStep = blockStep;
    }

    function hasStarted() external view returns (bool) {
        require(_blockInit != 0, "league not initialized");
        return _blockInit <= block.number;
    }
}