pragma solidity ^0.5.0;

import "./LeaguesState.sol";
import "./LeaguesScore.sol";
import "./LeaguesTactics.sol";
import "./Engine.sol";

contract LeaguesComputer is LeaguesState, LeaguesScore, LeaguesTactics {
    uint8 constant PLAYERS_PER_TEAM = 11;
    Engine private _engine;

    constructor(address engine) public {
        _engine = Engine(engine);
    }

    function getEngineContract() external view returns (address) {
        return address(_engine);
    }

    // TODO: add minimum checks
    function updateLeague(
        uint256 id, 
        bytes32 initStateHash,
        bytes32[] memory finalHashes,
        uint16[] memory scores
    ) 
        public 
    {
        _setInitStateHash(id, initStateHash);
        _setFinalTeamStateHashes(id, finalHashes);
        _setScores(id, scores);
    }

    function updatePlayerStatesAfterMatch(
        uint256[] memory homeTeamState,
        uint256[] memory visitorTeamState,
        uint8 homeGoals,
        uint8 visitorGoals
    )
        public 
        pure
        returns (uint256[] memory updatedHomeTeamState, uint256[] memory updatedVisitorTeamState) 
    {
        if (homeGoals == visitorGoals)
            return (homeTeamState, visitorTeamState);

        uint8 pointsWon = computePointsWon(
            homeTeamState,
            visitorTeamState,
            homeGoals,
            visitorGoals
        );

        if (homeGoals > visitorGoals){
            updatedHomeTeamState = teamStateEvolve(homeTeamState, pointsWon);             
            updatedVisitorTeamState = visitorTeamState;
        }
        else {
            updatedHomeTeamState = homeTeamState;
            updatedVisitorTeamState = teamStateEvolve(visitorTeamState, pointsWon);
        }
    }

    function computePointsWon(
        uint256[] memory homeTeamState, 
        uint256[] memory visitorTeamState,
        uint8 homeGoals,
        uint8 visitorGoals
    )
        public 
        pure
        returns (uint8 points)
    {
        require(isValidTeamState(homeTeamState), "home team state invalid");
        require(isValidTeamState(visitorTeamState), "visitor team state invalid");
        uint128 homeTeamRating = computeTeamRating(homeTeamState);
        uint128 visitorTeamRating = computeTeamRating(visitorTeamState);
        int256 ratingDiff = homeTeamRating - visitorTeamRating;
        if (ratingDiff == 0)
            return 5;
        int256 goalsDiff = homeGoals - visitorGoals;
        bool winnerWasBetter = (ratingDiff > 0 && goalsDiff > 0) || (ratingDiff < 0 && goalsDiff < 0);
        if (winnerWasBetter)
            return 2;
        return 10;
    }

    function computeStatesAtMatchday(
        uint256 id,
        uint256 leagueDay, 
        uint256[] memory initLeagueState, 
        uint256[3][] memory tactics,
        bytes32 seed
    )
        public
        view
        returns (uint16[] memory scores)
    {
        uint256 nMatchesPerMatchday = getMatchPerDay(id);
        for (uint256 i = 0; i < nMatchesPerMatchday ; i++)
        {
            uint16 score = computeScoreMatchInLeague(id, leagueDay, i, initLeagueState, tactics, seed);
            scores = addToDayScores(scores, score);
        }
    }

    function computeScoreMatchInLeague(
        uint256 id,
        uint256 leagueDay, 
        uint256 matchInLeagueDay,
        uint256[] memory initLeagueState, 
        uint256[3][] memory tactics,
        bytes32 seed
    )
        public
        view
        returns (uint16 score)
    {
        uint256 homeTeamIdx;
        uint256 visitorTeamIdx;
        uint8 homeGoals;
        uint8 visitorGoals;
        (homeTeamIdx, visitorTeamIdx) = getTeamsInMatch(id, leagueDay, matchInLeagueDay);
        (homeGoals, visitorGoals) = _engine.playMatch(
            seed, 
            getTeam(initLeagueState, homeTeamIdx), 
            getTeam(initLeagueState, visitorTeamIdx), 
            tactics[0], 
            tactics[1]
        );
        score = encodeScore(homeGoals, visitorGoals);
    }

    function computeAllMatchdayStates(
        uint256 id, 
        uint256[] memory initLeagueState, 
        uint256[3][] memory tactics // TODO: optimize data type
    )
        public 
        view 
        returns (uint16[] memory scores) 
    {
        uint256 nLeagueDays = countLeagueDays(id);
        for(uint256 day = 0 ; day < nLeagueDays ; day++)
        {
            bytes32 seed = getMatchDayBlockHash(id, day);
            uint16[] memory dayScores = computeStatesAtMatchday(id, day, initLeagueState, tactics, seed);
            scores = addToTournamentScores(scores, dayScores);
        }
    }

    function hashLeagueState(uint256[] memory leagueState) public pure returns (bytes32[] memory) {
        uint256 nTeams = countTeamsInState(leagueState);
        bytes32[] memory hashes = new bytes32[](nTeams);
        for (uint256 i = 0; i < nTeams ; i++){
            uint256[] memory teamState = getTeam(leagueState, i);
            hashes[i] = keccak256(abi.encode(teamState));
        }
        return hashes;
    }

    function hashInitState(uint256[] memory state) public pure returns (bytes32) {
        return keccak256(abi.encode(state));
    }

    function hashTeamState(uint256[] memory state) public pure returns (bytes32) {
        return keccak256(abi.encode(state));
    }

    function hashTactics(uint256[3][] memory tactics) public pure returns (bytes32) {
        return keccak256(abi.encode(tactics));
    }
}