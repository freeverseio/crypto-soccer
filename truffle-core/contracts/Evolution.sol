pragma solidity ^0.5.0;

import "./Assets.sol";
import "./EncodingSkills.sol";
import "./EngineLib.sol";
import "./EncodingMatchLog.sol";
import "./Engine.sol";
import "./EncodingTPAssignment.sol";
import "./EncodingSkillsSetters.sol";

contract Evolution is EncodingMatchLog, EncodingSkills, EngineLib, EncodingTPAssignment, EncodingSkillsSetters {

    uint8 constant public PLAYERS_PER_TEAM_MAX  = 25;
    uint8 public constant NO_OUT_OF_GAME_PLAYER  = 14;   // noone saw a card
    uint8 public constant RED_CARD = 3;   // noone saw a card
    uint256 constant public MAX_DIFF  = 10; // beyond this diff among team qualities, it's basically infinite
    uint256 constant public POINTS_FOR_HAVING_PLAYED  = 10; // beyond this diff among team qualities, it's basically infinite
    uint8 private constant IDX_IS_2ND_HALF      = 0; 
    uint256 private constant SECS_IN_DAY    = 86400; // 24 * 3600 
    uint8 constant public N_SKILLS = 5;
    uint8 constant public SK_SHO = 0;
    uint8 constant public SK_SPE = 1;
    uint8 constant public SK_PAS = 2;
    uint8 constant public SK_DEF = 3;
    uint8 constant public SK_END = 4;

    Assets private _assets;
    Engine private _engine;

    bool dummyBoolToEstimateCost;

    function setAssetsAddress(address addr) public {
        _assets = Assets(addr);
    }

    function setEngine(address addr) public {
        _engine = Engine(addr);
    }

    // function to call on 2nd half when we want to the matchlog to include the evolution points too.
    function play2ndHalfAndEvolve(
        uint256 seed,
        uint256 matchStartTime,
        uint256[PLAYERS_PER_TEAM_MAX][2] memory states,
        uint256[2] memory tactics,
        uint256[2] memory matchLog,
        bool[3] memory matchBools // [is2ndHalf, isHomeStadium, isPlayoff]
    )
        public view returns(uint256[2] memory)
    {
        require(matchBools[IDX_IS_2ND_HALF], "play with evolution should only be called in 2nd half games");
        return computeTrainingPoints(
            _engine.playHalfMatch(seed, matchStartTime, states, tactics, matchLog, matchBools)
        );
    }

    function computeTrainingPointsWithCost(uint256[2] memory matchLog) public returns (uint256[2] memory)
    {
        dummyBoolToEstimateCost = !dummyBoolToEstimateCost;
        return computeTrainingPoints(matchLog);
    }

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
                    (getOutOfGameType(matchLog[team], false) == RED_CARD ? 3 : 0)
                +   (getOutOfGameType(matchLog[team], true)  == RED_CARD ? 3 : 0)
                +   ((getYellowCard(matchLog[team], 0, false) < NO_OUT_OF_GAME_PLAYER) ? 1 : 0) 
                +   ((getYellowCard(matchLog[team], 1, false) < NO_OUT_OF_GAME_PLAYER) ? 1 : 0)
                +   ((getYellowCard(matchLog[team], 0, true)  < NO_OUT_OF_GAME_PLAYER) ? 1 : 0) 
                +   ((getYellowCard(matchLog[team], 1, true)  < NO_OUT_OF_GAME_PLAYER) ? 1 : 0);
        }
        
        // require(pointsNeg[0] == 10, "....");
        
        // subtract points, keeping them always non-negativre
        points[0] = (points[0] > pointsNeg[0]) ? (points[0] - pointsNeg[0]) : 0;
        points[1] = (points[1] > pointsNeg[1]) ? (points[1] - pointsNeg[1]) : 0;
        
        // +10% for each extra 50 points of lack of balance between teams
        uint256 teamSumSkills0 = getTeamSumSkills(matchLog[0]);
        uint256 teamSumSkills1 = getTeamSumSkills(matchLog[1]);

        if (teamSumSkills0 > teamSumSkills1) {
            points[0] = (points[0] * teamSumSkills1) / (teamSumSkills0);
            points[1] = (points[1] * teamSumSkills0) / (teamSumSkills1);
        } else if (teamSumSkills0 < teamSumSkills1) {
            points[0] = (points[0] * teamSumSkills1) / (teamSumSkills0);
            points[1] = (points[1] * teamSumSkills0) / (teamSumSkills1);
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
            quality +=  getShoot(state) + getSpeed(state) + getPass(state)
                    +   getDefence(state) + getEndurance(state);
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
    
    function getTeamEvolvedSkills(
        uint256[PLAYERS_PER_TEAM_MAX] memory states, 
        uint256 TPs, 
        uint256 userAssignment,
        uint256 matchStartTime
    ) 
        public
        view
        returns (uint256[PLAYERS_PER_TEAM_MAX] memory)
    {
        (uint8[25] memory weights, uint8 specialPlayer) = decodeTP(userAssignment);
        uint8[5] memory singleWeights;
        
        for (uint8 p = 0; p < PLAYERS_PER_TEAM_MAX; p++) {
            uint256 skills = states[p];
            if (skills == 0) continue; 
            uint8 offset = 0;
            if (p == specialPlayer) offset = 20;
            else if(getForwardness(skills) == IDX_GK) offset = 0;
            else if(getForwardness(skills) == IDX_D) offset = 5;
            else if(getForwardness(skills) == IDX_F) offset = 15;
            else offset = 10;
            for (uint8 w = 0; w < 5; w++) singleWeights[w] = weights[offset + w];
            states[p] = evolvePlayer(skills, TPs, singleWeights, matchStartTime);
        }    
        return states;
    }
    
    // formula: TP(i) = w(i)/100 * TP, where TP are the TPs, S(i) is the weight for skill(i)
    // deltaS(i)    = w(i)/100 * max[ TP,  TP* (pot * 4/3 - (age-16)/2) ] - max(0,(age-31)*8)
    // If age is in days, define Yd = year2days
    // deltaS(i)    = w(i)/100 * max[ TP,  TP * (pot * 8 * Yd - 3 * ageDays + 48 Yd)/ (6 Yd)] - max(0,(ageDays-31)*8/Yd)
    // If age is in secs, define Ys = year2secs
    // deltaS(i)    = w(i)/100 * max[ TP,  TP * (pot * 8 * Ys - 3 * ageInSecs + 48 Ys)/ (6 Ys)] - max(0,(ageInSecs-31)*8/Ys)
    // skill(i)     = max(0, skill(i) + deltaS(i))
    // 
    // shoot, speed, pass, defence, endurance
    function evolvePlayer(uint256 skills, uint256 TPs, uint8[5] memory weights, uint256 matchStartTime) public view returns(uint256) {
        uint256 ageInSecs = 7 * (matchStartTime - getBirthDay(skills) * 86400);  // 86400 = day2secs

        uint256 potential = getPotential(skills);
        uint256 deltaNeg = (ageInSecs > 977616000) ? ((ageInSecs-977616000)*8)/31536000 : 0;  // 977616000 = 31 * Ys, 31536000 = Ys
        uint256 multiplier;
        // if (potential * 2920 + 17520 > 3 * ageDays + 2190) {
        if (potential * 252288000 + 1513728000 > 3 * ageInSecs + 189216000) {  // 252288000 = 8 Ys,  1513728000 = 48 Ys, 189216000 = 6 Ys
            multiplier = (TPs*(potential * 252288000 + 1513728000 - 3 * ageInSecs))/189216000;
        } else {
            multiplier = TPs;
        }
        // 0: shoot
        if (getShoot(skills) + (multiplier * weights[SK_SHO])/100 > deltaNeg) {
            skills = setShoot(skills, getShoot(skills) + (multiplier * weights[SK_SHO])/100 - deltaNeg);
        } else {
            skills = setShoot(skills, 1);
        }
        // 1: speed
        if (getSpeed(skills) + (multiplier * weights[SK_SPE])/100 > deltaNeg) {
            skills = setSpeed(skills, getSpeed(skills) + (multiplier * weights[SK_SPE])/100 - deltaNeg);
        } else {
            skills = setSpeed(skills, 1);
        }
        // 2: pass
        if (getPass(skills) + (multiplier * weights[SK_PAS])/100 > deltaNeg) {
            skills = setPass(skills, getPass(skills) + (multiplier * weights[SK_PAS])/100 - deltaNeg);
        } else {
            skills = setPass(skills, 1);
        }
        // 3: defence
        if (getDefence(skills) + (multiplier * weights[SK_DEF])/100 > deltaNeg) {
            skills = setDefence(skills, getDefence(skills) + (multiplier * weights[SK_DEF])/100 - deltaNeg);
        } else {
            skills = setDefence(skills, 1);
        }
        // 4: endurance
        if (getEndurance(skills) + (multiplier * weights[SK_END])/100 > deltaNeg) {
            skills = setEndurance(skills, getEndurance(skills) + (multiplier * weights[SK_END])/100 - deltaNeg);
        } else {
            skills = setEndurance(skills, 1);
        }
        // 5: sumSkills
        skills = setSumOfSkills(skills, uint32(getShoot(skills) + getSpeed(skills) + getPass(skills) + getDefence(skills) + getEndurance(skills)));
        return generateChildIfNeeded(skills, ageInSecs, matchStartTime);
    } 

    function generateChildIfNeeded(uint256 skills, uint256 ageInSecs, uint256 matchStartTime) public view returns (uint256) {
        if ((getSumOfSkills(skills) > 200) && (ageInSecs < 1166832000)) {   // 1166832000 = 37 * Ys
            return skills;
        }
        uint256 dna = uint256(keccak256(abi.encode(skills, ageInSecs)));
        ageInSecs = 504576000 + (dna % 126144000);  // 504576000 = 16 * years2secs, 126144000 = 4 * years2secs
        uint256 dayOfBirth = (matchStartTime - ageInSecs / 7)/86400; // 86400 = 24 * 3600
        dna >>= 13; // log2(7300) = 12.8
        uint256 playerId = getPlayerIdFromSkills(skills);
        uint8 shirtNum = uint8(_assets.getCurrentShirtNum(_assets.getPlayerStateAtBirth(playerId)));
        (uint16[N_SKILLS] memory newSkills, uint8[4] memory birthTraits, uint32 sumSkills) = _assets.computeSkills(dna, shirtNum);
        encodePlayerSkills(newSkills, dayOfBirth, uint8(getGeneration(skills) + 1), playerId, birthTraits, false, false, 0, 0, false, sumSkills);
    }
}

