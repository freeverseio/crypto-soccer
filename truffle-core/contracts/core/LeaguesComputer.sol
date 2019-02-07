pragma solidity ^0.4.25;

import "./Leagues.sol";
import "./Engine.sol";

contract LeaguesComputer is Leagues {
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
    function computeLeagueFinalState (
        uint256 leagueId,
        uint256[] memory playersState,
        uint256[] memory playersPerTeam,
        uint256[3][] memory tactics
    )
        public 
        view 
        returns (uint256[2][] memory) 
    {
        uint256 nTeams = countTeams(leagueId);
        require(playersPerTeam.length == nTeams, "nTeams and size of playersPerTeam mismatch");
        require(tactics.length == nTeams, "nTeams and size of tactics mismatch");

        uint256 countPlayers = 0;
        for (uint256 i = 0; i < nTeams; i++){
            countPlayers += playersPerTeam[i];
        }
        require(playersState.length == countPlayers, "wrong number of players state");

        // uint256 initBlock = getInitBlock(leagueId);
        // uint256 step = getStep(leagueId);


        // uint256 nMatchdays = 2*(nTeams-1);
        // uint256 nMatchesPerMatchday = nTeams/2;
        uint256 nMatches = nTeams * (nTeams - 1);
        uint256[2][] memory scores = new uint256[2][](nMatches); 

        uint256[][] memory state = new uint256[][](nTeams);
        for (i = 0; i < nTeams; i++){
            uint256 nPlayers = playersPerTeam[i];
            state[i] = new uint256[](nPlayers);
            for (uint256 j = 0; j < nPlayers; j++){
                state[i][j] = 0;
            }
        }

        (scores[0][0], scores[0][1]) = _engine.playMatch(4353, state[0], state[1], tactics[0], tactics[1]);

        return scores;
    }

    function getTeamsInMatch(uint256 matchday, uint256 matchNumber, uint256 nTeams) private pure returns(uint256, uint256) {

    }

    function getTeamsInMatchFirstHalf(uint256 matchday, uint256 matchnumber, uint256 nTeams) private pure returns(uint256, uint256) {

    }

    function shiftBack(uint256 t, uint256 nTeams) private pure returns (uint256) {
        if (t < nTeams)
            return t;

        return t - (nTeams - 1);
    }
}