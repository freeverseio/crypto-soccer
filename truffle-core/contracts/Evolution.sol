pragma solidity ^0.5.0;

import "./EncodingSkills.sol";

contract Evolution is EncodingSkills {

    uint8 constant public PLAYERS_PER_TEAM_MAX  = 25;
    uint8 public constant NO_CARD  = 14;   // noone saw a card
    uint8 public constant RED_CARD = 3;   // noone saw a card
    uint256 constant public FREE_PLAYER_ID  = 1; // it never corresponds to a legit playerId due to its TZ = 0
    uint256 constant public MAX_DIFF  = 10; // beyond this diff among team qualities, it's basically infinite

    function computeTrainingPoints( 
        uint256[2] memory matchLog,
        uint256[PLAYERS_PER_TEAM_MAX][2] memory statesHalf1,
        uint256[PLAYERS_PER_TEAM_MAX][2] memory statesHalf2,
        uint256[2] memory tacticsHalf1,
        uint256[2] memory tacticsHalf2,
        bool isHomeStadium
    )
        public
        pure
        returns (uint256[2] memory points)
    {
        
        // +1 point for winning at home, +2 points for winning
        // away, or in a cup match. 0 points for drawing.
        uint256 nGoals0 = (matchLog[0] & 15);
        uint256 nGoals1 = (matchLog[1] & 15);
        if ( nGoals0 > nGoals1) {
            points[0] = isHomeStadium ? 1 : 2;    
        } else if (nGoals0 < nGoals1) {
            points[1] = isHomeStadium ? 2 : 2;    
        }

        // +6 for goal scored by GK/D; +5 for midfielder; +4 for attacker; +3 for each assist
        points[0] += pointsPerWhoScoredGoalsAndAssists(matchLog[0], nGoals0);
        points[1] += pointsPerWhoScoredGoalsAndAssists(matchLog[1], nGoals1);

        // +4 for GK/D who saw no card at all, +1 for Mids and Atts


        uint256[2] memory pointsNeg;
        // -1 for each opponent goal
        pointsNeg[0] = nGoals1;
        pointsNeg[1] = nGoals0;
        // -3 for redCards, -1 for yellows
        // ...note that offset for 1st half is 159, and for 2nd half is 179
        for (uint8 team = 0; team <2; team++) {
            for (uint8 offset = 159; offset < 180; offset += 20) {
                pointsNeg[team] += (((matchLog[0] >> offset) & 3) < RED_CARD) ? 3 : 0;
                pointsNeg[team] += (((matchLog[0] >> (offset + 2)) & 15) < NO_CARD) ? 1 : 0;
                pointsNeg[team] += (((matchLog[0] >> (offset + 6)) & 15) < NO_CARD) ? 1 : 0;
            }
        }
        
        points[0] = (points[0] > pointsNeg[0]) ? (points[0] - pointsNeg[0]) : 0;
        points[1] = (points[1] > pointsNeg[1]) ? (points[1] - pointsNeg[1]) : 0;
        
        // +10% for each extra 50 points of lack of balance between teams
        uint256 teamQuality0 = computeTeamQuality(statesHalf1[0]);
        uint256 teamQuality1 = computeTeamQuality(statesHalf1[1]);

        if (teamQuality0 > teamQuality1) {
            points[0] = (points[0] * teamQuality1 * 3) / (teamQuality0 * 4);
            points[1] = (points[1] * teamQuality0 * 4) / (teamQuality1 * 3);
        } else if (teamQuality0 < teamQuality1) {
            points[0] = (points[0] * teamQuality1 * 4) / (teamQuality0 * 3);
            points[1] = (points[1] * teamQuality0 * 3) / (teamQuality1 * 4);
        }
    }
    
    function computeTeamQuality(uint256[PLAYERS_PER_TEAM_MAX] memory states) public pure returns (uint256 quality) {
        uint256 state;
        for (uint8 p = 0; p < PLAYERS_PER_TEAM_MAX; p++) {
            state = states[p];
            if (state != FREE_PLAYER_ID) {
                quality +=  getShoot(state) + getSpeed(state) + getPass(state)
                        +   getDefence(state) + getEndurance(state);
            }
        }
    }
    
    // +6 for goal scored by GK/D; +5 for midfielder; +4 for attacker; +3 for each assist
    function pointsPerWhoScoredGoalsAndAssists(uint256 matchLog, uint256 nGoals) public pure returns(uint256 points) {
        for (uint8 goal = 0; goal < nGoals; goal++) {
            uint256 fwdPos = (matchLog >> 116 + 2 * goal) & 3;
            if (fwdPos < 2) {points += 6;}
            else if (fwdPos == 2) {points += 5;}
            else {points += 6;}
            // if assister is different the shooter, it was a true assist
            if (((matchLog >> 4 + 4 * goal) & 15) != ((matchLog >> 60 + 4 * goal) & 15)) {points += 3;}
            
        }
    }

}

