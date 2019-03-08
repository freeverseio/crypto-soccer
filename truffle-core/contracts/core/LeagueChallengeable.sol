pragma solidity ^0.5.0;

import "./LeagueUpdatable.sol";

contract LeagueChallengeable is LeagueUpdatable {
    uint256 constant private CHALLENGING_PERIOD_BLKS = 60;

    function getChallengePeriod() external view returns (uint256) {
        return CHALLENGING_PERIOD_BLKS;
    }

    function challengeMatchdayStates(
        uint256 id,
        uint256 leagueDay,
        uint256[] memory prevLeagueState
        // userInitData
        // userAlongData
    ) 
        public 
    {
        require(_isUpdated(id), "not updated league. No challenge allowed");
        require(!isVerified(id), "not challengeable league");

    }

    function isVerified(uint256 id) public view returns (bool) {
        uint256 endBlock = getEndBlock(id);
        return block.number > endBlock + CHALLENGING_PERIOD_BLKS;
    }
}