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
        returns (uint256[2][] memory) 
    {
        uint256 nTeams = countTeams(leagueId);
        require(leagueState.countTeams() == nTeams, "wrong number of teams");
        require(tactics.length == nTeams, "nTeams and size of tactics mismatch");

        uint256[][] memory state = new uint256[][](nTeams);
        state[0] = new uint256[](11);
        uint256 team;
        for (uint256 i = 0; i < leagueState.length - 1; i++){
            if(leagueState[i] == 0){
                team++;
                state[team] = new uint256[](11);
            } else {
                // state[team].push(leagueState[i]);
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

    function hashLeagueState(uint256[] memory leagueState) public pure returns (bytes32[] memory) {
        uint256 nTeams = leagueState.countTeams();
        bytes32[] memory hashes = new bytes32[](nTeams);
        for (uint256 i = 0; i < nTeams ; i++){
            uint256[] memory teamState = leagueState.getTeamState(i);
            hashes[i] = keccak256(abi.encode(teamState));
        }
        return hashes;
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

    function hashInitState(uint256[] memory state) public pure returns (bytes32) {
        return _hashState(state);
    }

    function hashTeamState(uint256[] memory state) public pure returns (bytes32) {
        return _hashState(state);
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

    function _hashState(uint256[] memory state) private pure returns (bytes32) {
        bytes memory origin;
        for(uint256 i = 0; i < state.length ; i++){
            origin = abi.encodePacked(origin, state[i]); 
        }
        return keccak256(origin);
    }
}