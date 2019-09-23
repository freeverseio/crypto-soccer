pragma solidity ^0.5.0;

contract Engine {
    
    uint8 public constant ROUNDS_PER_MATCH = 12;   // Number of relevant actions that happen during a game (12 equals one per 3.7 min)
    uint8 private constant BITS_PER_RND     = 36;   // Number of bits allowed for random numbers inside match decisisons
    uint256 public constant MAX_RND         = 68719476735; // Max random number allowed inside match decisions: 2^36-1
    uint256 public constant MAX_PENALTY     = 10000; // Idx used to identify normal player acting as GK, or viceversa.
    // // Idxs for vector of globSkills: [0=move2attack, 1=globSkills[IDX_CREATE_SHOOT], 2=globSkills[IDX_DEFEND_SHOOT], 3=blockShoot, 4=currentEndurance]
    uint8 private constant IDX_MOVE2ATTACK  = 0;        
    uint8 private constant IDX_CREATE_SHOOT = 1; 
    uint8 private constant IDX_DEFEND_SHOOT = 2; 
    uint8 private constant IDX_BLOCK_SHOOT  = 3; 
    uint8 private constant IDX_ENDURANCE    = 4; 
    uint256 private constant TENTHOUSAND    = uint256(10000); 
    uint256 private constant TENTHOUSAND_SQ = uint256(100000000); 
    // skills idxs: Defence, Speed, Pass, Shoot, Endurance
    uint8 constant public SK_SHO = 0;
    uint8 constant public SK_SPE = 1;
    uint8 constant public SK_PAS = 2;
    uint8 constant public SK_DEF = 3;
    uint8 constant public SK_END = 4;
    uint8 constant public N_SKILLS = 5;
    // prefPosition idxs: GoalKeeper, Defender, Midfielder, Forward, MidDefender, MidAttacker
    uint8 constant public IDX_GK = 0;
    uint8 constant public IDX_D  = 1;
    uint8 constant public IDX_M  = 2;
    uint8 constant public IDX_F  = 3;
    uint8 constant public IDX_MD = 4;
    uint8 constant public IDX_MF = 5;
    //  Leftishness:   0: 000, 1: 001, 2: 010, 3: 011, 4: 100, 5: 101, 6: 110, 7: 111
    uint8 constant public IDX_R = 1;
    uint8 constant public IDX_C = 2;
    uint8 constant public IDX_CR = 3;
    uint8 constant public IDX_L = 4;
    uint8 constant public IDX_LR = 5;
    uint8 constant public IDX_LC = 6;
    uint8 constant public IDX_LCR = 7;
    //  Bools
    uint8 constant public IS_2ND_HALF = 0;
    uint8 constant public IS_HOME_STADIUM = 1;



    
    bool dummyBoolToEstimateCost;

    /**
     * @dev playMatch returns the result of a match
     * @return the score of the match
     */
    function playMatch(
        uint256 seed,
        uint16[5][11][2] memory playerSkills,
        bool[10][2] memory extraAttack,
        uint8[11][2] memory forwardness,
        uint8[11][2] memory leftishness,
        uint8[9][2] memory playersPerZone,
        bool[2] memory matchBools // [is2ndHalf, isHomeStadium]
    )
        public
        pure
        returns (uint8[2] memory teamGoals) 
    {
        uint64[ROUNDS_PER_MATCH*4] memory rnds = getNRandsFromSeed(seed);
        uint256[5][2] memory globSkills;
        
        globSkills[0] = getTeamGlobSkills(playerSkills[0], playersPerZone[0], extraAttack[0], forwardness[0], leftishness[0]);
        globSkills[1] = getTeamGlobSkills(playerSkills[1], playersPerZone[1], extraAttack[1], forwardness[1], leftishness[1]);
        if (matchBools[IS_HOME_STADIUM]) {
            globSkills[IDX_ENDURANCE][0] = (globSkills[IDX_ENDURANCE][0] * 11500)/10000;
        }
        uint8 teamThatAttacks;
        for (uint8 round = 0; round < ROUNDS_PER_MATCH; round++){
            if (matchBools[IS_2ND_HALF] && ((round == 0) || (round == 5))) {
                (globSkills[0], globSkills[1]) = teamsGetTired(globSkills[0], globSkills[1]);
            }
            teamThatAttacks = throwDice(globSkills[0][IDX_MOVE2ATTACK], globSkills[1][IDX_MOVE2ATTACK], rnds[4*round]);
            if ( managesToShoot(teamThatAttacks, globSkills, rnds[4*round+1])) {
                if ( managesToScore(
                    playerSkills[teamThatAttacks],
                    playersPerZone[teamThatAttacks],
                    extraAttack[teamThatAttacks],
                    forwardness[teamThatAttacks],
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

    function getNAttackers(uint8[9] memory playersPerZone) private pure returns (uint8) {
        return 2 * playersPerZone[6] + playersPerZone[7];
    }


    /// @dev Rescales global skills of both teams according to their endurance
    function teamsGetTired(uint256[5] memory skillsTeamA, uint256[5]  memory skillsTeamB )
        public
        pure
        returns (uint256[5] memory, uint256[5] memory)
    {
        uint256 currentEnduranceA = skillsTeamA[IDX_ENDURANCE];
        uint256 currentEnduranceB = skillsTeamB[IDX_ENDURANCE];
        for (uint8 sk = IDX_MOVE2ATTACK; sk < IDX_ENDURANCE; sk++) {
            skillsTeamA[sk] = (skillsTeamA[sk] * currentEnduranceA) / 100;
            skillsTeamB[sk] = (skillsTeamB[sk] * currentEnduranceB) / 100;
        }
        return(skillsTeamA, skillsTeamB);
    }


    function getNRandsFromSeed(uint256 seed) public pure returns (uint64[ROUNDS_PER_MATCH*4] memory rnds) {
        uint256 currentBigRnd = uint256(keccak256(abi.encode(seed)));
        uint8 remainingBits = 255;
        for (uint8 n = 0; n < ROUNDS_PER_MATCH*4; n++) {
            if (remainingBits < BITS_PER_RND) {
                currentBigRnd = uint256(keccak256(abi.encode(seed, n)));
                remainingBits = 255;
            }
            rnds[n] = uint64(currentBigRnd & MAX_RND);
            currentBigRnd >>= BITS_PER_RND;
            remainingBits -= BITS_PER_RND;
        }
        return rnds;
    }


    /// @dev Throws a dice that returns 0 with probability weight0/(weight0+weight1), and 1 otherwise.
    /// @dev So, returning 0 has semantics: "the responsible for weight0 is selected".
    /// @dev We return a uint8, not bool, to allow the return to be used as an idx in an array by the callee.
    /// @dev The formula is derived as follows. Throw a random number R in the range [0,maxR].
    /// @dev Then, w0 wins if (w0+w1)*(R/maxR) < w0, and w1 wins otherise. 
    /// @dev MAX_RND controls the resolution or fine-graining of the algorithm.
    function throwDice(uint256 weight0, uint256 weight1, uint256 rndNum) public pure returns(uint8) {
        if( ( (weight0 + weight1) * rndNum ) < ( weight0 * (MAX_RND-1) ) ) {
            return 0;
        } else {
            return 1;
        }
    }

    /// @dev Generalization of the previous to any number of input weights
    /// @dev It therefore throws any number of dice and returns the winner's idx.
    function throwDiceArray11(uint256[11] memory weights, uint256 rndNum) public pure returns(uint8 w) {
        uint256 uniformRndInSumOfWeights;
        for (w = 0; w < 11; w++) {
            uniformRndInSumOfWeights += weights[w];
        }
        uniformRndInSumOfWeights *= rndNum;
        uint256 cumSum = 0;
        for (w = 0; w < 10; w++) {
            cumSum += weights[w];
            if( uniformRndInSumOfWeights < ( cumSum * (MAX_RND-1) )) {
                return w;
            }
        }
        return w;
    }


    /// @dev Decides if a team manages to shoot by confronting attack and defense via globSkills
    function managesToShoot(uint8 teamThatAttacks, uint256[5][2] memory globSkills, uint256 rndNum)
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

    function selectAssister(
        uint16[5][11] memory playerSkills,
        uint8[9] memory playersPerZone,
        bool[10] memory extraAttack,
        uint8 shooter,
        uint256 rnd
    )
        public
        pure
        returns (uint8)
    {
        uint256[11] memory weights;
        // if selected assister == selected shooter =>  
        //  there was no assist => individual play by shoorter
        weights[0] = 1;
        uint256 teamPassCapacity = 1;
        uint8 p = 1;
        for (uint8 i = 0; i < getNDefenders(playersPerZone); i++) {
            weights[p] = (extraAttack[p-1] ? 90 : 20 ) * playerSkills[p][SK_PAS];
            teamPassCapacity += weights[p];
            p++;
        }
        for (uint8 i = 0; i < getNMidfielders(playersPerZone); i++) {
            weights[p] = (extraAttack[p-1] ? 150 : 100 ) * playerSkills[p][SK_PAS];
            teamPassCapacity += weights[p];
            p++;
        }
        for (uint8 i = 0; i < getNAttackers(playersPerZone); i++) {
            weights[p] = 200 * playerSkills[p][SK_PAS];
            teamPassCapacity += weights[p];
            p++;
        }
        // on average: teamPassCapacity442 = (1 + 4 * 20 + 4 * 100 + 2 * 200) < getPass > = 881 <pass>_team
        // on average: shooterSumOfSkills = 5 * <skills>_shooter
        // so a good ratio is shooterSumOfSkills/teamPassCapacity442 = 5/881 * <skills_shooter>/<pass>_team
        // or better, to have an avg of 1: (shooterSumOfSkills*271)/(teamPassCapacity * 5) = <skills_shooter>/<pass>_team
        // or to have a 50% change, multiply by 10, and to have say, 1/3, multiply by 10/3
        // this is to be compensated by an overall factor of about.
        uint256 sum = playerSkills[shooter][0] + playerSkills[shooter][1] + playerSkills[shooter][2] + playerSkills[shooter][3] + playerSkills[shooter][4];
        weights[shooter] = (weights[shooter] * sum * 8810 )/ (N_SKILLS * (teamPassCapacity - weights[shooter]) * 3);
        return throwDiceArray11(weights, rnd);
    }


    function selectShooter(
        uint16[5][11] memory playerSkills,
        uint8[9] memory playersPerZone,
        bool[10] memory extraAttack,
        uint256 rnd
    )
        public
        pure
        returns (uint8)
    {
        uint256[11] memory weights;
        // GK has minimum weight, all others are relative to this.
        weights[0] = 1;
        uint8 p = 1;
        for (uint8 i = 0; i < getNDefenders(playersPerZone); i++) {
            weights[p] = (extraAttack[p-1] ? 15000 : 5000 ) * playerSkills[p][SK_SPE];
            p++;
        }
        for (uint8 i = 0; i < getNMidfielders(playersPerZone); i++) {
            weights[p] = (extraAttack[p-1] ? 50000 : 25000 ) * playerSkills[p][SK_SPE];
            p++;
        }
        for (uint8 i = 0; i < getNAttackers(playersPerZone); i++) {
            weights[p] = 75000 * playerSkills[p][SK_SPE];
            p++;
        }
        return throwDiceArray11(weights, rnd);
    }

    /// @dev Decides if a team that creates a shoot manages to score.
    /// @dev First: select attacker who manages to shoot. Second: challenge him with keeper
    function managesToScore(
        uint16[5][11] memory playerSkills,
        uint8[9] memory playersPerZone,
        bool[10] memory extraAttack,
        uint8[11] memory forwardness,
        uint256 blockShoot,
        uint256 rndNum1,
        uint256 rndNum2
    )
        public
        pure
        returns (bool)
    {
        uint8 shooter = selectShooter(playerSkills, playersPerZone, extraAttack, rndNum1);

        /// a goal is scored by confronting his shoot skill to the goalkeeper block skill
        uint256 shootPenalty = forwardness[shooter] == IDX_GK ? 10 : 1;
        return throwDice((playerSkills[shooter][SK_SHO]*7)/(shootPenalty*10), blockShoot, rndNum2) == 0;
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
        uint16[5][11] memory playerSkills,
        uint8[9] memory playersPerZone, 
        bool[10] memory extraAttack,
        uint8[11] memory forwardness,
        uint8[11] memory leftishness

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
        globSkills[IDX_ENDURANCE] = playerSkills[0][SK_END];
        if (computePenalty(0, playersPerZone,  forwardness[0], leftishness[0]) == 0) {globSkills[IDX_BLOCK_SHOOT] = 10;}
        else globSkills[IDX_BLOCK_SHOOT] = playerSkills[0][SK_SHO];
            
        
        uint256[3] memory fwdModFactors;
        uint8 p = 1;
        // loop over defenders
        for (uint8 i = 0; i < getNDefenders(playersPerZone); i++) {
            penalty = computePenalty(p, playersPerZone, forwardness[p], leftishness[p]);
            if (penalty != 0) {
                fwdModFactors = getExtraAttackFactors(extraAttack[p-1]);
                globSkills[IDX_MOVE2ATTACK] += ((playerSkills[p][SK_DEF] + playerSkills[p][SK_SPE] + playerSkills[p][SK_PAS]) * penalty * fwdModFactors[IDX_MOVE2ATTACK])/TENTHOUSAND_SQ;
                globSkills[IDX_DEFEND_SHOOT] += ((playerSkills[p][SK_DEF] + playerSkills[p][SK_SPE]) * penalty * fwdModFactors[IDX_MOVE2ATTACK])/TENTHOUSAND_SQ;
                globSkills[IDX_ENDURANCE]   += ((playerSkills[p][SK_END]) * penalty)/TENTHOUSAND;
            } else {
                globSkills[IDX_MOVE2ATTACK] += 30;
                globSkills[IDX_DEFEND_SHOOT] += 20;
                globSkills[IDX_ENDURANCE]   += 10;
            }
            p++;
        }
        // loop over midfielders
        for (uint8 i = 0; i < getNMidfielders(playersPerZone); i++) {
            penalty = computePenalty(p, playersPerZone, forwardness[p], leftishness[p]);
            if (penalty != 0) {
                fwdModFactors = getExtraAttackFactors(extraAttack[p-1]);
                globSkills[IDX_MOVE2ATTACK] += ((2*playerSkills[p][SK_DEF] + 2*playerSkills[p][SK_SPE] + 3*playerSkills[p][SK_PAS]) * penalty * fwdModFactors[IDX_MOVE2ATTACK])/TENTHOUSAND_SQ;
                globSkills[IDX_ENDURANCE]   += ((playerSkills[p][SK_END]) * penalty)/TENTHOUSAND;
            } else {
                globSkills[IDX_MOVE2ATTACK] += 50;
                globSkills[IDX_ENDURANCE]   += 10;
            }
            p++;
        }
        // loop over strikers
        for (uint8 i = 0; i < getNAttackers(playersPerZone); i++) {
            penalty = computePenalty(p, playersPerZone, forwardness[p], leftishness[p]);
            if (penalty != 0) {
                globSkills[IDX_MOVE2ATTACK] += ((playerSkills[p][SK_DEF]) * penalty * fwdModFactors[IDX_MOVE2ATTACK])/TENTHOUSAND_SQ;
                globSkills[IDX_CREATE_SHOOT] += ((playerSkills[p][SK_SPE] + playerSkills[p][SK_PAS]) * penalty * fwdModFactors[IDX_MOVE2ATTACK])/TENTHOUSAND_SQ;
                globSkills[IDX_ENDURANCE] += ((playerSkills[p][SK_END]) * penalty)/TENTHOUSAND;
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
    function getExtraAttackFactors(bool extraAttack) public pure returns (uint256[3] memory fwdModFactors) {
        if (extraAttack)    {fwdModFactors = [uint256(10500), uint256(10500), uint256(9500)];}
        else                {fwdModFactors = [TENTHOUSAND, TENTHOUSAND, TENTHOUSAND];}
    }
  
    // 0 penalty means no penalty
    // 1000 penalty means 10% penalty
    // etc... up to MAX_PENALTY
    function computePenalty(
        uint8 lineupPos, 
        uint8[9] memory playersPerZone, 
        uint8 forwardness,
        uint8 leftishness
    ) 
        public
        pure
        returns (uint256 penalty) 
    {
        require(lineupPos < 11, "wrong arg in computePenalty");
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

