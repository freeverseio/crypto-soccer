pragma solidity ^0.5.0;

import "./EncodingSkills.sol";

contract Engine is EncodingSkills{
    
    uint256 constant public FREE_PLAYER_ID  = 1; // it never corresponds to a legit playerId due to its TZ = 0
    uint8 public constant ROUNDS_PER_MATCH  = 12;   // Number of relevant actions that happen during a game (12 equals one per 3.7 min)
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
    uint16 private constant NO_EVENT        = 11; 

    
    bool dummyBoolToEstimateCost;

    // mock up to estimate cost of a match.
    // to be removed before deployment
    function playMatchWithCost(
        uint256 seed,
        uint256[PLAYERS_PER_TEAM_MAX][2] memory states,
        uint256[2] memory tactics,
        uint16[7][2] memory events1stHalf,
        bool is2ndHalf,
        bool isHomeStadium
    )
        public
    {
        playMatch(seed, states, tactics, events1stHalf, is2ndHalf, isHomeStadium);
        dummyBoolToEstimateCost = !dummyBoolToEstimateCost; 
    }
    
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
        uint256[2] memory tactics,
        uint16[7][2] memory events1stHalf,
        bool is2ndHalf,
        bool isHomeStadium
    )
        public
        pure
        returns (uint8[2] memory teamGoals) 
    {
        uint8[9][2] memory playersPerZone;
        uint64[] memory rnds = getNRandsFromSeed(seed, ROUNDS_PER_MATCH*4);
        uint256[5][2] memory globSkills;
        uint256[5][2] memory deltaSkills;
        bool[10][2] memory extraAttack;
        
        
        (states[0], extraAttack[0], playersPerZone[0]) = getLineUpAndPlayerPerZone(states[0], tactics[0], events1stHalf, is2ndHalf);
        (states[1], extraAttack[1], playersPerZone[1]) = getLineUpAndPlayerPerZone(states[1], tactics[1], events1stHalf, is2ndHalf);

        uint16[7][2] memory events;
        events[0] = computeExceptionalEvents(states[0], events1stHalf[0], seed);
        events[1] = computeExceptionalEvents(states[1], events1stHalf[1], seed);

        // events[0,1,2] contain the [playerPos, round, typeOfOutOfGame] data for a possible outOfGamePlayer. It is zero otherwise
        // - playerPos is set to NO_EVENT if nothing happened.
        (globSkills[0], deltaSkills[0]) = getTeamGlobSkills(states[0], playersPerZone[0], extraAttack[0], events[0][0]);
        (globSkills[1], deltaSkills[1]) = getTeamGlobSkills(states[1], playersPerZone[1], extraAttack[1], events[1][0]);
        if (isHomeStadium) {
            globSkills[IDX_ENDURANCE][0] = (globSkills[IDX_ENDURANCE][0] * 11500)/10000;
        }
        uint8 teamThatAttacks;
        for (uint8 round = 0; round < ROUNDS_PER_MATCH; round++){
            if (is2ndHalf && ((round == 0) || (round == 5))) {
                (globSkills[0], globSkills[1]) = teamsGetTired(globSkills[0], globSkills[1]);
            }
            teamThatAttacks = throwDice(globSkills[0][IDX_MOVE2ATTACK], globSkills[1][IDX_MOVE2ATTACK], rnds[4*round]);
            if ( managesToShoot(teamThatAttacks, globSkills, rnds[4*round+1])) {
                if ( managesToScore(
                    states[teamThatAttacks],
                    playersPerZone[teamThatAttacks],
                    extraAttack[teamThatAttacks],
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
    
    // // encodes: round (from 0 to 11), linedUpPos (from 0 to 10), eventType (from 1 to 3: redCard, softInjury, HardInjury)
    // // encoded value = 0 <==> no event
    // function encodeEvent(uint16 linedUpPos, uint16 round, uint16 eventType) public pure returns (uint16) {
    //     return (linedUpPos) + (round << 4) + (eventType << 8);
    // }
    
    // function decodeEvent(uint16 code) public pure returns (uint16[3] memory) {
    //     return [code & 15, (code >> 4) & 15, (code >> 8) & 3];
    // }
    
    function computeTypeOfEvent(uint256 rnd) private pure returns (uint8) {
        uint256[] memory weights = new uint256[](3);
        weights[0] = 1; // injuryHard   
        weights[1] = 2; // injuryLow
        weights[2] = 5; // redCard
        return 1 + throwDiceArray(weights, rnd);
    }

    // Over a game, we would like:
    //      - injuryHard = 1 per 100 games => 0.01 per game per player => 0.02 per game
    //      - injuryLow = 0.7 per 100 games => 0.007 per game per player => 0.04 per game
    //      - redCard 1/10 = 0.1 per game
    //      - yellowCard 2.5 per game 
    // We encode this in uint16[3] events, which applies to 1 half of the game only.
    //  - 1 possible event that leaves a player out of the match, encoded in:
    //          events[0, 1, 2] = [player (from 0 to 11), round, eventType (injuryHard, injuryLow, redCard)]
    //  - 2 possible events for yellow card:
    //          events[3, 4] = [player (from 0 to 11), round]
    //          events[5, 6] = [player (from 0 to 11), round]
    //  The player value is set to NO_EVENT ( = 11) if no event took place
    function computeExceptionalEvents
    (
        uint256[PLAYERS_PER_TEAM_MAX] memory states, 
        uint16[7] memory events1stHalf,
        uint256 seed
    ) 
        public 
        pure 
        returns 
    (
        uint16[7] memory events
    ) 
    {
        uint256[] memory weights = new uint256[](12);
        uint64[] memory rnds = getNRandsFromSeed(seed + 42, 4);
        for (uint8 p = 0; p < 11; p++) {
            weights[p] = 1 + getAggressiveness(states[p]); // weights must be > 0 to ever be selected
        }
        // events[0] => STUFF THAT REMOVES A PLAYER FROM FIELD: injuries and redCard 
        // average sumAggressiveness = 11 * 2.5 = 27.5
        // total = 0.07 per game = 0.035 per half => weight nothing happens = 758
        weights[11] = 758;
        uint8 outOfFieldPlayer = throwDiceArray(weights, rnds[0]);
        events[0] = outOfFieldPlayer;
        events[1] = uint8(rnds[0] % ROUNDS_PER_MATCH);
        events[2] = computeTypeOfEvent(rnds[1]);
        // two events for yellow cards
        // average sumAggressiveness = 11 * 2.5 = 27.5
        // total = 2.5 per game = 1.25 per half => 0.75 per dice thrown
        // weight nothing happens = 9
        weights[11] = 9;
        outOfFieldPlayer = throwDiceArray(weights, rnds[2]);
        events[3] = outOfFieldPlayer;
        events[4] = uint16(rnds[2] % ROUNDS_PER_MATCH);
        outOfFieldPlayer = throwDiceArray(weights, rnds[3]);
        events[5] = outOfFieldPlayer;
        events[6] = uint16(rnds[3] % ROUNDS_PER_MATCH);
    }
    


    // translates from a high level tacticsId (e.g. 442) to a format that describes how many
    // players play in each of the 9 zones in the field (Def, Mid, Forw) x (L, C, R), 
    // We impose left-right symmetry: DR = DL, MR = ML, FR = FL.
    // So we only manage 6 numbers: [DL, DM, ML, MM, FL, FM], and force 
    // this function returns an array of (still) size 25, but only with the correctly linedup first 11 entries filled.
    function getLineUpAndPlayerPerZone(
        uint256[PLAYERS_PER_TEAM_MAX] memory states, 
        uint256 tactics,
        uint16[7][2] memory events1stHalf,
        bool is2ndHalf
    ) 
        public 
        pure 
        returns (uint256[PLAYERS_PER_TEAM_MAX] memory outStates, bool[10] memory extraAttack, uint8[9] memory playersPerZone) 
    {
        uint8 tacticsId;
        uint8[11] memory lineup;
        uint8 changes;
        (lineup, extraAttack, tacticsId) = decodeTactics(tactics);
        for (uint8 p = 0; p < 11; p++) 
        {
            outStates[p] = states[lineup[p]];
            assertCanPlay(outStates[p]);
            if (is2ndHalf && !getAlignedLastHalf(outStates[p])) changes++;
        }
        require(changes < 4, "max allowed changes during the break is 3");
        return (states, extraAttack, getPlayersPerZone(tacticsId));
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


    function getNRandsFromSeed(uint256 seed, uint8 nRnds) public pure returns (uint64[] memory) {
        uint256 currentBigRnd = uint256(keccak256(abi.encode(seed)));
        uint8 remainingBits = 255;
        uint64[] memory rnds = new uint64[](nRnds);
        for (uint8 n = 0; n < nRnds; n++) {
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
    function throwDiceArray(uint256[] memory weights, uint256 rndNum) public pure returns(uint8 w) {
        uint256 uniformRndInSumOfWeights;
        for (w = 0; w < weights.length; w++) {
            uniformRndInSumOfWeights += weights[w];
        }
        uniformRndInSumOfWeights *= rndNum;
        uint256 cumSum = 0;
        for (w = 0; w < weights.length-1; w++) {
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
        uint256[PLAYERS_PER_TEAM_MAX] memory teamState,
        uint8[9] memory playersPerZone,
        bool[10] memory extraAttack,
        uint8 shooter,
        uint256 rnd
    )
        public
        pure
        returns (uint8)
    {
        uint256[] memory weights = new uint256[](11);
        // if selected assister == selected shooter =>  
        //  there was no assist => individual play by shoorter
        weights[0] = 1;
        uint256 teamPassCapacity = 1;
        uint8 p = 1;
        for (uint8 i = 0; i < getNDefenders(playersPerZone); i++) {
            weights[p] = (extraAttack[p-1] ? 90 : 20 ) * getPass(teamState[p]);
            teamPassCapacity += weights[p];
            p++;
        }
        for (uint8 i = 0; i < getNMidfielders(playersPerZone); i++) {
            weights[p] = (extraAttack[p-1] ? 150 : 100 ) * getPass(teamState[p]);
            teamPassCapacity += weights[p];
            p++;
        }
        for (uint8 i = 0; i < getNAttackers(playersPerZone); i++) {
            weights[p] = 200 * getPass(teamState[p]);
            teamPassCapacity += weights[p];
            p++;
        }
        // on average: teamPassCapacity442 = (1 + 4 * 20 + 4 * 100 + 2 * 200) < getPass > = 881 <pass>_team
        // on average: shooterSumOfSkills = 5 * <skills>_shooter
        // so a good ratio is shooterSumOfSkills/teamPassCapacity442 = 5/881 * <skills_shooter>/<pass>_team
        // or better, to have an avg of 1: (shooterSumOfSkills*271)/(teamPassCapacity * 5) = <skills_shooter>/<pass>_team
        // or to have a 50% change, multiply by 10, and to have say, 1/3, multiply by 10/3
        // this is to be compensated by an overall factor of about.
        weights[shooter] = (weights[shooter] * getSumOfSkills(teamState[shooter]) * 8810 )/ (N_SKILLS * (teamPassCapacity - weights[shooter]) * 3);
        return throwDiceArray(weights, rnd);
    }


    function selectShooter(
        uint256[PLAYERS_PER_TEAM_MAX] memory teamState,
        uint8[9] memory playersPerZone,
        bool[10] memory extraAttack,
        uint256 rnd
    )
        public
        pure
        returns (uint8)
    {
        uint256[] memory weights = new uint256[](11);
        // GK has minimum weight, all others are relative to this.
        weights[0] = 1;
        uint8 p = 1;
        for (uint8 i = 0; i < getNDefenders(playersPerZone); i++) {
            weights[p] = (extraAttack[p-1] ? 15000 : 5000 ) * getSpeed(teamState[p]);
            p++;
        }
        for (uint8 i = 0; i < getNMidfielders(playersPerZone); i++) {
            weights[p] = (extraAttack[p-1] ? 50000 : 25000 ) * getSpeed(teamState[p]);
            p++;
        }
        for (uint8 i = 0; i < getNAttackers(playersPerZone); i++) {
            weights[p] = 75000 * getSpeed(teamState[p]);
            p++;
        }
        return throwDiceArray(weights, rnd);
    }

    /// @dev Decides if a team that creates a shoot manages to score.
    /// @dev First: select attacker who manages to shoot. Second: challenge him with keeper
    function managesToScore(
        uint256[PLAYERS_PER_TEAM_MAX] memory teamState,
        uint8[9] memory playersPerZone,
        bool[10] memory extraAttack,
        uint256 blockShoot,
        uint256 rndNum1,
        uint256 rndNum2
    )
        public
        pure
        returns (bool)
    {
        uint8 shooter = selectShooter(teamState, playersPerZone, extraAttack, rndNum1);

        /// a goal is scored by confronting his shoot skill to the goalkeeper block skill
        uint256 shootPenalty = getForwardness(teamState[shooter]) == IDX_GK ? 10 : 1;
        return throwDice((getShoot(teamState[shooter])*7)/(shootPenalty*10), blockShoot, rndNum2) == 0;
    }
    
    function assertCanPlay(uint256 playerSkills) public pure {
        require(getPlayerIdFromSkills(playerSkills) != FREE_PLAYER_ID, "free player shirt has been aligned");
        require(!getRedCardLastGame(playerSkills) && getInjuryWeeksLeft(playerSkills) == 0, "player injured or sanctioned");
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
        bool[10] memory extraAttack,
        uint16 outOfGamePlayerPos
    )
        public
        pure
        returns 
    (
            uint256[5] memory globSkills,
            uint256[5] memory deltaSkills
    )
    {
        // for a keeper, the 'shoot skill' is interpreted as block skill
        // if for whatever reason, user places a non-GK as GK, the block skill is a terrible minimum.
        uint256 penalty;
        uint256 playerSkills = teamState[0];
        globSkills[IDX_ENDURANCE] = getEndurance(playerSkills);
        if (computePenaltyBadPositionAndCondition(0, playersPerZone, playerSkills) == 0) {globSkills[IDX_BLOCK_SHOOT] = 10;}
        else globSkills[IDX_BLOCK_SHOOT] = getShoot(playerSkills);
            
        
        uint256[3] memory fwdModFactors;
        uint8 p = 1;
        // loop over defenders
        for (uint8 i = 0; i < getNDefenders(playersPerZone); i++) {
            playerSkills = teamState[p];
            penalty = computePenaltyBadPositionAndCondition(p, playersPerZone, playerSkills);
            fwdModFactors = getExtraAttackFactors(extraAttack[p-1]);
            deltaSkills = computeDefenderGlobSkills(playerSkills, penalty, fwdModFactors);
            globSkills[IDX_MOVE2ATTACK]  += deltaSkills[IDX_MOVE2ATTACK];
            globSkills[IDX_DEFEND_SHOOT] += deltaSkills[IDX_DEFEND_SHOOT];
            globSkills[IDX_ENDURANCE]    += deltaSkills[IDX_ENDURANCE];
            p++;
        }
        // loop over midfielders
        for (uint8 i = 0; i < getNMidfielders(playersPerZone); i++) {
            playerSkills = teamState[p];
            penalty = computePenaltyBadPositionAndCondition(p, playersPerZone, playerSkills);
            fwdModFactors = getExtraAttackFactors(extraAttack[p-1]);
            deltaSkills = computeMidfielderGlobSkills(playerSkills, penalty, fwdModFactors);
            globSkills[IDX_MOVE2ATTACK] += deltaSkills[IDX_MOVE2ATTACK];
            globSkills[IDX_ENDURANCE]   += deltaSkills[IDX_ENDURANCE];
            p++;
        }
        // loop over strikers
        for (uint8 i = 0; i < getNAttackers(playersPerZone); i++) {
            playerSkills = teamState[p];
            penalty = computePenaltyBadPositionAndCondition(p, playersPerZone, playerSkills);
            fwdModFactors = getExtraAttackFactors(extraAttack[p-1]);
            deltaSkills = computeForwardsGlobSkills(playerSkills, penalty, fwdModFactors);
            globSkills[IDX_MOVE2ATTACK] += deltaSkills[IDX_MOVE2ATTACK];
            globSkills[IDX_CREATE_SHOOT]+= deltaSkills[IDX_CREATE_SHOOT];
            globSkills[IDX_ENDURANCE]   += deltaSkills[IDX_ENDURANCE];
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

    function computeDefenderGlobSkills(
        uint256 playerSkills, 
        uint256 penalty, 
        uint256[3] memory fwdModFactors
    ) 
        private 
        pure
        returns (uint256[5] memory deltaSkills) 
    {
        if (penalty != 0) {
            deltaSkills[IDX_MOVE2ATTACK] += ((getDefence(playerSkills) + getSpeed(playerSkills) + getPass(playerSkills)) * penalty * fwdModFactors[IDX_MOVE2ATTACK])/TENTHOUSAND_SQ;
            deltaSkills[IDX_DEFEND_SHOOT] += ((getDefence(playerSkills) + getSpeed(playerSkills)) * penalty * fwdModFactors[IDX_MOVE2ATTACK])/TENTHOUSAND_SQ;
            deltaSkills[IDX_ENDURANCE]   += ((getEndurance(playerSkills)) * penalty)/TENTHOUSAND;
        } else {
            deltaSkills[IDX_MOVE2ATTACK] += 30;
            deltaSkills[IDX_DEFEND_SHOOT] += 20;
            deltaSkills[IDX_ENDURANCE]   += 10;
        }
    }


    function computeMidfielderGlobSkills(
        uint256 playerSkills, 
        uint256 penalty, 
        uint256[3] memory fwdModFactors
    ) 
        private 
        pure
        returns (uint256[5] memory deltaSkills) 
    {
        if (penalty != 0) {
            deltaSkills[IDX_MOVE2ATTACK] += ((2*getDefence(playerSkills) + 2*getSpeed(playerSkills) + 3*getPass(playerSkills)) * penalty * fwdModFactors[IDX_MOVE2ATTACK])/TENTHOUSAND_SQ;
            deltaSkills[IDX_ENDURANCE]   += ((getEndurance(playerSkills)) * penalty)/TENTHOUSAND;
        } else {
            deltaSkills[IDX_MOVE2ATTACK] += 50;
            deltaSkills[IDX_ENDURANCE]   += 10;
        }
    }
    
    
    function computeForwardsGlobSkills(
        uint256 playerSkills, 
        uint256 penalty, 
        uint256[3] memory fwdModFactors
    ) 
        private 
        pure
        returns (uint256[5] memory deltaSkills) 
    {
        if (penalty != 0) {
            deltaSkills[IDX_MOVE2ATTACK] += ((getDefence(playerSkills)) * penalty * fwdModFactors[IDX_MOVE2ATTACK])/TENTHOUSAND_SQ;
            deltaSkills[IDX_CREATE_SHOOT] += ((getSpeed(playerSkills) + getPass(playerSkills)) * penalty * fwdModFactors[IDX_MOVE2ATTACK])/TENTHOUSAND_SQ;
            deltaSkills[IDX_ENDURANCE] += ((getEndurance(playerSkills)) * penalty)/TENTHOUSAND;
        } else {
            deltaSkills[IDX_MOVE2ATTACK] += 10;
            deltaSkills[IDX_CREATE_SHOOT] += 20;
            deltaSkills[IDX_ENDURANCE] += 10;
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

