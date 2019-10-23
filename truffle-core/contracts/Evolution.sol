pragma solidity ^0.5.0;

import "./EncodingSkills.sol";
import "./EngineLib.sol";
import "./EncodingMatchLog.sol";

contract Evolution is EncodingMatchLog, EncodingSkills, EngineLib {

    uint8 constant public PLAYERS_PER_TEAM_MAX  = 25;
    uint8 public constant NO_CARD  = 14;   // noone saw a card
    uint8 public constant RED_CARD = 3;   // noone saw a card
    uint256 constant public FREE_PLAYER_ID  = 1; // it never corresponds to a legit playerId due to its TZ = 0
    uint256 constant public MAX_DIFF  = 10; // beyond this diff among team qualities, it's basically infinite
    uint256 constant public POINTS_FOR_HAVING_PLAYED  = 10; // beyond this diff among team qualities, it's basically infinite

    function computeTrainingPoints(uint256[2] memory matchLog) public pure returns (uint256[2] memory)
    {
        
        // +11 point for winning at home, +22 points for winning
        // away, or in a cup match. 0 points for drawing.
        uint256 nGoals0 = getNGoals(matchLog[0]);
        uint256 nGoals1 = getNGoals(matchLog[1]);
        uint256[2] memory points;
        points[0] = POINTS_FOR_HAVING_PLAYED;
        points[1] = POINTS_FOR_HAVING_PLAYED;

        if (getWinner(matchLog[0])==0) { // we can get winner from [0] or [1], they are the same   
            points[0] += (getIsHomeStadium(matchLog[0]) ? 11 : 22); // we can get homeStadium from [0] or [1], they are the same   
        } else if (getWinner(matchLog[0])==1) {
            points[1] += (getIsHomeStadium(matchLog[0]) ? 22 : 22);    
        }

        // +6 for goal scored by GK/D; +5 for midfielder; +4 for attacker; +3 for each assist
        points[0] += pointsPerWhoScoredGoalsAndAssists(matchLog[0], nGoals0);
        points[1] += pointsPerWhoScoredGoalsAndAssists(matchLog[1], nGoals1);

        // if clean-sheet (opponent did not score):
        // +2 per half played by GK/D, +1 per half played for Mids and Atts
        if (nGoals1 == 0) points[0] += pointsPerCleanSheet(matchLog[0]);
        if (nGoals0 == 0) points[1] += pointsPerCleanSheet(matchLog[1]);

        uint256[2] memory pointsNeg;
        // -1 for each opponent goal
        pointsNeg[0] = nGoals1;
        pointsNeg[1] = nGoals0;
        // -3 for redCards, -1 for yellows
        for (uint8 team = 0; team <2; team++) {
            pointsNeg[team] += 
                    3 * (getOutOfGameType(matchLog[team], false) + getOutOfGameType(matchLog[team], true)) 
                +   (getYellowCard(matchLog[team], 0, false) < NO_CARD ? 1 : 0) 
                +   (getYellowCard(matchLog[team], 1, false) < NO_CARD ? 1 : 0)
                +   (getYellowCard(matchLog[team], 0, true)  < NO_CARD ? 1 : 0) 
                +   (getYellowCard(matchLog[team], 1, true)  < NO_CARD ? 1 : 0);
        }
        
        // subtract points, keeping them always non-negativre
        points[0] = (points[0] > pointsNeg[0]) ? (points[0] - pointsNeg[0]) : 0;
        points[1] = (points[1] > pointsNeg[1]) ? (points[1] - pointsNeg[1]) : 0;
        
        // +10% for each extra 50 points of lack of balance between teams
        uint256 teamSumSkills0 = getTeamSumSkills(matchLog[0]);
        uint256 teamSumSkills1 = getTeamSumSkills(matchLog[1]);

        if (teamSumSkills0 > teamSumSkills1) {
            points[0] = (points[0] * teamSumSkills1 * 3) / (teamSumSkills0 * 4);
            points[1] = (points[1] * teamSumSkills0 * 4) / (teamSumSkills1 * 3);
        } else if (teamSumSkills0 < teamSumSkills1) {
            points[0] = (points[0] * teamSumSkills1 * 4) / (teamSumSkills0 * 3);
            points[1] = (points[1] * teamSumSkills0 * 3) / (teamSumSkills1 * 4);
        }
        matchLog[0] = addTrainingPoints(matchLog[0], points[0]);
        matchLog[1] = addTrainingPoints(matchLog[1], points[1]);
        return matchLog;
    }
    
    // if clean-sheet (opponent did not score):
    // +2 per half played by GK/D, +1 per half played for Mids and Atts
    function pointsPerCleanSheet(uint256 matchLog) public pure returns (uint256) {
        // formula: (note that for a given half: 1 + nDef + nMid + nAtt = nTot)
        //      pointsPerHalf   = 2 (for GK) + 2 * nDef + nMid + nAtt 
        //                      = 2 + 2 * nDef + nTot - nDef - 1 = nTot + 1 + nDef
        //      note also that by constraint, nTot = 11 in the first half
        //      pointsPerMatch  = 2 + nTot1 + nTot2 + nDef1 + nDef2 = 13 + nTot2 + nDef1 + nDef2 
        return 13   + (getOutOfGameType(matchLog, false) == RED_CARD ? 10 : 11) 
                    +  getNDefs(matchLog, false) + getNDefs(matchLog, true);
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
            uint256 fwdPos = getForwardPos(matchLog, goal);
            if (fwdPos < 2) {points += 6;}
            else if (fwdPos == 2) {points += 5;}
            else {points += 4;}
            // if assister is different the shooter, it was a true assist
            if (getShooter(matchLog, goal) != getAssister(matchLog, goal)) {points += 3;}
        }
    }

}

