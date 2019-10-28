pragma solidity >=0.4.21 <0.6.0;

import "./Engine.sol";
import "./SortIdxs.sol";
import "./EncodingSkills.sol";
/**
 * @title Scheduling of leagues, and calls to Engine to resolve games.
 */

contract Championships is SortIdxs, EncodingSkills {
    
    uint8 constant public PLAYERS_PER_TEAM_MAX  = 25;
    uint8 constant public TEAMS_PER_LEAGUE = 8;
    uint8 constant public MATCHDAYS = 14;
    uint8 constant public MATCHES_PER_DAY = 4;
    uint8 constant public MATCHES_PER_LEAGUE = 56; // = 4 * 14 = 7*8
    uint256 constant public FREE_PLAYER_ID  = 1; // it never corresponds to a legit playerId due to its TZ = 0
    uint256 constant private INERTIA = 4;
    uint256 constant private ONE_MINUS_INERTIA = 6;
    uint256 constant private WEIGHT_SKILLS = 100;
    uint256 constant private SKILLS_AT_START = 900; // 18 players per team at start with 50 avg

    Engine private _engine;

    function setEngineAdress(address addr) public {
        _engine = Engine(addr);
    }

    function getEngineAddress() public view returns (address) {
        return address(_engine);
    }

    // groupIdx = 0,...,15
    // posInGroup = 0, ...7
    // teamIdx  = 0,...,128
    function getTeamIdxInCup(uint8 groupIdx, uint8 posInGroup) public pure returns(uint8) {
        if (groupIdx % 2 == 0) {
                return 8 * posInGroup + groupIdx / 2;
        } else {
                return 8 * posInGroup + (groupIdx - 1) / 2 + 64;
        }
    }

    function getGroupAndPosInGroup(uint8 teamIdxInCup) public pure returns(uint8 groupIdx, uint8 posInGroup) {
        if (teamIdxInCup < 64) {
            posInGroup = teamIdxInCup / 8;
            groupIdx = 2 * (teamIdxInCup % 8);
        } else {
            posInGroup = (teamIdxInCup-64) / 8;
            groupIdx = 1 + 2 * (teamIdxInCup % 8);
        }
    }

    // sortedTeamIdxInCup contains 64 teams, made up from the top 4 in each of the 16 leagues.
    // they are flattened by groupIdx, and then by final classifications in that group.
    // [groupIdx0_1st, ..., groupIdx0_4th; groupIdx1_1st, ...,]
    // so the index of the array is  groupIdx * 4 + classInGroup
    //      M(2m) 	= L(m mod M, 0) vs L(m+1 mod M, 3), 	m = 0,..., M-1,  M = 16
    //      M(2m+1) = L(m+2 mod M, 1) vs L(m+3 mod M, 2),	m = 0,..., M-1,  M = 16
    // returns indices in sortedTeamIdxInCup
    function getTeamsInCupPlayoffMatch(uint8 matchIdxInDay) public pure returns (uint8 team0, uint8 team1) {
        require(matchIdxInDay < 32, "there are only 32 mathches on day 9 of a cup");
        if (matchIdxInDay % 2 == 0) {
            team0 = (matchIdxInDay/2 % 16) * 4;
            team1 = ((matchIdxInDay/2 + 1) % 16) * 4 + 3;
        } else {
            team0 = (((matchIdxInDay-1)/2 + 2) % 16) * 4 + 1;
            team1 = (((matchIdxInDay-1)/2 + 3) % 16) * 4 + 2;
        }
    }
    
    // same as above, but returns teamIdxInCup as correctly provided by sortedTeamIdxInCup
    function getTeamsInCupPlayoffMatch(uint8 matchIdxInDay, uint8[64] memory sortedTeamIdxInCup) public pure returns (uint8 team0, uint8 team1) {
        (team0, team1) = getTeamsInCupPlayoffMatch(matchIdxInDay);
        return (sortedTeamIdxInCup[team0], sortedTeamIdxInCup[team1]);
    }



    function getTeamsInCupLeagueMatch(uint8 groupIdx, uint8 matchday, uint8 matchIdxInDay) public pure returns (uint8, uint8) 
    {
        require(matchday < MATCHDAYS/2, "wrong match day");
        (uint8 homeIdx, uint8 visitorIdx) = getTeamsInLeagueMatch(matchday, matchIdxInDay);
        return (getTeamIdxInCup(groupIdx, homeIdx), getTeamIdxInCup(groupIdx, visitorIdx)); 
    }
    
    function getTeamsInLeagueMatch(uint8 matchday, uint8 matchIdxInDay) public pure returns (uint8 homeIdx, uint8 visitorIdx) 
    {
        require(matchday < MATCHDAYS, "wrong match day");
        require(matchIdxInDay < MATCHES_PER_DAY, "wrong match");
        if (matchday < (TEAMS_PER_LEAGUE - 1))
            (homeIdx, visitorIdx) = _getTeamsInMatchFirstHalf(matchday, matchIdxInDay);
        else
            (visitorIdx, homeIdx) = _getTeamsInMatchFirstHalf(matchday - (TEAMS_PER_LEAGUE - 1), matchIdxInDay);
    }

    // TODO: do this by exact formula instead of brute force search
    function getMatchesForTeams(uint8 team0, uint8 team1) public pure returns (uint8 match0, uint8 match1) 
    {
        uint8 home;
        uint8 vist;
        for (uint8 m = 0; m < MATCHES_PER_LEAGUE; m++) {
            (home, vist) = getTeamsInLeagueMatch(m / 4, m % 4);
            if ((home == team0) && (vist == team1)) match0 = m;
            if ((home == team1) && (vist == team0)) match1 = m;
        }
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

    function computeTeamRankingPoints(
        uint256[PLAYERS_PER_TEAM_MAX] memory states,
        uint8 leagueRanking,
        uint256 prevPerfPoints
    ) 
        public
        pure
        returns (uint256 teamSkills)
    {
        for (uint8 p = 0; p < PLAYERS_PER_TEAM_MAX; p++) {
            if (states[p] != 0 && states[p] != FREE_PLAYER_ID)
                teamSkills += getSumOfSkills(states[p]);
        }
        
        // Nomenclature:    R = rankingPoints, W = Weight_Skills, SK = TeamSkills, SK0 = TeamSkillsAtStart, I = 
        //                  I = Inertia, I0 = inertia Max, P0 = prevPerfPoints, P1 = currenteLeaguePerfPoints
        // 
        // Note that we use P = [0, 20] instead of the user-facing P' = [-10, 10] to avoid negatives.
        // I/I0 is the percentage of the previous perfPoints that we carry here. 
        // Formula: R = W SK/SK0 + P - 10 = W SK/SK0 + (I P0 + (10-I) P1)/I0 - 10
        // So we can avoid dividing, and simply compute:  R * SK0 * I0 = W SK I0 + SK0 (I P0 + (10-I)P1) - 10 SK0 I0
        // Note that if we do not need to dive, we can just keep I = 4, I0 = 10
        //  R * SK0 * I0 = 10W SK + SK0 (I P0 + (10-I)P1 - 100) = 10 W SK + SK0 Pnow
        // And finall  RankingPoints = 10W SK + SK0 Pnow

        // The user knows that his performance points now are: (note I' = I/I0)
        //  Pnow' = I' P0' + (1-I')P1' = I' P0 + (1-I')P1 - 10 = Pnow/I0

        // Formula in terms of pos and neg terms:
        //   pos = 10 W SK + SK0 (I P0 + 10 P1),   neg = SK0 (I P1 + 100)
        uint256 perfPointsThisLeague = getPerfPoints(leagueRanking);
        uint256 pos = 10 * WEIGHT_SKILLS * teamSkills + SKILLS_AT_START * (INERTIA * prevPerfPoints + 10 * perfPointsThisLeague);
        uint256 neg = SKILLS_AT_START * (INERTIA * perfPointsThisLeague + 100);
        
        if (pos > neg) return pos-neg;
        else return 0;
    }

    function getPerfPoints(uint8 leagueRanking) public pure returns (uint256) {
        if (leagueRanking == 0) return 2;
        else if (leagueRanking == 1) return 5;
        else if (leagueRanking == 2) return 8;
        else if (leagueRanking == 3) return 10;
        else if (leagueRanking == 4) return 12;
        else if (leagueRanking == 5) return 15;
        else if (leagueRanking == 6) return 18;
        else return 20;
    }

    // returns two sorted lists, [best teamIdxInLeague, points], ....
    // idx in the N*(N-1) matrix, assuming t0 < t1
    function computeLeagueLeaderBoard(uint8[2][MATCHES_PER_LEAGUE] memory results, uint8 matchDay, uint256 matchDaySeed) public pure returns (
        uint8[TEAMS_PER_LEAGUE] memory ranking, uint256[TEAMS_PER_LEAGUE] memory points
    ) {
        require(matchDay < MATCHDAYS, "wrong matchDay");
        uint8 team0;
        uint8 team1;
        uint16[TEAMS_PER_LEAGUE]memory goals;
        for(uint8 m = 0; m < matchDay * 4; m++) {
            (team0, team1) = getTeamsInLeagueMatch(m / 4, m % 4); 
            goals[team0] += results[m][0];
            goals[team1] += results[m][1];
            if (results[m][0] == results[m][1]) {
                points[team0] += 1000000000;
                points[team1] += 1000000000;
            } else if (results[m][0] > results[m][1]) {
                points[team0] += 3000000000;
            } else {
                points[team1] += 3000000000;
            }
        }
        // note that both points and ranking are returned ordered: (but goals and goalsAverage remain with old idxs)
        for (uint8 i = 0; i < TEAMS_PER_LEAGUE; i++) ranking[i] = i;
        sortIdxs(points, ranking);
        uint8 lastNonTied;
        for (uint8 r = 0; r < TEAMS_PER_LEAGUE-1; r++) {
            if (points[r+1] != points[r] && lastNonTied == r) lastNonTied = r+1;
            else if (points[r+1] != points[r]) {
                computeSecondaryPoints(ranking, points, results, goals, lastNonTied, r, matchDaySeed);
                lastNonTied = r+1;
            }
        }
        if (points[TEAMS_PER_LEAGUE-1] == points[TEAMS_PER_LEAGUE-2]) {
            computeSecondaryPoints(ranking, points, results, goals, lastNonTied, TEAMS_PER_LEAGUE-1, matchDaySeed);
        }
        sortIdxs(points, ranking);
    }
    
    // Points = nPoints in league * 1e9 + bestDirects * 1e6 + nGoalsInLeague * 1e3 + random % 999
    function computeSecondaryPoints(
        uint8[TEAMS_PER_LEAGUE] memory ranking,
        uint256[TEAMS_PER_LEAGUE] memory points,
        uint8[2][MATCHES_PER_LEAGUE] memory results,
        uint16[TEAMS_PER_LEAGUE]memory goals,
        uint8 firstTeamInRank,
        uint8 lastTeamInRank,
        uint256 matchDaySeed
    ) public pure {
        for (uint8 team0 = firstTeamInRank; team0 <= lastTeamInRank; team0++) {
            points[team0] += uint256(goals[ranking[team0]])*1000 + (matchDaySeed >> team0 * 10) % 999;
            for (uint8 team1 = team0 + 1; team1 <= lastTeamInRank; team1++) {
                uint8 bestTeam = computeDirect(results, ranking[team0], ranking[team1]);
                if (bestTeam == 0) points[team0] += 1000000;
                else if (bestTeam == 1) points[team1] += 1000000;
            }        
        }
    }

    function computeDirect(uint8[2][MATCHES_PER_LEAGUE] memory results, uint8 team0, uint8 team1) public pure returns(uint8 bestTeam) {
        (uint8 match0, uint8 match1) = getMatchesForTeams(team0, team1);
        if (results[match0][0] + results[match1][1] > results[match0][1] + results[match1][0]) return 0;
        else if (results[match0][0] + results[match1][1] < results[match0][1] + results[match1][0]) return 1;
        else return 2;
    }
}