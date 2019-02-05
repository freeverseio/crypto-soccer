pragma solidity ^ 0.4.24;

import "./Leagues.sol";

contract Engine {
    uint8 constant NPLAYERS_PER_TEAM = 16;

    Leagues private _leagues;

    constructor(address leagues) public {
        _leagues = Leagues(leagues);
    }

    function getLeaguesContract() external view returns (address) {
        return address(_leagues);
    }
    
    /**
     * @dev playMatch returns the result of a match
     * @param seed the pseudo-random number to use as a seed for the match
     * @param stateTeam0 a vector with the state of the players of team 0
     * @param stateTeam1 a vector with the state of the players of team 1
     * @param tacticsTeam0 a vector[3] with the tactic (ex. [4,4,3]) of team 0 
     * @param tacticsTeam0 a vector[3] with the tactic (ex. [4,4,3]) of team 1
     * @return the score of the match
     */
    function playMatch(
        bytes32 seed,
        uint256[NPLAYERS_PER_TEAM] memory stateTeam0,
        uint256[NPLAYERS_PER_TEAM] memory stateTeam1, 
        uint256[3] memory tacticsTeam0, 
        uint256[3] memory tacticsTeam1
    ) 
        public 
        pure 
        returns (uint256, uint256) 
    {
        uint256 hash0 = uint256(seed) + stateTeam0[0];
        uint256 hash1 = uint256(seed) + stateTeam1[0];
        return (hash0 % 4, hash1 % 4);
    }

    /**
     * @dev compute the result of a league
     * @param leagueId id of the league to compute
     * @param initPlayerState initial state of the players of the league
     * @return result of every match
    */
    function computeLeagueFinalState (
        uint256 leagueId,
        uint256[NPLAYERS_PER_TEAM][] memory initPlayerState
    )
        public 
        view 
        returns (uint256[2][] memory) 
    {
            uint256 initBlock = _leagues.getInitBlock(leagueId);
            uint256 step = _leagues.getStep(leagueId);
            uint256[] memory teamIds = _leagues.getTeamIds(leagueId);
            uint256 nTeams = teamIds.length;
            uint256 nMatchdays = 2*(nTeams-1);
            uint256 nMatchesPerMatchday = nTeams/2;
            uint256[2][] memory scores; // TODO
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