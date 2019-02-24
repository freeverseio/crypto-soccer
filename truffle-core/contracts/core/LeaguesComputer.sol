pragma solidity ^0.5.0;

import "./LeaguesScheduler.sol";
import "./LeagueState.sol";
import "./Scores.sol";
import "./Engine.sol";

contract LeaguesComputer is LeaguesScheduler, Scores {
    using LeagueState for uint256[];

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

    function computeStatesAtMatchday(
        uint256 id,
        uint256 matchday, 
        uint256[] memory prevStates, 
        uint256[3][] memory tactics,
        bytes32 seed
    )
        public
        view
        returns (uint16[] memory scores)
    {
        uint256 nMatchesPerMatchday = getMatchPerDay(id);
        uint256 team0Idx;
        uint256 team1Idx;
        for (uint256 i = 0; i < nMatchesPerMatchday ; i++)
        {
            uint8 home;
            uint8 visitor;
            (team0Idx, team1Idx) = getTeamsInMatch(id, matchday, i);
            (home, visitor) = _engine.playMatch(
                seed, 
                prevStates.getTeam(team0Idx), 
                prevStates.getTeam(team1Idx), 
                tactics[0], 
                tactics[1]
            );
            uint16 score = encodeScore(home, visitor);
            scores = addToDayScores(scores, score);
        }
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
            // scores = scoresConcat(scores, dayScores);
        }
    }

    function hashLeagueState(uint256[] memory leagueState) public pure returns (bytes32[] memory) {
        uint256 nTeams = leagueState.countTeams();
        bytes32[] memory hashes = new bytes32[](nTeams);
        for (uint256 i = 0; i < nTeams ; i++){
            uint256[] memory teamState = leagueState.getTeam(i);
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