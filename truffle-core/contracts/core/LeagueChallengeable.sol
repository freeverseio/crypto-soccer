pragma solidity ^0.5.0;

import "./LeagueUpdatable.sol";

contract LeagueChallengeable is LeagueUpdatable {
    uint256 constant private CHALLENGING_PERIOD_BLKS = 60;

    function getChallengePeriod() external view returns (uint256) {
        return CHALLENGING_PERIOD_BLKS;
    }

    function challengeMatchdayStates(uint256 id) public {
        require(_isUpdated(id), "not updated league. No challenge allowed");
    }
}