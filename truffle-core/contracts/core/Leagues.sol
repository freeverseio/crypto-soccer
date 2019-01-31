pragma solidity ^ 0.4.24;

contract Leagues {
    // teams ids of the league
    uint256[] private _teamIds;
    // init state hash of the league
    uint256 private _init;
    // final state hash of the league
    uint256 private _final;

    function getInit() external view returns (uint256) {
        return _init;
    }

    function getFinal() external view returns (uint256) {
        return _final;
    }

    function getTeamIds() external view returns (uint256[] memory) {
        return _teamIds;
    }

    function create(uint256 init, uint256[] memory teamIds) public {
        require(init != 0, "invalid league init state");
        require(teamIds.length > 1, "minimum 2 teams per league");
        _teamIds = teamIds;
        _init = init;
    }

    function update(uint256 current) public {
        require(current != 0, "invalid league current state");
        require(_init != 0, "league not created");
        _final = current;
    }
}