pragma solidity ^ 0.4.24;

import "./Leagues.sol";

contract Engine {
    Leagues private _leagues;

    constructor(address leagues) public {
        _leagues = Leagues(leagues);
    }

    function getLeaguesContract() external view returns (address) {
        return address(_leagues);
    }

    function computeLeagueFinalState (
        uint256 leagueId,
        uint256 initPlayerState
        ) public view {
            uint256 initBlock = _leagues.getInitBlock(leagueId);
            uint256 step = _leagues.getStep(leagueId);
            uint256[] memory teamIds = _leagues.getTeamIds(leagueId);
            uint256 nTeams = teamIds.length;
            uint256 nMatchdays = 2*(nTeams-1);
            uint256 nMatchesPerMatchday = nTeams/2;
            //uint256[] memory scores = new uint256[](nMatchdays);
    }
}