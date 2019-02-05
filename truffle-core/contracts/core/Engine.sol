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

    function computeLeagueFinalState (
        uint256 leagueId,
        uint256[NPLAYERS_PER_TEAM][] memory initPlayerState
        ) public view returns (uint256[2][] memory) {
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

    function playMatch(uint256[] memory stateTeam0, uint256[] memory stateTeam1, bytes32 seed) public pure returns (uint256, uint256) {
        uint256 hash1 = uint256(keccak256(stateTeam0));
        uint256 hash2 = uint256(keccak256(stateTeam1));
        return (hash1 % 4, hash2 % 4);
    }
}