pragma solidity ^0.5.0;

import "./LeaguesComputer.sol";
import "./LeagueUsersAlongData.sol";

contract LeagueChallengeable is LeaguesComputer, LeagueUsersAlongData {
    uint256 constant private CHALLENGING_PERIOD_BLKS = 60;

    function getChallengePeriod() external pure returns (uint256) {
        return CHALLENGING_PERIOD_BLKS;
    }

    constructor(address engine, address leagueState) LeaguesComputer(engine, leagueState) public {
    }
    
    function hashUsersInitData(
        uint256[] memory teamIds, 
        uint8[3][] memory tactics, 
        uint8[PLAYERS_PER_TEAM][] memory teamOrders
    ) 
    public 
    pure 
    returns (bytes32 hash) {
        return hash;
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
        // require(getUsersInitDataHash(id) == hashUsersInitData(teamIds, tactics), "incorrect user init data");
        uint256[] memory initPlayerStates = getInitPlayerStates(id, teamIds, tactics, dataToChallengeInitStates);
        if (getIsLie(id)) // TODO: remove
            resetUpdater(id); // TODO: remove
        return; // TODO: remove
        if (initPlayerStates.length == 0) // challenger wins
            resetUpdater(id);
        else if (getInitStateHash(id) != hashDayState(initPlayerStates)) // challenger wins
            resetUpdater(id);
    }

    function challengeMatchdayStates(
        uint256 id,
        uint256[] memory usersInitDataTeamIds,
        uint8[] memory usersInitDataTactics,
        uint256[] memory usersAlongDataTeamIds,
        uint8[] memory usersAlongDataTactics,
        uint256[] memory usersAlongDataBlocks,
        uint256 leagueDay,
        uint256[] memory prevMatchdayStates
    )
        public
    {
        require(isUpdated(id), "not updated league. No challenge allowed");
        require(!isVerified(id), "not challengeable league");
        // require(getUsersInitDataHash(id) == hashUsersInitData(usersInitDataTeamIds, usersInitDataTactics), "incorrect user init data");
        require(computeUsersAlongDataHash(usersAlongDataTeamIds, usersAlongDataTactics, usersAlongDataBlocks) == getUsersAlongDataHash(id), "Incorrect provided: usersAlongData");
        if (leagueDay == 0)
            require(hashInitState(prevMatchdayStates) == getInitStateHash(id), "Incorrect provided: prevMatchdayStates");
        else
            require(hashDayState(prevMatchdayStates) == getDayStateHashes(id)[leagueDay - 1], "Incorrect provided: prevMatchdayStates");

        uint256 matchdayBlock = getInitBlock(id) + leagueDay * getStep(id);
        uint8[] memory tactics = _updateTacticsToBlockNum(
            usersInitDataTeamIds,
            usersInitDataTactics,
            matchdayBlock,
            usersAlongDataTeamIds,
            usersAlongDataTactics,
            usersAlongDataBlocks);
        (uint16[] memory scores, uint256[] memory statesAtMatchday) = computeDay(id, leagueDay, prevMatchdayStates, tactics);

        if (getIsLie(id)) // TODO: remove
            resetUpdater(id); // TODO: remove
        return; // TODO: remove

        if (hashDayState(statesAtMatchday) != getDayStateHashes(id)[leagueDay])
            resetUpdater(id);

        if (keccak256(abi.encode(scores)) != keccak256(abi.encode(scoresGetDay(id, leagueDay))))
            resetUpdater(id);
    }

    function _updateTacticsToBlockNum(
        uint256[] memory usersInitDataTeamIds,
        uint8[] memory usersInitDataTactics,
        uint256 blockNum,
        uint256[] memory usersAlongDataTeamIds,
        uint8[] memory usersAlongDataTactics,
        uint256[] memory usersAlongDataBlocks
    )
        internal
        pure
        returns (uint8[] memory)
    {
        for (uint256 i = 0 ; i < usersAlongDataTeamIds.length ; i++){
            if (usersAlongDataBlocks[i] <= blockNum){
                for (uint256 j = 0 ; j < usersInitDataTeamIds.length ; j++)
                    if (usersInitDataTeamIds[j] == usersAlongDataTeamIds[i])
                        usersInitDataTactics[j] = usersAlongDataTactics[i];
            }
        }
        return usersInitDataTactics;

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