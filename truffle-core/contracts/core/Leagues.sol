pragma solidity ^ 0.4.24;

contract Leagues {
    // teams ids in the league
    uint256[] private _teamIds;
    // init block of the league
    uint256 private _blockInit;
    // step blocks of the league
    uint256 private _blockStep;

    function getBlockInit() external view returns (uint256) {
        return _blockInit;
    }

    function getBlockStep() external view returns (uint256) {
        return _blockStep;
    }

    function getTeamIds() external view returns (uint256[] memory) {
        return _teamIds;
    }

    function create(uint256 blockInitDelta, uint256 blockStep, uint256[] memory teamIds) public {
        require(blockInitDelta > 0, "invalid init block");
        require(blockStep > 0, "invalid block step");
        require(teamIds.length > 1, "minimum 2 teams per league");
        _teamIds = teamIds;
        _blockInit = block.number + blockInitDelta;
        _blockStep = blockStep;
    }
}