pragma solidity ^0.5.0;

import "../core/LeaguesComputer.sol";

contract LeaguesComputerMock is LeaguesComputer {
    constructor(address engine, address state) 
    public 
    LeaguesComputer(engine, state)
    {
    }    

    function evolveTeams(
        uint256[] memory homeTeamState,
        uint256[] memory visitorTeamState,
        uint8 homeGoals,
        uint8 visitorGoals
    )
        public
        view
        returns (uint256[] memory updatedHomeTeamState, uint256[] memory updatedVisitorTeamState) 
    {
        return _evolveTeams(homeTeamState, visitorTeamState, homeGoals, visitorGoals);
    }

    function computePoints(
        uint256[] memory homeTeamState, 
        uint256[] memory visitorTeamState,
        uint8 homeGoals,
        uint8 visitorGoals
    )
        public 
        view
        returns (uint8 homePoints, uint8 visitorPoints)
    {
        return _computePoints(homeTeamState, visitorTeamState, homeGoals, visitorGoals);
    }

    function computeDayWithSeed(
        uint256 id,
        uint256 leagueDay, 
        uint256[] memory initleagueState, 
        uint8[3][] memory tactics,
        uint256 seed
    )
        public
        view
        returns (uint16[] memory scores, uint256[] memory finalleagueState)
    {
        return _computeDayWithSeed(id, leagueDay, initleagueState, tactics, seed);
    }

    function computeMatch(
        uint256[] memory homeTeamState,
        uint256[] memory visitorTeamState,
        uint8[3][] memory tactics,
        uint256 seed
    )
        public
        view
        returns (uint16 score, uint256[] memory newHomeState, uint256[] memory newVisitorState)
    {
        return _computeMatch(homeTeamState, visitorTeamState, tactics, seed);
    }

}