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
        uint256[] memory playersPerTeam
    )
        public 
        view 
        returns (uint256[2][] memory) 
    {
        uint256 nTeams = countTeams(leagueId);
        require(nTeams == playersPerTeam.length, "nTeams and size of playersPerTeam mismatch");

        // uint256 initBlock = getInitBlock(leagueId);
        // uint256 step = getStep(leagueId);


        // uint256 nMatchdays = 2*(nTeams-1);
        // uint256 nMatchesPerMatchday = nTeams/2;
        uint256 nMatches = nTeams * (nTeams - 1);
        uint256[2][] memory scores = new uint256[2][](nMatches); 

        uint256[] memory team0;
        uint256[] memory team1;
        uint256[3] memory tactic0;
        uint256[3] memory tactic1;
        _engine.playMatch(4353, team0, team1, tactic0, tactic1);

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