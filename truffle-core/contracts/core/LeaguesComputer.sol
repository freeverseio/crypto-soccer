pragma solidity ^0.5.0;

import "../state/DayState.sol";
import "./LeaguesScore.sol";
import "./Engine.sol";

contract LeaguesComputer is LeaguesScore {
    DayState private _leagueState;
    Engine private _engine;

    uint8 constant PLAYERS_PER_TEAM = 11;

    constructor(address engine, address leagueState) public {
        _engine = Engine(engine);
        _leagueState = DayState(leagueState);
    }

    function getEngineContract() external view returns (address) {
        return address(_engine);
    }

    function computeDay(
        uint256 leagueId,
        uint256 leagueDay, 
        uint256[] memory initDayState, 
        uint256[3][] memory tactics
    )
        public
        view
        returns (uint16[] memory scores, uint256[] memory finalDayState)
    {
        bytes32 seed = getMatchDayBlockHash(leagueId, leagueDay);
        return _computeStatesAtMatchday(
            leagueId,
            leagueDay,
            initDayState,
            tactics,
            seed
        );
    } 

    function _updatePlayerStatesAfterMatch(
        uint256[] memory homeTeamState,
        uint256[] memory visitorTeamState,
        uint8 homeGoals,
        uint8 visitorGoals
    )
        internal
        view
        returns (uint256[] memory updatedHomeTeamState, uint256[] memory updatedVisitorTeamState) 
    {
        if (homeGoals == visitorGoals)
            return (homeTeamState, visitorTeamState);

        uint8 pointsWon = _computePointsWon(
            homeTeamState,
            visitorTeamState,
            homeGoals,
            visitorGoals
        );

        if (homeGoals > visitorGoals){
            updatedHomeTeamState = _leagueState.teamStateEvolve(homeTeamState, pointsWon);             
            updatedVisitorTeamState = visitorTeamState;
        }
        else {
            updatedHomeTeamState = homeTeamState;
            updatedVisitorTeamState = _leagueState.teamStateEvolve(visitorTeamState, pointsWon);
        }
    }

    function _computePointsWon(
        uint256[] memory homeTeamState, 
        uint256[] memory visitorTeamState,
        uint8 homeGoals,
        uint8 visitorGoals
    )
        internal
        view
        returns (uint8 points)
    {
        require(_leagueState.isValidTeamState(homeTeamState), "home team state invalid");
        require(_leagueState.isValidTeamState(visitorTeamState), "visitor team state invalid");
        uint128 homeTeamRating = _leagueState.computeTeamRating(homeTeamState);
        uint128 visitorTeamRating = _leagueState.computeTeamRating(visitorTeamState);
        int256 ratingDiff = homeTeamRating - visitorTeamRating;
        if (ratingDiff == 0)
            return 5;
        int256 goalsDiff = homeGoals - visitorGoals;
        bool winnerWasBetter = (ratingDiff > 0 && goalsDiff > 0) || (ratingDiff < 0 && goalsDiff < 0);
        if (winnerWasBetter)
            return 2;
        return 10;
    }

    function _computeStatesAtMatchday(
        uint256 id,
        uint256 leagueDay, 
        uint256[] memory initDayState, 
        uint256[3][] memory tactics,
        bytes32 seed
    )
        internal
        view
        returns (uint16[] memory scores, uint256[] memory finalDayState)
    {
        uint256 nMatchesPerMatchday = getMatchPerDay(id);
        for (uint256 i = 0; i < nMatchesPerMatchday ; i++)
        {
            (uint256 homeTeamIdx, uint256 visitorTeamIdx) = getTeamsInMatch(id, leagueDay, i);
            uint256[] memory homeTeamState = _leagueState.dayStateAt(initDayState, homeTeamIdx);
            uint256[] memory visitorTeamState = _leagueState.dayStateAt(initDayState, visitorTeamIdx);
            (uint16 score,
            uint256[] memory updatedHomeState,
            uint256[] memory updatedVisitorState) = _computeScoreMatchInLeague(
                homeTeamState,
                visitorTeamState, 
                tactics, 
                seed
            );
            scores = addToDayScores(scores, score);
            finalDayState = _leagueState.dayStateUpdate(initDayState, homeTeamIdx, updatedHomeState);
            finalDayState = _leagueState.dayStateUpdate(initDayState, visitorTeamIdx, updatedVisitorState);
        }
    }

    function _computeScoreMatchInLeague(
        uint256[] memory homeTeamState,
        uint256[] memory visitorTeamState,
        uint256[3][] memory tactics,
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
        (newHomeState, newVisitorState) = _updatePlayerStatesAfterMatch(
            homeTeamState,
            visitorTeamState,
            homeGoals,
            visitorGoals
        );
    }
        
    // function computeAllMatchdayStates(
    //     uint256 id, 
    //     uint256[] memory initDayState, 
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
    //         // uint16[] memory dayScores = computeStatesAtMatchday(id, day, initDayState, tactics, seed);
    //         // scores = addToTournamentScores(scores, dayScores);
    //     }
    // }
}