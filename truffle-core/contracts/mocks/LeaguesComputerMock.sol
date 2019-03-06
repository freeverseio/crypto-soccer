pragma solidity ^0.5.0;

import "../core/LeaguesComputer.sol";

contract LeaguesComputerMock is LeaguesComputer {
    constructor(address engine, address state) 
    public 
    LeaguesComputer(engine, state)
    {
    }    

    function updatePlayerStatesAfterMatch(
        uint256[] memory homeTeamState,
        uint256[] memory visitorTeamState,
        uint8 homeGoals,
        uint8 visitorGoals
    )
        public
        view
        returns (uint256[] memory updatedHomeTeamState, uint256[] memory updatedVisitorTeamState) 
    {
        return _updatePlayerStatesAfterMatch(homeTeamState, visitorTeamState, homeGoals, visitorGoals);
    }

    function computePointsWon(
        uint256[] memory homeTeamState, 
        uint256[] memory visitorTeamState,
        uint8 homeGoals,
        uint8 visitorGoals
    )
        public 
        view
        returns (uint8 points)
    {
        return _computePointsWon(homeTeamState, visitorTeamState, homeGoals, visitorGoals);
    }

    function computeStatesAtMatchday(
        uint256 id,
        uint256 leagueDay, 
        uint256[] memory initDayState, 
        uint256[3][] memory tactics,
        bytes32 seed
    )
        public
        view
        returns (uint16[] memory scores, uint256[] memory finalDayState)
    {
        return _computeStatesAtMatchday(id, leagueDay, initDayState, tactics, seed);
    }

    function computeScoreMatchInLeague(
        uint256[] memory homeTeamState,
        uint256[] memory visitorTeamState,
        uint256[3][] memory tactics,
        bytes32 seed
    )
        public
        view
        returns (uint16 score, uint256[] memory newHomeState, uint256[] memory newVisitorState)
    {
        return _computeScoreMatchInLeague(homeTeamState, visitorTeamState, tactics, seed);
    }

}