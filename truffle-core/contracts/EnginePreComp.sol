pragma solidity ^0.5.0;

import "./EncodingSkills.sol";
import "./EngineLib.sol";

contract EnginePreComp is EngineLib {
    uint256 private constant ONE256       = 1; 
    // // Idxs for vector of globSkills: [0=move2attack, 1=globSkills[IDX_CREATE_SHOOT], 2=globSkills[IDX_DEFEND_SHOOT], 3=blockShoot, 4=currentEndurance]
    uint8 private constant IDX_MOVE2ATTACK  = 0;        
    uint8 private constant IDX_CREATE_SHOOT = 1; 
    uint8 private constant IDX_DEFEND_SHOOT = 2; 
    uint8 private constant IDX_BLOCK_SHOOT  = 3; 
    uint8 private constant IDX_ENDURANCE    = 4; 
    uint256 private constant TEN_TO_4       = uint256(10000); 
    uint256 private constant TEN_TO_10      = uint256(10000000000); 
    uint256 private constant TEN_TO_14      = uint256(100000000000000); 
    uint256 private constant SECS_IN_DAY    = 86400; // 24 * 3600 



    // Over a game, we would like:
    //      - injuryHard = 1 per 100 games => 0.01 per game per player => 0.02 per game
    //      - injuryLow = 0.7 per 100 games => 0.007 per game per player => 0.04 per game
    //      - redCard 1/10 = 0.1 per game
    //      - yellowCard 2.5 per game 
    // We encode this in uint16[3] events, which applies to 1 half of the game only.
    //  - 1 possible event that leaves a player out of the match, encoded in:
    //          events[0, 1] = [player (from 0 to 13), eventType (injuryHard, injuryLow, redCard)]
    //  - 2 possible events for yellow card:
    //          events[2] = player (from 0 to 13)
    //          events[3] = player (from 0 to 13)
    // The player value is set to NO_EVENT ( = 14) if no event took place
    // If we're on the 2nd half, the idx are events[4,5,6,7]
    // for out of game:
    //      it cannot return 0
    //      injuryHard:  1
    //      injuryLow:   2
    //      redCard:     3
    function computeExceptionalEvents
    (
        uint256 matchLog,
        uint256[PLAYERS_PER_TEAM_MAX] memory states,
        bool is2ndHalf,
        uint256 seed
    ) 
        public 
        pure 
        returns 
    (
        uint256
    ) 
    {
        uint8 offset = is2ndHalf ? 169 : 151;
        uint256[] memory weights = new uint256[](15);
        uint64[] memory rnds = getNRandsFromSeed(seed + 42, 4);
        for (uint8 p = 0; p < 14; p++) {
            if (states[p] != 0) weights[p] = 1 + getAggressiveness(states[p]); // weights must be > 0 to ever be selected
        }
        // events[0] => STUFF THAT REMOVES A PLAYER FROM FIELD: injuries and redCard 
        // average sumAggressiveness = 11 * 2.5 = 27.5
        // total = 0.07 per game = 0.035 per half => weight nothing happens = 758
        weights[14] = 758;
        matchLog |= uint256(throwDiceArray(weights, rnds[0])) << offset;
        matchLog |= uint256(computeRound(rnds[0]+1)) << offset + 4;
        matchLog |= uint256(computeTypeOfEvent(rnds[1])) << (offset + 8);
        // next: two events for yellow cards
        // average sumAggressiveness = 11 * 2.5 = 27.5
        // total = 2.5 per game = 1.25 per half => 0.75 per dice thrown
        // weight nothing happens = 9
        weights[14] = 9;
        matchLog |= uint256(throwDiceArray(weights, rnds[2])) << (offset + 10);
        matchLog |= uint256(throwDiceArray(weights, rnds[3])) << (offset + 14);
        
        return matchLog;
    }
    
    function computeRound(uint256 seed) private pure returns (uint8 round) {
        return uint8(seed % 12);
    }

    // it cannot return 0.
    // injuryHard:  1
    // injuryLow:   2
    // redCard:     3
    function computeTypeOfEvent(uint256 rnd) private pure returns (uint8) {
        uint256[] memory weights = new uint256[](3);
        weights[0] = 1; // injuryHard   
        weights[1] = 2; // injuryLow
        weights[2] = 5; // redCard
        return 1 + throwDiceArray(weights, rnd);
    }

    function computePenalties(
        uint256[2] memory matchLog, 
        uint256[PLAYERS_PER_TEAM_MAX][2] memory states, 
        uint256 block0, 
        uint256 block1, 
        uint64 seed
    )
        public 
        pure 
        returns(uint256[2] memory) 
    {
        uint64[] memory rnds = getNRandsFromSeed(seed * 7, 14);
        uint8[2] memory totalGoals;
        for (uint8 round = 0; round < 6; round++) {
            if (throwDice(block1, 3 * getShoot(states[0][10-round]), rnds[2 *round]) == 1) {
                matchLog[0] |= (ONE256 << 144 + round);
                totalGoals[0] += 1;
            }
            if (throwDice(block0, 3 * getShoot(states[1][10-round]), rnds[2 *round + 1]) == 1) {
                matchLog[1] |= (ONE256 << 144 + round);
                totalGoals[1] += 1;
            }
            if ((round > 3) && (totalGoals[0] != totalGoals[1])) return matchLog;
        }
        if (throwDice(block0 + getShoot(states[0][4]), block0 + getShoot(states[0][4]), rnds[13]) == 1) {
            matchLog[0] |= (ONE256 << 144 + 6);
        } else {
            matchLog[1] |= (ONE256 << 144 + 6);
        }
        return matchLog;
    }

    /// @dev Computes basic data, including globalSkills, needed during the game.
    /// @dev Basically implements the formulas:
    // move2attack =    defence(defenders + 2*midfields + attackers) +
    //                  speed(defenders + 2*midfields) +
    //                  pass(defenders + 3*midfields)
    // globSkills[IDX_CREATE_SHOOT] =    speed(attackers) + pass(attackers)
    // globSkills[IDX_DEFEND_SHOOT] =    speed(defenders) + defence(defenders);
    // blockShoot  =    shoot(keeper);
    function getTeamGlobSkills(
        uint256[PLAYERS_PER_TEAM_MAX] memory states, 
        uint8[9] memory playersPerZone, 
        bool[10] memory extraAttack,
        uint256 matchStartTime
    )
        public
        pure
        returns 
    (
        uint256[5] memory globSkills
    )
    {
        // for a keeper, the 'shoot skill' is interpreted as block skill
        // if for whatever reason, user places a non-GK as GK, the block skill is a terrible minimum.
        uint256 penalty;
        uint256 playerSkills = states[0];
        globSkills[IDX_ENDURANCE] = getEndurance(playerSkills);
        if (computePenaltyBadPositionAndCondition(0, playersPerZone, playerSkills) == 0) {
            globSkills[IDX_BLOCK_SHOOT] = (10 * penaltyPerAge(playerSkills, matchStartTime))/1000000;
        }
        else globSkills[IDX_BLOCK_SHOOT] = (getShoot(playerSkills) * penaltyPerAge(playerSkills, matchStartTime))/1000000;
        
        uint256[3] memory fwdModFactors;

        for (uint8 p = 1; p < 11; p++){
            playerSkills = states[p];
            penalty = computePenaltyBadPositionAndCondition(p, playersPerZone, playerSkills) * penaltyPerAge(playerSkills, matchStartTime);
            fwdModFactors = getExtraAttackFactors(extraAttack[p-1]);
            if (p < 1 + getNDefenders(playersPerZone)) {computeDefenderGlobSkills(globSkills, playerSkills, penalty, fwdModFactors);}
            else if (p < 1 + getNDefenders(playersPerZone) + getNMidfielders(playersPerZone)) {computeMidfielderGlobSkills(globSkills, playerSkills, penalty, fwdModFactors);}
            else {computeForwardsGlobSkills(globSkills, playerSkills, penalty, fwdModFactors);}       
        }
        // endurance is converted to a percentage, 
        // used to multiply (and hence decrease) the start endurance.
        // 100 is super-endurant (1500), 70 is bad, for an avg starting team (550).
        if (globSkills[IDX_ENDURANCE] < 500) {
            globSkills[IDX_ENDURANCE] = 70;
        } else if (globSkills[IDX_ENDURANCE] < 1400) {
            globSkills[IDX_ENDURANCE] = 100 - (1400-globSkills[IDX_ENDURANCE])/30;
        } else {
            globSkills[IDX_ENDURANCE] = 100;
        }
    }

    // 0 penalty means no penalty
    // 1000 penalty means 10% penalty
    // etc... up to MAX_PENALTY
    function computePenaltyBadPositionAndCondition(
        uint8 lineupPos, 
        uint8[9] memory playersPerZone, 
        uint256 playerSkills
    ) 
        public
        pure
        returns (uint256 penalty) 
    {
        require(lineupPos < 11, "wrong arg in computePenaltyBadPositionAndCondition");
        uint256 forwardness = getForwardness(playerSkills);
        uint256 leftishness = getLeftishness(playerSkills);
        if (forwardness == IDX_GK && lineupPos > 0 || forwardness != IDX_GK && lineupPos == 0) return 0;
        uint8[9] memory playersBelow = playersBelowZones(playersPerZone);
        lineupPos--; // remove the offset due to the GK
        if (lineupPos < playersBelow[0]) { 
            // assigned to defense left
            penalty = penaltyForDefenders(forwardness);
            penalty += penaltyForLefts(leftishness);
        } else if (lineupPos < playersBelow[1]) { 
            // assigned to defense center
            penalty = penaltyForDefenders(forwardness);
            penalty += penaltyForCenters(leftishness);
        } else if (lineupPos < playersBelow[2]) { 
            // assigned to defense left
            penalty = penaltyForDefenders(forwardness);
            penalty += penaltyForRights(leftishness);
        } else if (lineupPos < playersBelow[3]) { 
            // assigned to mid left
            penalty = penaltyForMids(forwardness);
            penalty += penaltyForLefts(leftishness);
        } else if (lineupPos < playersBelow[4]) { 
            // assigned to mid center
            penalty = penaltyForMids(forwardness);
            penalty += penaltyForCenters(leftishness);
        } else if (lineupPos < playersBelow[5]) { 
            // assigned to mid right
            penalty = penaltyForMids(forwardness);
            penalty += penaltyForRights(leftishness);
        } else if (lineupPos < playersBelow[6]) { 
            // assigned to attack left
            penalty = penaltyForAttackers(forwardness);
            penalty += penaltyForLefts(leftishness);
        } else if (lineupPos < playersBelow[7]) { 
            // assigned to attack center
            penalty = penaltyForAttackers(forwardness);
            penalty += penaltyForCenters(leftishness);
        } else { 
            // assigned to attack right
            penalty = penaltyForAttackers(forwardness);
            penalty += penaltyForRights(leftishness);
        }
        uint256 gamesNonStop = getGamesNonStopping(playerSkills);
        if (gamesNonStop > 5) {
            return 5000 - penalty;
        } else {
            return 10000 - gamesNonStop * 1000 - penalty;
        }
    }

    // recall order: [MOVE2ATTACK, CREATE_SHOOT, DEFEND_SHOOT, BLOCK_SHOOT, ENDURANCE]
    // the forward modifier factors only change the first 3.
    function getExtraAttackFactors(bool extraAttack) public pure returns (uint256[3] memory fwdModFactors) {
        if (extraAttack)    {fwdModFactors = [uint256(10500), uint256(10500), uint256(9500)];}
        else                {fwdModFactors = [TEN_TO_4, TEN_TO_4, TEN_TO_4];}
    }
  
    function computeDefenderGlobSkills(
        uint256[5] memory globSkills,
        uint256 playerSkills, 
        uint256 penalty, 
        uint256[3] memory fwdModFactors
    ) 
        private 
        pure
    {
        if (penalty != 0) {
            globSkills[IDX_MOVE2ATTACK] += ((getDefence(playerSkills) + getSpeed(playerSkills) + getPass(playerSkills)) * penalty * fwdModFactors[IDX_MOVE2ATTACK])/TEN_TO_14;
            globSkills[IDX_DEFEND_SHOOT] += ((getDefence(playerSkills) + getSpeed(playerSkills)) * penalty * fwdModFactors[IDX_MOVE2ATTACK])/TEN_TO_14;
            globSkills[IDX_ENDURANCE]   += ((getEndurance(playerSkills)) * penalty)/TEN_TO_10;
        } else {
            globSkills[IDX_MOVE2ATTACK] += 30;
            globSkills[IDX_DEFEND_SHOOT] += 20;
            globSkills[IDX_ENDURANCE]   += 10;
        }
    }


    function computeMidfielderGlobSkills(
        uint256[5] memory globSkills,
        uint256 playerSkills, 
        uint256 penalty, 
        uint256[3] memory fwdModFactors
    ) 
        private 
        pure
    {
        if (penalty != 0) {
            globSkills[IDX_MOVE2ATTACK] += ((2*getDefence(playerSkills) + 2*getSpeed(playerSkills) + 3*getPass(playerSkills)) * penalty * fwdModFactors[IDX_MOVE2ATTACK])/TEN_TO_14;
            globSkills[IDX_ENDURANCE]   += ((getEndurance(playerSkills)) * penalty)/TEN_TO_10;
        } else {
            globSkills[IDX_MOVE2ATTACK] += 50;
            globSkills[IDX_ENDURANCE]   += 10;
        }
    }
    
    
    function computeForwardsGlobSkills(
        uint256[5] memory globSkills,
        uint256 playerSkills, 
        uint256 penalty, 
        uint256[3] memory fwdModFactors
    ) 
        private 
        pure
    {
        if (penalty != 0) {
            globSkills[IDX_MOVE2ATTACK] += ((getDefence(playerSkills)) * penalty * fwdModFactors[IDX_MOVE2ATTACK])/TEN_TO_14;
            globSkills[IDX_CREATE_SHOOT] += ((getSpeed(playerSkills) + getPass(playerSkills)) * penalty * fwdModFactors[IDX_MOVE2ATTACK])/TEN_TO_14;
            globSkills[IDX_ENDURANCE] += ((getEndurance(playerSkills)) * penalty)/TEN_TO_10;
        } else {
            globSkills[IDX_MOVE2ATTACK] += 10;
            globSkills[IDX_CREATE_SHOOT] += 20;
            globSkills[IDX_ENDURANCE] += 10;
        }
    }

    function playersBelowZones(uint8[9] memory playersPerZone) private pure returns(uint8[9] memory  playersBelow) {
        playersBelow[0] = playersPerZone[0];    
        for (uint8 z = 1; z < 9; z++) {
            playersBelow[z] = playersBelow[z-1] + playersPerZone[z];
        }
    }

    function penaltyForLefts(uint256 leftishness) private pure returns(uint16) {
            if (leftishness == IDX_C || leftishness == IDX_CR) {return 1000;} 
            else if (leftishness == IDX_R) {return 2000;}
    }

    function penaltyForCenters(uint256 leftishness) private pure returns(uint16) {
            if (leftishness == IDX_L || leftishness == IDX_R) {return 1000;} 
    }

    function penaltyForRights(uint256 leftishness) private pure returns(uint16) {
            if (leftishness == IDX_C || leftishness == IDX_LC) {return 1000;} 
            else if (leftishness == IDX_L) {return 2000;}
    }
    
    function penaltyForDefenders(uint256 forwardness) private pure returns(uint16) {
            if (forwardness == IDX_M || forwardness == IDX_MF) {return 1000;}
            else if (forwardness == IDX_F) {return 2000;}
    }

    function penaltyForMids(uint256 forwardness) private pure returns(uint16) {
            if (forwardness == IDX_D || forwardness == IDX_F) {return 1000;}
    }

    function penaltyForAttackers(uint256 forwardness) private pure returns(uint16) {
            if (forwardness == IDX_M || forwardness == IDX_MD) {return 1000;}
            else if (forwardness == IDX_D) {return 2000;}
    }


}

