pragma solidity ^0.4.25;

import "./Leagues.sol";
import "./Engine.sol";

contract LeaguesComputer is Leagues {
    uint8 constant PLAYERS_PER_TEAM = 11;
    Engine private _engine;

    constructor(address engine) public {
        _engine = Engine(engine);
    }

    function getEngineContract() external view returns (address) {
        return address(_engine);
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
        uint256[] memory playersState, // TODO: if teams are fixed size it can be optimized
        uint256[3][] memory tactics // TODO: optimize data type
    )
        public 
        view 
        returns (uint256[2][] memory) 
    {
        uint256 nTeams = countTeams(leagueId);
        require(playersState.length == nTeams * PLAYERS_PER_TEAM, "wrong number of players");
        require(tactics.length == nTeams, "nTeams and size of tactics mismatch");

        uint256 i;
        uint256[][] memory state = new uint256[][](nTeams);
        for (i = 0; i < nTeams; i++){
            state[i] = new uint256[](PLAYERS_PER_TEAM);
            for (uint256 j = 0; j < PLAYERS_PER_TEAM; j++){
                state[i][j] = playersState[i*PLAYERS_PER_TEAM + j];
            }
        }

        uint256 nMatches = nTeams * (nTeams - 1);
        uint256[2][] memory scores = new uint256[2][](nMatches); 
        uint256 leagueInitBlock = getInitBlock(leagueId);
        bytes32 seed = blockhash(leagueInitBlock);
        require(seed != 0, "can't retrive league init block hash");
        for (i = 0; i < nMatches; i++)
            (scores[i][0], scores[i][1]) = _engine.playMatch(seed, state[0], state[1], tactics[0], tactics[1]);

        return scores;
    }

    // TODO: function name => hashResult ?
    function calculateFinalHash(uint256[2][] memory scores) public pure returns (bytes32) {
        bytes memory origin;
        for(uint256 i = 0; i < scores.length ; i++){
            origin = abi.encodePacked(origin, scores[i][0]); 
            origin = abi.encodePacked(origin, scores[i][1]); 
        }
        return keccak256(origin);
    }

    function hashTactics(uint256[3][] memory tactics) public pure returns (bytes32) {
        bytes memory origin;
        for(uint256 i = 0; i < tactics.length ; i++){
            origin = abi.encodePacked(origin, tactics[i][0]); 
            origin = abi.encodePacked(origin, tactics[i][1]); 
            origin = abi.encodePacked(origin, tactics[i][2]); 
        }
        return keccak256(origin);
    }

    function updateLeague(
        uint256 id, 
        bytes32[] memory finalHashes,
        uint256[2][] memory scores
    ) 
        public 
    {
        _setFinalHashes(id, finalHashes);
        _setScores(id, scores);
    }
}