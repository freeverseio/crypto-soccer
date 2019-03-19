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
        uint8[3][] memory tactics,
        uint256[] memory dataToChallengeInitStates
    )
        public
    {
        require(isUpdated(id), "not updated league. No challenge allowed");
        require(!isVerified(id), "not challengeable league");
        require(getUsersInitDataHash(id) == hashUsersInitData(teamIds, tactics), "incorrect user init data");
        uint256[] memory initPlayerStates = getInitPlayerStates(id, teamIds, tactics, dataToChallengeInitStates);
        if (initPlayerStates.length == 0) // challenger wins
            resetUpdater(id);
        else if (getInitStateHash(id) != hashState(initPlayerStates)) // challenger wins
            resetUpdater(id);
    }

    function challengeMatchdayStates(
        uint256 id,
        uint256 leagueDay,
        uint256[] memory prevMatchdayStates
    )
        public
    {
        require(isUpdated(id), "not updated league. No challenge allowed");
        require(!isVerified(id), "not challengeable league");
        resetUpdater(id);
    }

    function getInitPlayerStates(
        uint256 id,
        uint256[] memory teamIds,
        uint8[3][] memory tactics,
        uint256[] memory dataToChallengeInitStates
    )
        public
        returns (uint256[] memory state)
    {
        uint256 nTeams = getNTeams(id);
    }

    function getLastChallengeBlock(uint256 id) public view returns (uint256) {
        require(isUpdated(id), "not updated league");
        return getUpdateBlock(id) + CHALLENGING_PERIOD_BLKS;
    }

    function isVerified(uint256 id) public view returns (bool) {
        if (!isUpdated(id))
            return false;
        return block.number > getLastChallengeBlock(id);
    }
}