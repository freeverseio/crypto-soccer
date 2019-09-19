pragma solidity ^0.5.0;

import "./Encoding.sol";

contract Engine is Encoding{
    
    uint8 private constant ROUNDS_PER_MATCH = 18;   // Number of relevant actions that happen during a game (18 equals one per 5 min)
    uint8 private constant RNDS_PER_UINT    = 18;   // Num of short nums that fit in a bignum = (256/ BITS_PER_RND);
    uint256 private constant BITS_PER_RND   = 14;   // Number of bits allowed for random numbers inside match decisisons
    uint256 private constant MAX_RND        = 16383;// Max random number allowed inside match decisions = 2^BITS_PER_RND-1 
    uint256 private constant MASK           = (1 << BITS_PER_RND)-1; // = (2**bits)-1, MASK used to extract short nums from bignum
    uint256 public constant MAX_PENALTY    = 10000; // Idx used to identify normal player acting as GK, or viceversa.
    // // Idxs for vector of globSkills: [0=move2attack, 1=globSkills[IDX_CREATE_SHOOT], 2=globSkills[IDX_DEFEND_SHOOT], 3=blockShoot, 4=currentEndurance]
    uint8 private constant IDX_MOVE2ATTACK  = 0;        
    uint8 private constant IDX_CREATE_SHOOT = 1; 
    uint8 private constant IDX_DEFEND_SHOOT = 2; 
    uint8 private constant IDX_BLOCK_SHOOT  = 3; 
    uint8 private constant IDX_ENDURANCE    = 4; 
    uint256 private constant TENTHOUSAND    = uint256(10000); 
    uint256 private constant TENTHOUSAND_SQ = uint256(100000000); 

    
    bool dummyBoolToEstimateCost;

    // mock up to estimate cost of a match.
    // to be removed before deployment
    // function playMatchWithCost(
    //     uint256 seed,
    //     uint256[PLAYERS_PER_TEAM_MAX][2] memory states,
    //     uint256[2] memory tactics
    // )
    //     public
    // {
    //     playMatch(seed, states, tactics);
    //     dummyBoolToEstimateCost = !dummyBoolToEstimateCost; 
    // }
    /**
     * @dev playMatch returns the result of a match
     * @param seed the pseudo-random number to use as a seed for the match
     * @param states a 2-vector, each of the 2 being vector with the state of the players of team 0
     * @param tactics a 2-vector with the tacticId (ex. 0 for [4,4,2]) for each team
     * @return the score of the match
     */
    function playMatch(
        uint256 seed,
        uint256[PLAYERS_PER_TEAM_MAX][2] memory states,
        uint256[2] memory tactics
    )
        public
        pure
        returns (uint8[2] memory teamGoals) 
    {
        uint8[11][2] memory lineups;
        uint8[9][2] memory playersPerZone;
        uint16[] memory rnds = getNRandsFromSeed(ROUNDS_PER_MATCH*4, seed);
        uint[5][2] memory globSkills;
        uint8[10][2] memory fwdModifiers;
        
        (lineups[0], fwdModifiers[0], playersPerZone[0]) = getLineUpAndPlayerPerZone(tactics[0]);
        (lineups[1], fwdModifiers[1], playersPerZone[1]) = getLineUpAndPlayerPerZone(tactics[1]);
        globSkills[0] = getTeamGlobSkills(states[0], playersPerZone[0], lineups[0], fwdModifiers[0]);
        globSkills[1] = getTeamGlobSkills(states[1], playersPerZone[1], lineups[1], fwdModifiers[0]);
        uint8 teamThatAttacks;
        for (uint8 round = 0; round < ROUNDS_PER_MATCH; round++){
            if ((round == 8) || (round == 13)) {
                (globSkills[0], globSkills[1]) = teamsGetTired(globSkills[0], globSkills[1]);
            }
            teamThatAttacks = throwDice(globSkills[0][IDX_MOVE2ATTACK], globSkills[1][IDX_MOVE2ATTACK], rnds[4*round]);
            if ( managesToShoot(teamThatAttacks, globSkills, rnds[4*round+1])) {
                if ( managesToScore(
                    states[teamThatAttacks],
                    playersPerZone[teamThatAttacks],
                    lineups[teamThatAttacks],
                    globSkills[1-teamThatAttacks][IDX_BLOCK_SHOOT],
                    rnds[4*round+2],
                    rnds[4*round+3]
                    )
                ) 
                {
                    teamGoals[teamThatAttacks]++;
                }
            }
        }
        return teamGoals;
    }
    
    function getNDefenders(uint8[9] memory playersPerZone) private pure returns (uint8) {
        return 2 * playersPerZone[0] + playersPerZone[1];
    }

    function getNMidfielders(uint8[9] memory playersPerZone) private pure returns (uint8) {
        return 2 * playersPerZone[3] + playersPerZone[4];
    }

    function getNAtackers(uint8[9] memory playersPerZone) private pure returns (uint8) {
        return 2 * playersPerZone[6] + playersPerZone[7];
    }

    // translates from a high level tacticsId (e.g. 442) to a format that describes how many
    // players play in each of the 9 zones in the field (Def, Mid, Forw) x (L, C, R), 
    // We impose left-right symmetry: DR = DL, MR = ML, FR = FL.
    // So we only manage 6 numbers: [DL, DM, ML, MM, FL, FM], and force 
    function getLineUpAndPlayerPerZone(uint256 tactics) 
        public 
        pure 
        returns (uint8[11] memory lineup, uint8[10] memory fwdModifiers, uint8[9] memory playersPerZone) 
    {
        uint8 tacticsId;
        (lineup, fwdModifiers, tacticsId) = decodeTactics(tactics);
        return (lineup, fwdModifiers, getPlayersPerZone(tacticsId));
    }

    // TODO: can this be expressed as
    // translates from a high level tacticsId (e.g. 442) to a format that describes how many
    // players play in each of the 9 zones in the field (Def, Mid, Forw) x (L, C, R), 
    // We impose left-right symmetry: DR = DL, MR = ML, FR = FL.
    // So we only manage 6 numbers: [DL, DM, ML, MM, FL, FM], and force 
    function getPlayersPerZone(uint8 tacticsId) internal pure returns (uint8[9] memory) {
        require(tacticsId < 4, "we currently support only 4 different tactics");
        if (tacticsId == 0) return [1,2,1,1,2,1,0,2,0];  // 0 = 442
        if (tacticsId == 1) return [1,3,1,1,2,1,0,1,0];  // 0 = 541
        if (tacticsId == 2) return [1,2,1,1,1,1,1,1,1];  // 0 = 433
        if (tacticsId == 3) return [1,2,1,1,3,1,0,1,0];  // 0 = 451
    }

    /// @dev Rescales global skills of both teams according to their endurance
    function teamsGetTired(uint[5] memory skillsTeamA, uint[5]  memory skillsTeamB )
        public
        pure
        returns (uint[5] memory, uint[5] memory)
    {
        uint currentEnduranceA = skillsTeamA[IDX_ENDURANCE];
        uint currentEnduranceB = skillsTeamB[IDX_ENDURANCE];
        for (uint8 sk = IDX_MOVE2ATTACK; sk < IDX_ENDURANCE; sk++) {
            skillsTeamA[sk] = (skillsTeamA[sk] * currentEnduranceA) / 100;
            skillsTeamB[sk] = (skillsTeamB[sk] * currentEnduranceB) / 100;
        }
        return(skillsTeamA, skillsTeamB);
    }


    function getNRandsFromSeed(uint16 nRands, uint256 seed) public pure returns (uint16[] memory rnds) {
        rnds = new uint16[](nRands);
        uint256 currentBigRnd = uint(keccak256(abi.encode(seed)));
        uint8 rndsFromSameBigRnd = 0;
        for (uint8 n = 0; n < nRands; n++) {
            if (rndsFromSameBigRnd == RNDS_PER_UINT) {
                currentBigRnd = uint(keccak256(abi.encode(seed+1)));
                rndsFromSameBigRnd = 0;
            }
            rnds[n] = uint16(currentBigRnd & MASK);
            currentBigRnd >>= BITS_PER_RND;
            rndsFromSameBigRnd ++;
        }
        return rnds;
    }


    /// @dev Throws a dice that returns 0 with probability weight1/(weight1+weight2), and 1 otherwise.
    /// @dev So, returning 0 has semantics: "the responsible for weight1 is selected".
    /// @dev We return a uint8, not bool, to allow the return to be used as an idx in an array by the callee.
    /// @dev The formula is derived as follows. Throw a random number R in the range [0,maxR].
    /// @dev Then, w1 wins if (w1+w2)*(R/maxR) < w1, and w2 wins otherise. 
    /// @dev MAX_RND controls the resolution or fine-graining of the algorithm.
    function throwDice(uint weight1, uint weight2, uint rndNum) public pure returns(uint8) {
        if( ( (weight1 + weight2) * rndNum ) < ( weight1 * (MAX_RND-1) ) ) {
            return 0;
        } else {
            return 1;
        }
    }

    /// @dev Generalization of the previous to any number of input weights
    /// @dev It therefore throws any number of dice and returns the winner's idx.
    function throwDiceArray11(uint[11] memory weights, uint rndNum) public pure returns(uint8 w) {
        uint uniformRndInSumOfWeights;
        for (w = 0; w < 11; w++) {
            uniformRndInSumOfWeights += weights[w];
        }
        uniformRndInSumOfWeights *= rndNum;
        uint cumSum = 0;
        for (w = 0; w < 10; w++) {
            cumSum += weights[w];
            if( uniformRndInSumOfWeights < ( cumSum * (MAX_RND-1) )) {
                return w;
            }
        }
        return w;
    }


    /// @dev Decides if a team manages to shoot by confronting attack and defense via globSkills
    function managesToShoot(uint8 teamThatAttacks, uint[5][2] memory globSkills, uint rndNum)
        public
        pure
        returns (bool)
    {
        return throwDice(
            globSkills[1-teamThatAttacks][IDX_DEFEND_SHOOT],       // globSkills[IDX_DEFEND_SHOOT] of defending team against...
            (globSkills[teamThatAttacks][IDX_CREATE_SHOOT]*6)/10,  // globSkills[IDX_CREATE_SHOOT] of attacking team.
            rndNum
        ) == 1 ? true : false;
    }


    function selectShooter(
        uint256[PLAYERS_PER_TEAM_MAX] memory teamState,
        uint8[9] memory playersPerZone,
        uint8[11] memory lineup,
        uint rndNum1
    )
        public
        pure
        returns (uint8)
    {
        /// attacker who actually shoots is selected weighted by his speed
        uint256[11] memory weights;
        for (uint8 p = 11 - getNAtackers(playersPerZone); p < 11; p++) {
            weights[p] = getSpeed(teamState[lineup[p]]);
        }
        return throwDiceArray11(weights, rndNum1);
    }

    /// @dev Decides if a team that creates a shoot manages to score.
    /// @dev First: select attacker who manages to shoot. Second: challenge him with keeper
    function managesToScore(
        uint256[PLAYERS_PER_TEAM_MAX] memory teamState,
        uint8[9] memory playersPerZone,
        uint8[11] memory lineup,
        uint blockShoot,
        uint rndNum1,
        uint rndNum2
    )
        public
        pure
        returns (bool)
    {
        uint8 shooter = selectShooter(teamState, playersPerZone, lineup, rndNum1);

        /// a goal is scored by confronting his shoot skill to the goalkeeper block skill
        return throwDice((getSpeed(teamState[lineup[shooter]])*7)/10, blockShoot, rndNum2) == 0;
        // return false;
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
        uint256[PLAYERS_PER_TEAM_MAX] memory teamState, 
        uint8[9] memory playersPerZone, 
        uint8[11] memory lineup,
        uint8[10] memory fwdModifiers
    )
        public
        pure
        returns (
            uint256[5] memory globSkills
        )
    {
        // for a keeper, the 'shoot skill' is interpreted as block skill
        // if for whatever reason, user places a non-GK as GK, the block skill is a terrible minimum.
        uint256 penalty;
        uint256 playerSkills = teamState[lineup[0]];
        globSkills[IDX_ENDURANCE] = getEndurance(playerSkills);
        if (computePenalty(0, playersPerZone, playerSkills) == 0) {globSkills[IDX_BLOCK_SHOOT] = 10;}
        else globSkills[IDX_BLOCK_SHOOT] = getShoot(playerSkills);
            
        
        uint256[3] memory fwdModFactors;
        uint8 p = 1;
        // loop over defenders
        for (uint8 i = 0; i < getNDefenders(playersPerZone); i++) {
            playerSkills = teamState[lineup[p]];
            penalty = computePenalty(p, playersPerZone, playerSkills);
            if (penalty != 0) {
                fwdModFactors = getFwdModFactors(fwdModifiers[p-1]);
                globSkills[IDX_MOVE2ATTACK] += ((getDefence(playerSkills) + getSpeed(playerSkills) + getPass(playerSkills)) * penalty * fwdModFactors[IDX_MOVE2ATTACK])/TENTHOUSAND_SQ;
                globSkills[IDX_DEFEND_SHOOT] += ((getDefence(playerSkills) + getSpeed(playerSkills)) * penalty * fwdModFactors[IDX_MOVE2ATTACK])/TENTHOUSAND_SQ;
                globSkills[IDX_ENDURANCE]   += ((getEndurance(playerSkills)) * penalty)/TENTHOUSAND;
            } else {
                globSkills[IDX_MOVE2ATTACK] += 30;
                globSkills[IDX_DEFEND_SHOOT] += 20;
                globSkills[IDX_ENDURANCE]   += 10;
            }
            p++;
        }
        // loop over midfielders
        for (uint8 i = 0; i < getNMidfielders(playersPerZone); i++) {
            playerSkills = teamState[lineup[p]];
            penalty = computePenalty(p, playersPerZone, playerSkills);
            fwdModFactors = getFwdModFactors(fwdModifiers[p-1]);
            if (penalty != 0) {
                penalty = computePenalty(p, playersPerZone, playerSkills);
                globSkills[IDX_MOVE2ATTACK] += ((2*getDefence(playerSkills) + 2*getSpeed(playerSkills) + 3*getPass(playerSkills)) * penalty * fwdModFactors[IDX_MOVE2ATTACK])/TENTHOUSAND_SQ;
                globSkills[IDX_ENDURANCE]   += ((getEndurance(playerSkills)) * penalty)/TENTHOUSAND;
            } else {
                globSkills[IDX_MOVE2ATTACK] += 50;
                globSkills[IDX_ENDURANCE]   += 10;
            }
            p++;
        }
        // loop over strikers
        for (uint8 i = 0; i < getNAtackers(playersPerZone); i++) {
            playerSkills = teamState[lineup[p]];
            penalty = computePenalty(p, playersPerZone, playerSkills);
            if (penalty != 0) {
                penalty = computePenalty(p, playersPerZone, playerSkills);
                globSkills[IDX_MOVE2ATTACK] += ((getDefence(playerSkills)) * penalty * fwdModFactors[IDX_MOVE2ATTACK])/TENTHOUSAND_SQ;
                globSkills[IDX_CREATE_SHOOT] += ((getSpeed(playerSkills) + getPass(playerSkills)) * penalty * fwdModFactors[IDX_MOVE2ATTACK])/TENTHOUSAND_SQ;
                globSkills[IDX_ENDURANCE] += ((getEndurance(playerSkills)) * penalty)/TENTHOUSAND;
            } else {
                globSkills[IDX_MOVE2ATTACK] += 10;
                globSkills[IDX_CREATE_SHOOT] += 20;
                globSkills[IDX_ENDURANCE] += 10;
            }
            p++;
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

    // recall order: [MOVE2ATTACK, CREATE_SHOOT, DEFEND_SHOOT, BLOCK_SHOOT, ENDURANCE]
    // the forward modifier factors only change the first 3.
    function getFwdModFactors(uint8 fwdModifier) public pure returns (uint256[3] memory fwdModFactors) {
        if (fwdModifier == 0)       {fwdModFactors = [TENTHOUSAND, TENTHOUSAND, TENTHOUSAND];}
        else if (fwdModifier == 1)  {fwdModFactors = [uint256(9500),  uint256(9500),  uint256(10500)];}
        else if (fwdModifier == 2)  {fwdModFactors = [uint256(10500), uint256(10500), uint256(9500)];}
    }

  
    // 0 penalty means no penalty
    // 1000 penalty means 10% penalty
    // etc... up to MAX_PENALTY
    function computePenalty(
        uint8 lineupPos, 
        uint8[9] memory playersPerZone, 
        uint256 playerSkills
    ) 
        public
        pure
        returns (uint256 penalty) 
    {
        require(lineupPos < 11, "wrong arg in computePenalty");
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
        return 10000-penalty; 
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

