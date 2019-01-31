pragma solidity ^ 0.4.24;

contract Leagues {
    // teams ids of the league
    uint256[] _teamIds;
    uint256 private _init;
    uint256 private _final;

    function getInit() public view returns (uint256) {
        return _init;
    }

    function getFinal() public view returns (uint256) {
        return _final;
    }

    function getTeamIds() external view returns (uint256[] memory) {
        return _teamIds;
    }

    function createLeague(uint256 init, uint256[] memory teamIds) public {
        require(init != 0, "invalid league init state");
        _init = init;
    }

    function updateLeague(uint256 current) public {
        require(current != 0, "invalid league current state");
        require(_init != 0, "league not created");
        _final = current;
    }
}