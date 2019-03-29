pragma solidity ^0.5.0;

import "../state/LeagueState.sol";
import "./LeaguesScore.sol";
import "./Engine.sol";

contract LeaguesComputer is LeaguesScore {
    LeagueState private _leagueState;
    Engine private _engine;

    uint8 constant PLAYERS_PER_TEAM = 11;

    constructor(address engine, address leagueState) public {
        _engine = Engine(engine);
        _leagueState = LeagueState(leagueState);
    }

    function getEngineContract() external view returns (address) {
        return address(_engine);
    }

    function computeDay(
        uint256 leagueId,
        uint256 leagueDay, 
        uint256[] memory initLeagueState, 
        uint8[3][] memory tactics
    )
        public
        view
        returns (uint16[] memory scores, uint256[] memory finalLeagueState)
    {
        bytes32 seed = getMatchDayBlockHash(leagueId, leagueDay);
        return _computeDayWithSeed(
            leagueId,
            leagueDay,
            initLeagueState,
            tactics,
            seed
        );
    } 

    /// @dev evolves the teams of a match in front of the scores
    /// @param homeTeamState initial state of the home team
    /// @param visitorTeamState initial state of the visitor team
    /// @param homeGoals goals scored by home team
    /// @param visitorTeamState goals scored by visitor team
    /// @return updated home team state and updated visitor team state
    function _evolveTeams(
        uint256[] memory homeTeamState,
        uint256[] memory visitorTeamState,
        uint8 homeGoals,
        uint8 visitorGoals
    )
        internal
        view
        returns (uint256[] memory updatedHomeTeamState, uint256[] memory updatedVisitorTeamState) 
    {
        (uint8 homeTeamPoints, uint8 visitorTeamPoints) = _computePoints(
            homeTeamState,
            visitorTeamState,
            homeGoals,
            visitorGoals
        );

        updatedHomeTeamState = _leagueState.teamStateEvolve(homeTeamState, homeTeamPoints);             
        updatedVisitorTeamState = _leagueState.teamStateEvolve(visitorTeamState, visitorTeamPoints);
    }

    /// compute points per team in front of goals
    /// @return home and visitor points
    function _computePoints(
        uint256[] memory homeTeamState, 
        uint256[] memory visitorTeamState,
        uint8 homeGoals,
        uint8 visitorGoals
    )
        internal
        view
        returns (uint8 homePoints, uint8 visitorPoints)
    {
        if (homeGoals == visitorGoals)
            return (0, 0);

        uint256 homeTeamRating = _leagueState.computeTeamRating(homeTeamState);
        uint256 visitorTeamRating = _leagueState.computeTeamRating(visitorTeamState);

        if (homeTeamRating == visitorTeamRating)
            return homeGoals > visitorGoals ? (5, 0) : (0, 5);
        else if (homeTeamRating > visitorTeamRating)
            return homeGoals > visitorGoals ? (2, 0) : (0, 8);
        else 
            return homeGoals > visitorGoals ? (8, 0) : (0, 2);
    }

    function _computeDayWithSeed(
        uint256 id,
        uint256 leagueDay, 
        uint256[] memory initLeagueState, 
        uint8[3][] memory tactics,
        bytes32 seed
    )
        internal
        view
        returns (uint16[] memory scores, uint256[] memory finalLeagueState)
    {
        uint256 nMatchesPerMatchday = getMatchPerDay(id);
        finalLeagueState = initLeagueState; 
        for (uint256 i = 0; i < nMatchesPerMatchday ; i++)
        {
            (uint256 homeTeamIdx, uint256 visitorTeamIdx) = getTeamsInMatch(id, leagueDay, i);
            uint256[] memory homeTeamState = _leagueState.leagueStateAt(initLeagueState, homeTeamIdx);
            uint256[] memory visitorTeamState = _leagueState.leagueStateAt(initLeagueState, visitorTeamIdx);
            (uint16 score,
            uint256[] memory updatedHomeState,
            uint256[] memory updatedVisitorState) = _computeMatch(
                homeTeamState,
                visitorTeamState, 
                tactics, 
                seed
            );
            scores = scoresAppend(scores, score);
            finalLeagueState = _leagueState.leagueStateUpdate(finalLeagueState, homeTeamIdx, updatedHomeState);
            finalLeagueState = _leagueState.leagueStateUpdate(finalLeagueState, visitorTeamIdx, updatedVisitorState);
        }
    }

    /// compute score of a match and evolve home and visitor team states
    function _computeMatch(
        uint256[] memory homeTeamState,
        uint256[] memory visitorTeamState,
        uint8[3][] memory tactics,
        bytes32 seed
    )
        internal
        view
        returns (uint16 score, uint256[] memory newHomeState, uint256[] memory newVisitorState)
    {
        (uint8 homeGoals, uint8 visitorGoals) = _engine.playMatch(
            seed, 
            homeTeamState, 
            visitorTeamState, 
            tactics[0], 
            tactics[1]
        );
        score = encodeScore(homeGoals, visitorGoals);
        (newHomeState, newVisitorState) = _evolveTeams(
            homeTeamState,
            visitorTeamState,
            homeGoals,
            visitorGoals
        );
    }
        
    // function computeAllMatchleagueStates(
    //     uint256 id, 
    //     uint256[] memory initleagueState, 
    //     uint256[3][] memory tactics // TODO: optimize data type
    // )
    //     public 
    //     view 
    //     returns (uint16[] memory scores) 
    // {
    //     uint256 nLeagueDays = countLeagueDays(id);
    //     for(uint256 day = 0 ; day < nLeagueDays ; day++)
    //     {
    //         bytes32 seed = getMatchDayBlockHash(id, day);
    //         // uint16[] memory dayScores = computeStatesAtMatchday(id, day, initleagueState, tactics, seed);
    //         // scores = addToTournamentScores(scores, dayScores);
    //     }
    // }
}