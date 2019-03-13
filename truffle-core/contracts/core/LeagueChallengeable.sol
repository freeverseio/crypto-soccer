pragma solidity ^0.5.0;

import "./LeagueUpdatable.sol";

contract LeagueChallengeable is LeagueUpdatable {
    uint256 constant private CHALLENGING_PERIOD_BLKS = 60;

    function getChallengePeriod() external pure returns (uint256) {
        return CHALLENGING_PERIOD_BLKS;
    }

    function challengeInitStates(
        uint256 id,
        uint256[] memory teamIds,
        uint8[3][] memory tactics
    )
        public
    {
        require(_isUpdated(id), "not updated league. No challenge allowed");
        require(!isVerified(id), "not challengeable league");
        require(getUsersInitDataHash(id) == hashUsersInitData(teamIds, tactics), "incorrect user init data");
        uint256 nTeams = getNTeams(id);
    }

    function isVerified(uint256 id) public view returns (bool) {
        uint256 endBlock = getEndBlock(id);
        return block.number > endBlock + CHALLENGING_PERIOD_BLKS;
    }
}