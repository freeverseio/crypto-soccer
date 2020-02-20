pragma solidity >=0.5.12 <0.6.2;

import "./Assets.sol";
import "./Market.sol";
import "./EngineLib.sol";
import "./EncodingMatchLog.sol";
import "./Engine.sol";
import "./EncodingTPAssignment.sol";
import "./EncodingSkills.sol";
import "./EncodingSkillsSetters.sol";
import "./EncodingTacticsPart2.sol";

contract TrainingPoints is EncodingMatchLog, EngineLib, EncodingTPAssignment, EncodingSkills, EncodingSkillsSetters, EncodingTacticsPart2 {

    // uint8 constant public PLAYERS_PER_TEAM_MAX  = 25;
    uint8 public constant NO_OUT_OF_GAME_PLAYER  = 14;   // noone saw a card
    uint8 public constant RED_CARD = 3;   // noone saw a card
    uint256 constant public POINTS_FOR_HAVING_PLAYED  = 10; // beyond this diff among team qualities, it's basically infinite
    // uint8 constant public N_SKILLS = 5;

    Assets private _assets;
    Market private _market;

    function setAssetsAddress(address addr) public {
        _assets = Assets(addr);
    }
    function setMarketAddress(address addr) public {
        _market = Market(addr);
    }

    function computeTrainingPoints(uint256 matchLog0, uint256 matchLog1) public pure returns (uint256, uint256)
    {
        // +11 point for winning at home, +22 points for winning
        // away, or in a cup match. 0 points for drawing.
        uint256 nGoals0 = getNGoals(matchLog0);
        uint256 nGoals1 = getNGoals(matchLog1);
        uint256[2] memory points;
        points[0] = POINTS_FOR_HAVING_PLAYED;
        points[1] = POINTS_FOR_HAVING_PLAYED;

        if (getWinner(matchLog0)==0) { // we can get winner from [0] or [1], they are the same   
            points[0] += (getIsHomeStadium(matchLog0) ? 11 : 22); // we can get homeStadium from [0] or [1], they are the same   
        } else if (getWinner(matchLog0)==1) {
            points[1] += (getIsHomeStadium(matchLog0) ? 22 : 22);    
        }
        
        // +6 for goal scored by GK/D; +5 for midfielder; +4 for attacker; +3 for each assist
        points[0] += pointsPerWhoScoredGoalsAndAssists(matchLog0, nGoals0);
        points[1] += pointsPerWhoScoredGoalsAndAssists(matchLog1, nGoals1);

        // if clean-sheet (opponent did not score):
        // +2 per half played by GK/D, +1 per half played for Mids and Atts
        if (nGoals1 == 0) points[0] += pointsPerCleanSheet(matchLog0);
        if (nGoals0 == 0) points[1] += pointsPerCleanSheet(matchLog1);

        uint256[2] memory pointsNeg;
        // -1 for each opponent goal
        pointsNeg[0] = nGoals1;
        pointsNeg[1] = nGoals0;
        // -3 for redCards, -1 for yellows
        for (uint8 team = 0; team <2; team++) {
            uint256 thisLog = (team == 0 ? matchLog0 : matchLog1);
            pointsNeg[team] += 
                    (getOutOfGameType(thisLog, false) == RED_CARD ? 3 : 0)
                +   (getOutOfGameType(thisLog, true)  == RED_CARD ? 3 : 0)
                +   ((getYellowCard(thisLog, 0, false) < NO_OUT_OF_GAME_PLAYER) ? 1 : 0) 
                +   ((getYellowCard(thisLog, 1, false) < NO_OUT_OF_GAME_PLAYER) ? 1 : 0)
                +   ((getYellowCard(thisLog, 0, true)  < NO_OUT_OF_GAME_PLAYER) ? 1 : 0) 
                +   ((getYellowCard(thisLog, 1, true)  < NO_OUT_OF_GAME_PLAYER) ? 1 : 0);
        }
        
        // subtract points, keeping them always non-negativre
        points[0] = (points[0] > pointsNeg[0]) ? (points[0] - pointsNeg[0]) : 0;
        points[1] = (points[1] > pointsNeg[1]) ? (points[1] - pointsNeg[1]) : 0;
        
        // +10% for each extra 50 points of lack of balance between teams
        uint256 teamSumSkills0 = getTeamSumSkills(matchLog0);
        uint256 teamSumSkills1 = getTeamSumSkills(matchLog1);

        if (teamSumSkills0 > teamSumSkills1) {
            points[0] = (points[0] * teamSumSkills1) / (teamSumSkills0);
            points[1] = (points[1] * teamSumSkills0) / (teamSumSkills1);
        } else if (teamSumSkills0 < teamSumSkills1) {
            points[0] = (points[0] * teamSumSkills1) / (teamSumSkills0);
            points[1] = (points[1] * teamSumSkills0) / (teamSumSkills1);
        }
        matchLog0 = addTrainingPoints(matchLog0, points[0]);
        matchLog1 = addTrainingPoints(matchLog1, points[1]);
        return (matchLog0, matchLog1);
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
    
    
    
    function computeTeamQuality(uint256[PLAYERS_PER_TEAM_MAX] memory teamSkills) public pure returns (uint256 quality) {
        uint256 skills;
        for (uint8 p = 0; p < PLAYERS_PER_TEAM_MAX; p++) {
            skills = teamSkills[p];
            if (skills != 0) {
                for (uint8 sk = 0; sk < N_SKILLS; sk++) quality += getSkill(skills, sk); 
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
    
    function applyTrainingPoints(
        uint256[PLAYERS_PER_TEAM_MAX] memory teamSkills, 
        uint256 assignedTPs,
        uint256 tactics,
        uint256 matchStartTime,
        uint16 earnedTPs
    ) 
        public
        view
        returns (uint256[PLAYERS_PER_TEAM_MAX] memory)
    {
        if (assignedTPs == 0) return teamSkills;
        (uint16[25] memory TPperSkill, uint8 specialPlayer, uint16 TP) = decodeTP(assignedTPs);
        require(earnedTPs == TP, "assignedTPs used an amount of TP that does not match the earned TPs in previous match");
        uint16[5] memory singleTPperSkill;
        (uint8[PLAYERS_PER_TEAM_MAX] memory staminas,,) = getItemsData(tactics);
                
        // note that if no special player was selected => specialPlayer = PLAYERS_PER_TEAM_MAX 
        // ==> it will never be processed in this loop
        for (uint8 p = 0; p < PLAYERS_PER_TEAM_MAX; p++) {
            uint256 thisSkills = teamSkills[p];
            if (thisSkills == 0) continue; 
            if (staminas[p] > 0) thisSkills = reduceGamesNonStopping(thisSkills, staminas[p]);
            uint8 offset = 0;
            if (p == specialPlayer) offset = 20; 
            else if(getForwardness(thisSkills) == IDX_GK) offset = 0;
            else if(getForwardness(thisSkills) == IDX_D) offset = 5;
            else if(getForwardness(thisSkills) == IDX_F) offset = 15;
            else offset = 10;
            for (uint8 s = 0; s < 5; s++) singleTPperSkill[s] = TPperSkill[offset + s];
            teamSkills[p] = evolvePlayer(thisSkills, singleTPperSkill, matchStartTime);
            
        }    
        return teamSkills;
    }
    
    // deltaS(i)    = max[ TP(i), TP(i) * (pot * 4/3 - (age-16)/2) ] - max(0,(age-31)*8)
    // If age is in days, define Yd = year2days
    // deltaS(i)    = max[ TP(i), TP(i) * (pot * 8 * Yd - 3 * ageDays + 48 Yd)/ (6 Yd)] - max(0,(ageDays-31)*8/Yd)
    // If age is in secs, define Ys = year2secs
    // deltaS(i)    = max[ TP(i), TP(i) * (pot * 8 * Ys - 3 * ageInSecs + 48 Ys)/ (6 Ys)] - max(0,(ageInSecs-31)*8/Ys)
    // skill(i)     = max(0, skill(i) + deltaS(i))
    // deltaS(i)    = max[ TP(i), TP(i) * numerator / denominator] - max(0,(ageInSecs-31)*8/Ys)
    // skill(i)     = max(0, skill(i) + deltaS(i))
    // shoot, speed, pass, defence, endurance
    function evolvePlayer(uint256 skills, uint16[5] memory TPperSkill, uint256 matchStartTime) public view returns(uint256) {
        uint256 ageInSecs = 7 * (matchStartTime - getBirthDay(skills) * 86400);  // 86400 = day2secs
        uint256 deltaNeg = (ageInSecs > 977616000) ? ((ageInSecs-977616000)*8)/31536000 : 0;  // 977616000 = 31 * Ys, 31536000 = Ys
        uint256 numerator;
        if (getPotential(skills) * 252288000 + 1513728000 > 3 * ageInSecs) {  // 252288000 = 8 Ys,  1513728000 = 48 Ys, 189216000 = 6 Ys
            numerator = (getPotential(skills) * 252288000 + 1513728000 - 3 * ageInSecs);
        } else {
            numerator = 0;
        }
        uint256 sum;
        for (uint8 sk = 0; sk < N_SKILLS; sk++) {
            uint256 newSkill = getNewSkill(getSkill(skills, sk), TPperSkill[sk], numerator, 189216000, deltaNeg);
            skills = setSkill(skills, newSkill, sk);
            sum += newSkill;
        }
        skills = setSumOfSkills(skills, uint32(sum));
        return generateChildIfNeeded(skills, ageInSecs, matchStartTime);
    } 
    
    // stamina = 0 => do not reduce
    // stamina = 1 => reduce 2 games
    // stamina = 2 => reduce 4 games
    // stamina = 3 => full recovery
    function reduceGamesNonStopping(uint256 skills, uint8 stamina) public pure returns (uint256) {
        require(stamina < 4, "stamina value too large");
        uint8 gamesNonStopping = getGamesNonStopping(skills);
        if (gamesNonStopping == 0) return skills;
        if ((stamina == 3) || (gamesNonStopping <= 2 * stamina)) return setGamesNonStopping(skills, 0);
        return setGamesNonStopping(skills, gamesNonStopping - 2 * stamina);
    }

    function getNewSkill(uint256 oldSkill, uint16 TPthisSkill, uint256 numerator, uint256 denominator, uint256 deltaNeg) private pure returns (uint256) {
        uint256 term1 = (TPthisSkill*numerator) / denominator;
        term1 = (term1 > TPthisSkill) ? term1 : TPthisSkill;
        if ((oldSkill + term1) > deltaNeg) return oldSkill + term1 - deltaNeg;
        return 1;
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
        uint8 shirtNum = uint8(_assets.getCurrentShirtNum(_market.getPlayerStateAtBirth(playerId)));
        (uint16[N_SKILLS] memory newSkills, uint8[4] memory birthTraits, uint32 sumSkills) = _assets.computeSkills(dna, shirtNum);
        // if dna is even => leads to child, if odd => leads to academy player
        uint8 generation = uint8((getGeneration(skills) % 32) + 1 + (dna % 2 == 0 ? 0 : 32));
        return encodePlayerSkills(newSkills, dayOfBirth, generation, playerId, birthTraits, false, false, 0, 0, false, sumSkills);
    }
}

