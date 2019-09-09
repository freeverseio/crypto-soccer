pragma solidity >=0.4.21 <0.6.0;

import "./Assets.sol";
import "./Engine.sol";

contract Leagues is Assets {
    
    uint8 constant public MATCHDAYS = 14;
    uint8 constant public MATCHES_PER_DAY = 4;
    Engine private _engine;

    function setEngineAdress(address addr) public {
        _engine = Engine(addr);
    }

    function getEngineAddress() public view returns (address) {
        return address(_engine);
    }

    function getTeamsInMatch(uint8 matchday, uint8 matchIdxInDay) public pure returns (uint8 homeIdx, uint8 visitorIdx) 
    {
        require(matchday < MATCHDAYS, "wrong match day");
        require(matchIdxInDay < MATCHES_PER_DAY, "wrong match");
        if (matchday < (TEAMS_PER_LEAGUE - 1))
            (homeIdx, visitorIdx) = _getTeamsInMatchFirstHalf(matchday, matchIdxInDay);
        else
            (visitorIdx, homeIdx) = _getTeamsInMatchFirstHalf(matchday - (TEAMS_PER_LEAGUE - 1), matchIdxInDay);
    }

    function _shiftBack(uint8 t) private pure returns (uint8)
    {
        if (t < TEAMS_PER_LEAGUE)
            return t;
        else
            return t-(TEAMS_PER_LEAGUE-1);
    }

    function _getTeamsInMatchFirstHalf(uint8 matchday, uint8 matchIdxInDay) private pure returns (uint8, uint8) 
    {
        uint8 team1 = 0;
        if (matchIdxInDay > 0)
            team1 = _shiftBack(TEAMS_PER_LEAGUE-matchIdxInDay+matchday);

        uint8 team2 = _shiftBack(matchIdxInDay+1+matchday);
        if ( (matchday % 2) == 0)
            return (team1, team2);
        else
            return (team2, team1);
    }

    /// compute points per team in front of goals
    /// @return home and visitor points
    function computeEvolutionPoints(
        uint256[] memory homeTeamState, 
        uint256[] memory visitorTeamState,
        uint8 homeGoals,
        uint8 visitorGoals
    )
        public
        pure
        returns (uint8 homePoints, uint8 visitorPoints)
    {
        if (homeGoals == visitorGoals)
            return (0, 0);

        uint256 homeTeamRating = computeTeamRating(homeTeamState);
        uint256 visitorTeamRating = computeTeamRating(visitorTeamState);

        if (homeTeamRating == visitorTeamRating)
            return homeGoals > visitorGoals ? (5, 0) : (0, 5);
        else if (homeTeamRating > visitorTeamRating)
            return homeGoals > visitorGoals ? (2, 0) : (0, 8);
        else 
            return homeGoals > visitorGoals ? (8, 0) : (0, 2);
    }

    function computeTeamRating(uint256[] memory teamState) public pure returns (uint256 rating) {
        for(uint256 i = 0 ; i < teamState.length ; i++){
            uint256 playerSkills = teamState[i];
            if (getPlayerIdFromSkills(playerSkills) != FREE_PLAYER_ID) {
                uint16[5] memory skills = getSkillsVec(playerSkills);
                for (uint8 sk = 0; sk < N_SKILLS; sk++) {
                    rating += skills[sk];
                }
            }
        }
    }


}