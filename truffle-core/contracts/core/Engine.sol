pragma solidity ^ 0.4.24;

import "./Leagues.sol";

contract Engine {
    Leagues private _leagues;

    constructor(address leagues) public {
        _leagues = Leagues(leagues);
    }

    function getLeaguesContract() external view returns (address) {
        return address(_leagues);
    }

    function computeLeagueFinalState (
        uint256 initBlock,
        uint256 step,
        uint256 initPlayerState
        ) public view {
    }
}