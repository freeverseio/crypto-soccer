pragma solidity ^0.4.25;

import "./LeaguesScheduler.sol";
import "./LeagueState.sol";
import "./Engine.sol";

contract LeaguesComputer is LeaguesScheduler {
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
        uint256[2][] memory scores
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
        returns (uint256[2][] memory scores)
    {
        uint256 nMatchesPerMatchday = getMatchPerDay(id);
        scores = new uint256[2][](nMatchesPerMatchday);
        uint256 team0Idx;
        uint256 team1Idx;
        for (uint256 i = 0; i < nMatchesPerMatchday ; i++)
        {
            (team0Idx, team1Idx) = getTeamsInMatch(id, matchday, i);
            (scores[i][0], scores[i][1]) = _engine.playMatch(
                seed, 
                prevStates.getTeam(team0Idx), 
                prevStates.getTeam(team1Idx), 
                tactics[0], 
                tactics[1]
            );
        }
    }


    /**
     * @dev compute the result of a league
     * @param leagueId id of the league to compute
     * @return result of every match
    */
    // TODO: better function name ex. computeLeague ?
    // TODO: does a team play always with 11 players ?
    function computeLeagueFinalState (
        uint256 leagueId,
        uint256[] memory leagueState, 
        uint256[3][] memory tactics // TODO: optimize data type
    )
        public 
        view 
        returns (uint256[2][] memory scores) 
    {
        uint256 nTeams = countTeams(leagueId);
        require(leagueState.countTeams() == nTeams, "wrong number of teams");
        require(tactics.length == nTeams, "nTeams and size of tactics mismatch");

        uint256[][] memory state = new uint256[][](nTeams);
        for (uint256 i = 0; i < nTeams; i++){
            state[i] = new uint256[](leagueState.countTeamPlayers(i));
            uint256[] memory teamState = leagueState.getTeam(i);
            for (uint256 j = 0; j < teamState.length ; j++)
                state[i][j] = teamState[j];
        }

        uint256 nMatches = nTeams * (nTeams - 1);
        scores = new uint256[2][](nMatches); 
        uint256 leagueInitBlock = getInitBlock(leagueId);
        bytes32 seed = blockhash(leagueInitBlock);
        require(seed != 0, "can't retrive league init block hash");
        for (i = 0; i < nMatches; i++)
            (scores[i][0], scores[i][1]) = _engine.playMatch(seed, state[0], state[1], tactics[0], tactics[1]);
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